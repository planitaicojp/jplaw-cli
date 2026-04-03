package list

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
	flagLawType          []string
	flagEra              string
	flagTitle            string
	flagPromulgationFrom string
	flagPromulgationTo   string
	flagCategory         []string
	flagRepealStatus     []string
	flagLimit            int
	flagOffset           int
)

var Cmd = &cobra.Command{
	Use:   "list",
	Short: "法令一覧を取得",
	Long:  "条件に一致する法令の一覧を取得します。",
	RunE:  run,
}

func init() {
	Cmd.Flags().StringSliceVar(&flagLawType, "law-type", nil, "法令種別（法律,政令,府省令等、複数可）")
	Cmd.Flags().StringVar(&flagEra, "era", "", "時代（明治,大正,昭和,平成,令和）")
	Cmd.Flags().StringVar(&flagTitle, "title", "", "法令名（部分一致）")
	Cmd.Flags().StringVar(&flagPromulgationFrom, "promulgation-from", "", "公布日FROM（YYYY-MM-DD）")
	Cmd.Flags().StringVar(&flagPromulgationTo, "promulgation-to", "", "公布日TO（YYYY-MM-DD）")
	Cmd.Flags().StringSliceVar(&flagCategory, "category", nil, "分類コード")
	Cmd.Flags().StringSliceVar(&flagRepealStatus, "repeal-status", nil, "廃止状態")
	Cmd.Flags().IntVar(&flagLimit, "limit", 100, "結果数上限")
	Cmd.Flags().IntVar(&flagOffset, "offset", 0, "開始位置")
}

func run(cmd *cobra.Command, args []string) error {
	client, err := cmdutil.NewClient()
	if err != nil {
		return err
	}

	params := &api.LawsParams{
		LawTitle:         flagTitle,
		PromulgationFrom: flagPromulgationFrom,
		PromulgationTo:   flagPromulgationTo,
		CategoryCd:       flagCategory,
		Limit:            flagLimit,
		Offset:           flagOffset,
	}

	for _, lt := range flagLawType {
		apiType := model.LawTypeLabelFromAPI(lt)
		if apiType == "" {
			params.LawType = append(params.LawType, model.LawType(lt))
		} else {
			params.LawType = append(params.LawType, apiType)
		}
	}

	if flagEra != "" {
		era := model.EraFromLabel(flagEra)
		if era == "" {
			params.LawNumEra = model.Era(flagEra)
		} else {
			params.LawNumEra = era
		}
	}

	for _, rs := range flagRepealStatus {
		params.RepealStatus = append(params.RepealStatus, model.RepealStatus(rs))
	}

	resp, err := client.ListLaws(params)
	if err != nil {
		return err
	}

	format := cmdutil.GetFormat(cmd, "table")

	if format == "json" {
		return output.New("json").Format(os.Stdout, resp)
	}

	if len(resp.Laws) == 0 {
		fmt.Fprintln(os.Stderr, "該当する法令がありません")
		return nil
	}

	rows := make([]model.LawListRow, len(resp.Laws))
	for i, law := range resp.Laws {
		rows[i] = model.LawListRow{
			LawID:        law.LawInfo.LawID,
			LawNum:       law.LawInfo.LawNum,
			LawTitle:     law.RevisionInfo.LawTitle,
			LawType:      law.LawInfo.LawType.Label(),
			Promulgation: law.LawInfo.Promulgation,
		}
	}

	fmt.Fprintf(os.Stderr, "法令一覧: %d/%d件\n", resp.Count, resp.TotalCount)
	return output.New("table").Format(os.Stdout, rows)
}
