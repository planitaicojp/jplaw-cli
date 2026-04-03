package attachment

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/planitaicojp/jplaw-cli/cmd/cmdutil"
)

var (
	flagSrc    string
	flagOutput string
)

var Cmd = &cobra.Command{
	Use:   "attachment <リビジョンID>",
	Short: "添付ファイルをダウンロード",
	Long:  "法令リビジョンに添付されたファイルをダウンロードします。",
	Args:  cmdutil.ExactArgs(1),
	RunE:  run,
}

func init() {
	Cmd.Flags().StringVar(&flagSrc, "src", "", "特定添付ファイルのsrc属性値")
	Cmd.Flags().StringVarP(&flagOutput, "output", "o", "", "保存先パス")
}

func run(cmd *cobra.Command, args []string) error {
	client, err := cmdutil.NewClient()
	if err != nil {
		return err
	}

	data, err := client.GetAttachment(args[0], flagSrc)
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
