package get

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/planitaicojp/jplaw-cli/cmd/cmdutil"
	"github.com/planitaicojp/jplaw-cli/internal/api"
	"github.com/planitaicojp/jplaw-cli/internal/lawtext"
	"github.com/planitaicojp/jplaw-cli/internal/output"
)

var (
	flagAsof   string
	flagElm    string
	flagFile   string
	flagOutput string
)

var Cmd = &cobra.Command{
	Use:   "get <法令IDまたは法令番号>",
	Short: "法令本文を取得",
	Long:  "法令IDまたは法令番号を指定して法令本文を取得します。",
	Args:  cmdutil.ExactArgs(1),
	RunE:  run,
}

func init() {
	Cmd.Flags().StringVar(&flagAsof, "asof", "", "時点指定（YYYY-MM-DD）")
	Cmd.Flags().StringVar(&flagElm, "elm", "", "特定条文指定（例: 第一条, MainProvision）")
	Cmd.Flags().StringVar(&flagFile, "file", "", "ファイルダウンロード（xml,json,html,rtf,docx）")
	Cmd.Flags().StringVarP(&flagOutput, "output", "o", "", "ファイル保存先パス")
}

func run(cmd *cobra.Command, args []string) error {
	if err := cmdutil.ValidateDate(flagAsof, "asof"); err != nil {
		return err
	}

	client, err := cmdutil.NewClient()
	if err != nil {
		return err
	}

	idOrNum := args[0]

	if flagFile != "" {
		return downloadFile(client, idOrNum)
	}

	params := &api.LawDataParams{
		Asof: flagAsof,
		Elm:  flagElm,
	}

	resp, err := client.GetLawData(idOrNum, params)
	if err != nil {
		return err
	}

	format := cmdutil.GetFormat(cmd, "text")

	switch format {
	case "json":
		return output.New("json").Format(os.Stdout, resp)
	case "xml":
		data, err := client.GetLawFile("xml", idOrNum, flagAsof)
		if err != nil {
			return err
		}
		_, err = os.Stdout.Write(data)
		return err
	default:
		if _, err := fmt.Fprintf(os.Stdout, "%s\n（%s）\n\n", resp.RevisionInfo.LawTitle, resp.LawInfo.LawNum); err != nil {
			return err
		}

		text, err := lawtext.Convert(resp.LawFullText)
		if err != nil {
			// Fallback to raw JSON
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			enc.SetEscapeHTML(false)
			return enc.Encode(resp.LawFullText)
		}
		if text != "" {
			_, err = fmt.Fprint(os.Stdout, text)
		}
		return err
	}
}

func downloadFile(client *api.Client, idOrNum string) error {
	data, err := client.GetLawFile(flagFile, idOrNum, flagAsof)
	if err != nil {
		return err
	}

	if flagOutput == "" {
		_, err = os.Stdout.Write(data)
		return err
	}

	if err := os.WriteFile(flagOutput, data, 0644); err != nil {
		return fmt.Errorf("ファイル書き込みエラー: %w", err)
	}
	fmt.Fprintf(os.Stderr, "保存しました: %s (%d bytes)\n", flagOutput, len(data))
	return nil
}
