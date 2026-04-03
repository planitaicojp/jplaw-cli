package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "バージョン情報を表示",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("jplaw version %s\n", version)
		fmt.Println("e-Gov法令API (https://laws.e-gov.go.jp/api/2) CLI")
	},
}
