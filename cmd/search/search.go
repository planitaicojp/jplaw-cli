package search

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/planitaicojp/jplaw-cli/cmd/cmdutil"
	"github.com/planitaicojp/jplaw-cli/internal/api"
	"github.com/planitaicojp/jplaw-cli/internal/model"
	"github.com/planitaicojp/jplaw-cli/internal/output"
)

var (
	flagLawType  []string
	flagEra      string
	flagAsof     string
	flagCategory []string
	flagLimit    int
	flagOffset   int
)

var Cmd = &cobra.Command{
	Use:   "search <キーワード>",
	Short: "法令をキーワードで全文検索",
	Long:  "法令本文をキーワードで全文検索します。ワイルドカード、AND、OR、NOT検索に対応。",
	Args:  cmdutil.ExactArgs(1),
	RunE:  run,
}

func init() {
	Cmd.Flags().StringSliceVar(&flagLawType, "law-type", nil, "法令種別（法律,政令,府省令等、複数可）")
	Cmd.Flags().StringVar(&flagEra, "era", "", "時代（明治,大正,昭和,平成,令和）")
	Cmd.Flags().StringVar(&flagAsof, "asof", "", "時点指定（YYYY-MM-DD）")
	Cmd.Flags().StringSliceVar(&flagCategory, "category", nil, "分類コード")
	Cmd.Flags().IntVar(&flagLimit, "limit", 100, "結果数上限")
	Cmd.Flags().IntVar(&flagOffset, "offset", 0, "開始位置")
}

func run(cmd *cobra.Command, args []string) error {
	if err := cmdutil.ValidateDate(flagAsof, "asof"); err != nil {
		return err
	}

	client, err := cmdutil.NewClient()
	if err != nil {
		return err
	}

	params := &api.KeywordParams{
		Keyword:    args[0],
		Limit:      flagLimit,
		Offset:     flagOffset,
		Asof:       flagAsof,
		CategoryCd: flagCategory,
	}

	params.LawType = cmdutil.ParseLawTypes(flagLawType)
	params.LawNumEra = cmdutil.ParseEra(flagEra)

	resp, err := client.SearchKeyword(params)
	if err != nil {
		return err
	}

	format := cmdutil.GetFormat(cmd, "table")

	if format == "json" {
		return output.New("json").Format(os.Stdout, resp)
	}

	if len(resp.Items) == 0 {
		fmt.Fprintln(os.Stderr, "検索結果がありません")
		return nil
	}

	rows := make([]model.KeywordListRow, 0)
	for _, item := range resp.Items {
		for _, s := range item.Sentences {
			text := stripHTML(s.Text)
			rows = append(rows, model.KeywordListRow{
				LawNum:   item.LawInfo.LawNum,
				LawTitle: item.RevisionInfo.LawTitle,
				Match:    text,
			})
		}
	}

	fmt.Fprintf(os.Stderr, "検索結果: %d件\n", resp.TotalCount)
	return output.New("table").Format(os.Stdout, rows)
}

func stripHTML(s string) string {
	result := s
	for {
		start := strings.Index(result, "<")
		if start == -1 {
			break
		}
		end := strings.Index(result[start:], ">")
		if end == -1 {
			break
		}
		result = result[:start] + result[start+end+1:]
	}
	return result
}
