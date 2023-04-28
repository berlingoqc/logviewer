package cmd

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/berlingoqc/logexplorer/pkg/log/client"
	"github.com/berlingoqc/logexplorer/pkg/log/k8s"
	"github.com/berlingoqc/logexplorer/pkg/log/opensearch"
	"github.com/berlingoqc/logexplorer/pkg/log/printer"
	"github.com/berlingoqc/logexplorer/pkg/ty"
	"github.com/spf13/cobra"
)

var (
	target opensearch.OpenSearchTarget

    k8sNamespace string
    k8sPod string
    k8sContainer string
    k8sPrevious bool
    k8sTimestamp bool

	from string
	to   string
	last string

	fields    []string
	fieldsOps []string
    regex       string

	size int

	refreshOptions client.RefreshOptions
    refresh bool

	outputter printer.PrintPrinter
)

func stringArrayEnvVariable(strs []string, maps *ty.MS) error {
	for _, f := range strs {
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

    var system string

    if refresh && refreshOptions.Duration == "" {
        refreshOptions.Duration = "5s"
    }

    if !refresh && refreshOptions.Duration != "" {
        refresh = true
    }

    if target.Endpoint != "" {
        system = "opensearch"
    } else if k8sNamespace != "" {
        system = "k8s"
    } else {
        return nil, errors.New("failed to select a system for logging provide --openseach-endpoint or --k8s-namespace")
    }


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
		Size:           size,
		Range:          client.SearchRange{Lte: lte, Gte: gte, Last: last},
		Tags:           ty.MS{},
		TagsCondition:  ty.MS{},
		RefreshOptions: refreshOptions,
        Options: ty.MI{},
        TagExtraction: client.TagExtraction{
            Regex: regex,
        },
	}

	stringArrayEnvVariable(fields, &searchRequest.Tags)
	stringArrayEnvVariable(fieldsOps, &searchRequest.TagsCondition)

    var logClient client.LogClient

    if system == "opensearch" {


        logClient = opensearch.GetClient(target)
    } else {

        searchRequest.Options[k8s.FieldContainer] = k8sContainer
        searchRequest.Options[k8s.FieldNamespace] = k8sNamespace
        searchRequest.Options[k8s.FieldPod] = k8sPod
        searchRequest.Options[k8s.FieldPrevious] = k8sPrevious
        searchRequest.Options[k8s.OptionsTimestamp] = k8sTimestamp

        logClient, err = k8s.GetLogClient()
        if err != nil { return nil, err }
    }


	searchResult, err2 := logClient.Get(searchRequest)
	if err2 != nil {
		return nil, err2
	}

	return searchResult, nil

}

var queryTagCommand = &cobra.Command{
	Use:   "field",
	Short: "Dispaly available field for filtering of logs",
	Run: func(cmd *cobra.Command, args []string) {
		searchResult, err1 := resolveSearch()

		if err1 != nil {
			panic(err1)
		}
        searchResult.GetEntries(context.Background())
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
	Use:   "log",
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
	Use:   "query",
	Short: "Query a login system for logs and available fields",
	Run:   func(cmd *cobra.Command, args []string) {},
}

func init() {
	target = opensearch.OpenSearchTarget{}

    // IMPL SPECIFIQUE

    // K8S
    queryCommand.PersistentFlags().StringVar(&k8sNamespace, "k8s-namespace", "", "K8s namespace")
    queryCommand.PersistentFlags().StringVar(&k8sPod, "k8s-pod", "", "K8s pod")
    queryCommand.PersistentFlags().StringVar(&k8sContainer, "k8s-container", "", "K8s container")
    queryCommand.PersistentFlags().BoolVar(&k8sPrevious, "k8s-previous", false, "K8s log of previous container")
    queryCommand.PersistentFlags().BoolVar(&k8sTimestamp, "k8s-timestamp", false, "K8s include RFC3339 timestamp")
    // OPENSEARCH
	queryCommand.PersistentFlags().StringVar(&target.Endpoint, "opensearch-endpoint", "", "Opensearch endpoint")
	queryCommand.PersistentFlags().StringVar(&target.Index, "opensearch-index", "", "Opensearch index to search")


    // RANGE
	queryCommand.PersistentFlags().StringVar(&from, "from", "", "Get entry gte datetime date >= from")
	queryCommand.PersistentFlags().StringVar(&to, "to", "", "Get entry lte datetime date <= to")
	queryCommand.PersistentFlags().StringVar(&last, "last", "15m", "Get entry in the last duration")

    // SIZE
	queryCommand.PersistentFlags().IntVar(&size, "size", 100, "Get entry max size")

    // FIELD validation
	queryCommand.PersistentFlags().StringArrayVarP(&fields, "fields", "f", []string{}, "Field for selection field=value")
	queryCommand.PersistentFlags().StringArrayVarP(
		&fieldsOps, "fields-condition", "c", []string{}, "Field Ops for selection field=value (match, exists, wildcard, regex)",
	)
    queryCommand.PersistentFlags().StringVar(&regex, "fields-regex", "", "Regex to extract field from log text, using named group \".*(?P<Level>INFO|WARN|ERROR).*\"")

    // LIVE DATA OPTIONS
	queryLogCommand.PersistentFlags().StringVar(
		&refreshOptions.Duration, "refresh-rate", "", "If provide refresh log at the rate provide (ex: 30s)")
    queryLogCommand.PersistentFlags().BoolVar(&refresh, "refresh", false, "If provide activate live data")


    // OUTPUT FORMATTING
	queryLogCommand.PersistentFlags().StringVar(
		&outputter.Options.Template,
		"format",
		"[{{.Timestamp.Format \"15:04:05\" }}][{{.Level}}] {{.Message}}", "Format for the log entry")


	queryCommand.AddCommand(queryLogCommand)
	queryCommand.AddCommand(queryTagCommand)

}
