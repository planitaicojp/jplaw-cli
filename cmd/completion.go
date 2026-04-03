package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/planitaicojp/jplaw-cli/cmd/cmdutil"
)

var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "シェル補完スクリプトを生成",
	Long: `指定されたシェルの補完スクリプトを生成します。

使用例:
  # Bash
  jplaw completion bash > /etc/bash_completion.d/jplaw

  # Zsh
  jplaw completion zsh > "${fpath[1]}/_jplaw"

  # Fish
  jplaw completion fish > ~/.config/fish/completions/jplaw.fish`,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Args:                  cmdutil.ExactArgs(1),
	DisableFlagsInUseLine: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		switch args[0] {
		case "bash":
			return rootCmd.GenBashCompletion(os.Stdout)
		case "zsh":
			return rootCmd.GenZshCompletion(os.Stdout)
		case "fish":
			return rootCmd.GenFishCompletion(os.Stdout, true)
		case "powershell":
			return rootCmd.GenPowerShellCompletionWithDesc(os.Stdout)
		default:
			return nil
		}
	},
}
