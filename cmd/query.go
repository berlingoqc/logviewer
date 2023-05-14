package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strings"

	"github.com/berlingoqc/logexplorer/pkg/log/client"
	"github.com/berlingoqc/logexplorer/pkg/log/config"
	"github.com/berlingoqc/logexplorer/pkg/log/elk/kibana"
	"github.com/berlingoqc/logexplorer/pkg/log/elk/opensearch"
	"github.com/berlingoqc/logexplorer/pkg/log/factory"
	"github.com/berlingoqc/logexplorer/pkg/log/k8s"
	"github.com/berlingoqc/logexplorer/pkg/log/local"
	"github.com/berlingoqc/logexplorer/pkg/log/printer"
	"github.com/berlingoqc/logexplorer/pkg/log/ssh"
	"github.com/berlingoqc/logexplorer/pkg/ty"
	"github.com/berlingoqc/logexplorer/pkg/views"

	"github.com/spf13/cobra"
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

	// resolve this from args
	searchRequest := client.LogSearch{
		Fields:          ty.MS{},
		FieldsCondition: ty.MS{},
		Options:         ty.MI{},
	}
	if size > 0 {
		searchRequest.Size.S(size)
	}
	if duration != "" {
		searchRequest.Refresh.Duration.S(duration)
	}
	if regex != "" {
		searchRequest.FieldExtraction.Regex.S(regex)
	}

	if to != "" {
		searchRequest.Range.Lte.S(to)
	}

	if from != "" {
		searchRequest.Range.Gte.S(from)
	}

	if last != "" {
		searchRequest.Range.Last.S(last)
	}

	if len(fields) > 0 {
		stringArrayEnvVariable(fields, &searchRequest.Fields)
	}

	if len(fieldsOps) > 0 {
		stringArrayEnvVariable(fieldsOps, &searchRequest.FieldsCondition)
	}

	if index != "" {
		searchRequest.Options["Index"] = index
	}

	if k8sContainer != "" {
		searchRequest.Options[k8s.FieldContainer] = k8sContainer
	}

	if k8sNamespace != "" {
		searchRequest.Options[k8s.FieldNamespace] = k8sNamespace
	}

	if k8sPod != "" {
		searchRequest.Options[k8s.FieldPod] = k8sPod
	}

	if k8sPrevious {
		searchRequest.Options[k8s.FieldPrevious] = k8sPrevious
	}

	if k8sTimestamp {
		searchRequest.Options[k8s.OptionsTimestamp] = k8sTimestamp
	}

	if cmd != "" {
		searchRequest.Options[local.OptionsCmd] = cmd
	}

	if template != "" {
		searchRequest.PrinterOptions.Template.S(template)
	}

	if len(contextIds) != 1 {
		return nil, errors.New("-i required only exactly one element when doing a query log or query tag")
	}
	if contextPath != "" || contextIds[0] != "" {
		var config config.ContextConfig
		if err := ty.ReadJsonFile(contextPath, &config); err != nil {
			return nil, err
		}

		clientFactory, err := factory.GetLogClientFactory(config.Clients)
		if err != nil {
			return nil, err
		}

		searchFactory, err := factory.GetLogSearchFactory(clientFactory, config)
		if err != nil {
			return nil, err
		}

		sr, err := searchFactory.GetSearchResult(contextIds[0], inherits, searchRequest)

		return sr, err
	} else {
		if len(inherits) > 0 {
			return nil, errors.New("--inherits is only when using --config")
		}
	}

	var err error
	var system string

	if endpointOpensearch != "" {
		system = "opensearch"
	} else if endpointKibana != "" {
		system = "kibana"
	} else if k8sNamespace != "" {
		system = "k8s"
	} else if cmd != "" {
		if sshOptions.Addr != "" {
			system = "ssh"
		} else {
			system = "local"
		}
	} else {
		return nil, errors.New(`
        failed to select a system for logging provide one of the following:
			* --kibana-endpoint
            * --openseach-endpoint
            * --k8s-namespace
            * --ssh-addr
            * --cmd
        `)
	}

	var logClient client.LogClient

	if system == "opensearch" {
		logClient, err = opensearch.GetClient(opensearch.OpenSearchTarget{Endpoint: endpointOpensearch})
	} else if system == "kibana" {
		logClient, err = kibana.GetClient(kibana.KibanaTarget{Endpoint: endpointKibana})
	} else if system == "k8s" {
		logClient, err = k8s.GetLogClient(k8s.K8sLogClientOptions{})
	} else if system == "ssh" {
		logClient, err = ssh.GetLogClient(sshOptions)
	} else {
		logClient, err = local.GetLogClient()
	}
	if err != nil {
		return nil, err
	}

	searchResult, err2 := logClient.Get(searchRequest)
	if err2 != nil {
		return nil, err2
	}

	return searchResult, nil

}

var queryFieldCommand = &cobra.Command{
	Use:    "field",
	Short:  "Dispaly available field for filtering of logs",
	PreRun: onCommandStart,
	Run: func(cmd *cobra.Command, args []string) {
		searchResult, err1 := resolveSearch()

		if err1 != nil {
			panic(err1)
		}
		searchResult.GetEntries(context.Background())
		fields, _ := searchResult.GetFields()

		for k, b := range fields {
			fmt.Printf("%s \n", k)
			for _, r := range b {
				fmt.Println("    " + r)
			}
		}

	},
}

var queryLogCommand = &cobra.Command{
	Use:    "log",
	Short:  "Display logs for system",
	PreRun: onCommandStart,
	Run: func(cmd *cobra.Command, args []string) {
		searchResult, err1 := resolveSearch()

		if err1 != nil {
			panic(err1)
		}
		outputter := printer.PrintPrinter{}
		outputter.Display(context.Background(), searchResult)
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c

	},
}

var queryCommand = &cobra.Command{
	Use:    "query",
	Short:  "Query a login system for logs and available fields",
	PreRun: onCommandStart,
	Run: func(cmd *cobra.Command, args []string) {
		var config config.ContextConfig
		if err := ty.ReadJsonFile(contextPath, &config); err != nil {
			panic(err)
		}

		if err := views.RunQueryViewApp(config, contextIds); err != nil {
			panic(err)
		}
	},
}
