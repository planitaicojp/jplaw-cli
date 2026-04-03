package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/planitaicojp/jplaw-cli/cmd/attachment"
	"github.com/planitaicojp/jplaw-cli/cmd/get"
	"github.com/planitaicojp/jplaw-cli/cmd/history"
	"github.com/planitaicojp/jplaw-cli/cmd/list"
	"github.com/planitaicojp/jplaw-cli/cmd/search"
	"github.com/planitaicojp/jplaw-cli/internal/api"
	cerrors "github.com/planitaicojp/jplaw-cli/internal/errors"
)

var (
	version = "dev"

	flagFormat  string
	flagVerbose bool
	flagNoColor bool
	flagNoInput bool
)

var rootCmd = &cobra.Command{
	Use:           "jplaw",
	Short:         "e-Gov法令API CLI",
	Long:          "e-Gov法令APIを操作するためのコマンドラインツール",
	SilenceUsage:  true,
	SilenceErrors: true,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		api.UserAgent = "jplaw-cli/" + version
		if flagVerbose {
			api.SetDebugLevel(api.DebugVerbose)
		}
	},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&flagFormat, "format", "", "出力フォーマット: table, json, text")
	rootCmd.PersistentFlags().BoolVar(&flagVerbose, "verbose", false, "詳細出力（HTTPデバッグ）")
	rootCmd.PersistentFlags().BoolVar(&flagNoColor, "no-color", false, "色出力を無効化")
	rootCmd.PersistentFlags().BoolVar(&flagNoInput, "no-input", false, "対話プロンプトを無効化")

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(completionCmd)
	rootCmd.AddCommand(search.Cmd)
	rootCmd.AddCommand(list.Cmd)
	rootCmd.AddCommand(get.Cmd)
	rootCmd.AddCommand(history.Cmd)
	rootCmd.AddCommand(attachment.Cmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(cerrors.GetExitCode(err))
	}
}
