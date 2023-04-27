package cmd

import (
	"fmt"
	"os"

	"github.com/rivo/tview"
	"github.com/spf13/cobra"
)


var rootCmd = &cobra.Command{
  Use:   "logviewer",
  Short: "Log viewer for different backend (OpenSearch, SSH, Local Files)",
  Long: ``,
  Run: func(cmd *cobra.Command, args []string) {

      box := tview.NewBox().SetBorder(true).SetTitle("logviewer")
      if err := tview.NewApplication().SetRoot(box, true).Run(); err != nil {
          panic(err)
      }
  },
}


func Execute() {
  if err := rootCmd.Execute(); err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
}


func init() {

    rootCmd.AddCommand(queryCommand)
    rootCmd.AddCommand(versionCommand)
}
