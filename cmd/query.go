package cmd

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"git.tmaws.io/tmconnect/logexplorer/pkg/log/client"
	"git.tmaws.io/tmconnect/logexplorer/pkg/log/opensearch"
	"git.tmaws.io/tmconnect/logexplorer/pkg/log/printer"
	"git.tmaws.io/tmconnect/logexplorer/pkg/ty"
	"github.com/spf13/cobra"
)

var (
	target opensearch.OpenSearchTarget

	from string
	to   string
	last string

    fields []string
    fieldsOps []string

	size int

    refreshOptions client.RefreshOptions
    
    outputter printer.PrintPrinter
)

func stringArrayEnvVariable(strs []string, maps *ty.MS) error {
    for _ , f := range strs {
        items := strings.Split(f, "=")
        if len(items) < 2 {
            return errors.New("invalid value : " + f)
        }
        (*maps)[items[0]] = strings.Join(items[1:], "")
    }
    return nil
}

func resolveSearch() (client.LogSearchResult, error) {
	var lte, gte string
	var fromDate time.Time
	var err error

	if from != "" {
		fromDate, err = time.Parse(ty.Format, from)
		if err != nil {
			return nil, errors.New("failed to parsed --to date : " + err.Error())
		}
		lte = from
	} else {
		fromDate = time.Now()
		lte = fromDate.Format(ty.Format)
	}

	if to != "" {
		gte = to
	} else {
		if duration, err1 := time.ParseDuration(last); err1 == nil {
			gte = fromDate.Add(-duration).Format(ty.Format)
		} else {
			return nil, errors.New("invalid --last value : " + err1.Error())
		}
	}
   
    searchRequest := client.LogSearch{
        Size: size,
        Range: client.SearchRange{Lte: lte, Gte: gte},
        Tags: ty.MS{},
        TagsCondition: ty.MS{},
        RefreshOptions: refreshOptions,
    }

    stringArrayEnvVariable(fields, &searchRequest.Tags)
    stringArrayEnvVariable(fieldsOps, &searchRequest.TagsCondition)

	openSearchClient := opensearch.GetClient(target)
	searchResult, err2 := openSearchClient.Get(searchRequest)
	if err2 != nil {
		return nil, err2
	}

	return searchResult, nil

}

var queryTagCommand = &cobra.Command{
	Use: "field",
    Short: "Dispaly available field for filtering of logs",
	Run: func(cmd *cobra.Command, args []string) {
		searchResult, err1 := resolveSearch()

		if err1 != nil {
			panic(err1)
		}
		tags, _ := searchResult.GetTags()

		for k, b := range tags {
            fmt.Printf("%s \n", k)
            for _, r := range b {
                fmt.Println("    " + r)
            }
		}

	},
}

var queryLogCommand = &cobra.Command{
	Use: "log",
    Short: "Display logs for system",
	Run: func(cmd *cobra.Command, args []string) {
		searchResult, err1 := resolveSearch()

		if err1 != nil {
			panic(err1)
		}

        outputter.Append(searchResult)
	},
}

var queryCommand = &cobra.Command{
	Use: "query",
    Short: "Query a login system for logs and available fields",
	Run: func(cmd *cobra.Command, args []string) {},
}

func init() {
	target = opensearch.OpenSearchTarget{}

	queryCommand.PersistentFlags().StringVar(&target.Endpoint, "opensearch-endpoint", "", "Opensearch endpoint")
	queryCommand.PersistentFlags().StringVar(&target.Index, "opensearch-index", "", "Opensearch index to search")


	queryCommand.PersistentFlags().StringVar(&from, "from", "", "Get entry gte datetime date >= from")
	queryCommand.PersistentFlags().StringVar(&to, "to", "", "Get entry lte datetime date <= to")
	queryCommand.PersistentFlags().StringVar(&last, "last", "15m", "Get entry in the last duration")
	queryCommand.PersistentFlags().IntVar(&size, "size", 100, "Get entry max size")
    queryCommand.PersistentFlags().StringArrayVarP(&fields, "fields", "f", []string{}, "Field for selection field=value")
    queryCommand.PersistentFlags().StringArrayVarP(
        &fieldsOps, "fields-condition", "c", []string{}, "Field Ops for selection field=value (match, exists, wildcard, regex)",
    )

    queryLogCommand.PersistentFlags().StringVar(
        &refreshOptions.Duration, "refresh-rate", "", "If provide refresh log at the rate provide (ex: 30s)")
    queryLogCommand.PersistentFlags().StringVar(
        &outputter.Options.Template,
        "format",
        "[{{.Timestamp.Format \"15:04:05\" }}][{{.Fields.applicationName}}][{{.Level}}] {{.Message}}", "Format for the log entry")



	queryCommand.MarkFlagRequired("opensearch-endpoint")
	queryCommand.MarkFlagRequired("opensearch-index")



	queryCommand.AddCommand(queryLogCommand)
	queryCommand.AddCommand(queryTagCommand)

}
