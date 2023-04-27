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
	rootCmd.PersistentFlags().StringVar(&target.Endpoint, "opensearch-endpoint", "", "opensearch endpoint")
	rootCmd.PersistentFlags().StringVar(&target.Index, "opensearch-index", "", "Search index")


    rootCmd.PersistentFlags().StringVar(&refreshOptions.Duration, "refresh-rate", "", "If provide refresh log at the rate provide (ex: 30s)")
	rootCmd.PersistentFlags().StringVar(&from, "from", "", "Get entry gte datetime date >= from")
	rootCmd.PersistentFlags().StringVar(&to, "to", "", "Get entry lte datetime date <= to")
	rootCmd.PersistentFlags().StringVar(&last, "last", "15m", "Get entry in the last duration")
	rootCmd.PersistentFlags().IntVar(&size, "size", 100, "Get entry max size")
    rootCmd.PersistentFlags().StringArrayVarP(&fields, "fields", "f", []string{}, "field for selection field=value")

	rootCmd.MarkFlagRequired("opensearch-endpoint")
	rootCmd.MarkFlagRequired("opensearch-index")


    rootCmd.AddCommand(queryCommand)
    rootCmd.AddCommand(versionCommand)
}
