package history

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/planitaicojp/jplaw-cli/cmd/cmdutil"
	"github.com/planitaicojp/jplaw-cli/internal/api"
	"github.com/planitaicojp/jplaw-cli/internal/model"
	"github.com/planitaicojp/jplaw-cli/internal/output"
)

var (
	flagAmendmentFrom string
	flagAmendmentTo   string
)

var Cmd = &cobra.Command{
	Use:   "history <法令IDまたは法令番号>",
	Short: "法令の改正履歴を取得",
	Long:  "法令の改正履歴一覧を取得します。",
	Args:  cmdutil.ExactArgs(1),
	RunE:  run,
}

func init() {
	Cmd.Flags().StringVar(&flagAmendmentFrom, "amendment-from", "", "改正施行日FROM（YYYY-MM-DD）")
	Cmd.Flags().StringVar(&flagAmendmentTo, "amendment-to", "", "改正施行日TO（YYYY-MM-DD）")
}

func run(cmd *cobra.Command, args []string) error {
	if err := cmdutil.ValidateDate(flagAmendmentFrom, "amendment-from"); err != nil {
		return err
	}
	if err := cmdutil.ValidateDate(flagAmendmentTo, "amendment-to"); err != nil {
		return err
	}

	client, err := cmdutil.NewClient()
	if err != nil {
		return err
	}

	params := &api.RevisionsParams{
		AmendmentDateFrom: flagAmendmentFrom,
		AmendmentDateTo:   flagAmendmentTo,
	}

	resp, err := client.GetRevisions(args[0], params)
	if err != nil {
		return err
	}

	format := cmdutil.GetFormat(cmd, "table")

	if format == "json" {
		return output.New("json").Format(os.Stdout, resp)
	}

	if len(resp.Revisions) == 0 {
		fmt.Fprintln(os.Stderr, "改正履歴がありません")
		return nil
	}

	rows := make([]model.RevisionListRow, len(resp.Revisions))
	for i, rev := range resp.Revisions {
		rows[i] = model.RevisionListRow{
			AmendmentDate:  rev.AmendmentEnforcementDate,
			AmendmentTitle: rev.AmendmentLawTitle,
			AmendmentNum:   rev.AmendmentLawNum,
			AmendmentType:  string(rev.AmendmentType),
		}
	}

	fmt.Fprintf(os.Stderr, "%s の改正履歴: %d件\n", resp.LawInfo.LawNum, len(resp.Revisions))
	return output.New("table").Format(os.Stdout, rows)
}
