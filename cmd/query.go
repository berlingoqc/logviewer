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
	"github.com/berlingoqc/logexplorer/pkg/log/factory"
	"github.com/berlingoqc/logexplorer/pkg/log/k8s"
	"github.com/berlingoqc/logexplorer/pkg/log/local"
	"github.com/berlingoqc/logexplorer/pkg/log/opensearch"
	"github.com/berlingoqc/logexplorer/pkg/log/printer"
	"github.com/berlingoqc/logexplorer/pkg/log/ssh"
	"github.com/berlingoqc/logexplorer/pkg/ty"
	"github.com/berlingoqc/logexplorer/pkg/views"

	"github.com/spf13/cobra"
)

var (
	target opensearch.OpenSearchTarget
	index  string

	k8sNamespace string
	k8sPod       string
	k8sContainer string
	k8sPrevious  bool
	k8sTimestamp bool

	sshOptions ssh.SSHLogClientOptions
	cmd        string

	from string
	to   string
	last string

	fields    []string
	fieldsOps []string
	inherits  []string
	regex     string

	size int

	duration string
	refresh  bool

	template string

	contextPath string
	contextId   string
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
		Tags:          ty.MS{},
		TagsCondition: ty.MS{},
		Options:       ty.MI{},
	}
	if size > 0 {
		searchRequest.Size.S(size)
	}
	if duration != "" {
		searchRequest.Refresh.Duration.S(duration)
	}
	if regex != "" {
		searchRequest.TagExtraction.Regex.S(regex)
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
		stringArrayEnvVariable(fields, &searchRequest.Tags)
	}

	if len(fieldsOps) > 0 {
		stringArrayEnvVariable(fieldsOps, &searchRequest.TagsCondition)
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

	if contextPath != "" || contextId != "" {
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

		sr, err := searchFactory.GetSearchResult(contextId, inherits, searchRequest)

		return sr, err
	} else {
		if len(inherits) > 0 {
			return nil, errors.New("--inherits is only when using --config")
		}
	}

	var err error
	var system string

	if target.Endpoint != "" {
		system = "opensearch"
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
            * --openseach-endpoint
            * --k8s-namespace
            * --ssh-addr
            * --cmd
        `)
	}

	var logClient client.LogClient

	if system == "opensearch" {

		logClient = opensearch.GetClient(target)
	} else if system == "k8s" {

		logClient, err = k8s.GetLogClient(k8s.K8sLogClientOptions{})
		if err != nil {
			return nil, err
		}
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
		outputter := printer.PrintPrinter{}
		outputter.Display(context.Background(), searchResult)
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c

	},
}

var queryCommand = &cobra.Command{
	Use:   "query",
	Short: "Query a login system for logs and available fields",
	Run: func(cmd *cobra.Command, args []string) {
		var config config.ContextConfig
		if err := ty.ReadJsonFile(contextPath, &config); err != nil {
			panic(err)
		}

		if err := views.RunQueryViewApp(config); err != nil {
			panic(err)
		}
	},
}

func init() {
	target = opensearch.OpenSearchTarget{}

	// CONFIG

	queryCommand.PersistentFlags().StringVarP(&contextPath, "config", "c", "", "Config for preconfigure context for search")
	queryCommand.PersistentFlags().StringVarP(&contextId, "id", "i", "", "Context id to execute")

	// IMPL SPECIFIQUE

	// K8S
	queryCommand.PersistentFlags().StringVar(&k8sNamespace, "k8s-namespace", "", "K8s namespace")
	queryCommand.PersistentFlags().StringVar(&k8sPod, "k8s-pod", "", "K8s pod")
	queryCommand.PersistentFlags().StringVar(&k8sContainer, "k8s-container", "", "K8s container")
	queryCommand.PersistentFlags().BoolVar(&k8sPrevious, "k8s-previous", false, "K8s log of previous container")
	queryCommand.PersistentFlags().BoolVar(&k8sTimestamp, "k8s-timestamp", false, "K8s include RFC3339 timestamp")
	// OPENSEARCH
	queryCommand.PersistentFlags().StringVar(&target.Endpoint, "opensearch-endpoint", "", "Opensearch endpoint")
	queryCommand.PersistentFlags().StringVar(&index, "opensearch-index", "", "Opensearch index to search")
	// SSH
	queryCommand.PersistentFlags().StringVar(&sshOptions.Addr, "ssh-addr", "", "SSH address and port localhost:22")
	queryCommand.PersistentFlags().StringVar(&sshOptions.User, "ssh-user", "", "SSH user")
	queryCommand.PersistentFlags().StringVar(&sshOptions.PrivateKey, "ssh-identifiy", "", "SSH private key , by default $HOME/.ssh/id_rsa")

	// COMMAND
	queryCommand.PersistentFlags().StringVar(&cmd, "cmd", "", "If using ssh or local , manual command to run")

	// RANGE
	queryCommand.PersistentFlags().StringVar(&from, "from", "", "Get entry gte datetime date >= from")
	queryCommand.PersistentFlags().StringVar(&to, "to", "", "Get entry lte datetime date <= to")
	queryCommand.PersistentFlags().StringVar(&last, "last", "", "Get entry in the last duration")

	// SIZE
	queryCommand.PersistentFlags().IntVar(&size, "size", 0, "Get entry max size")

	// FIELD validation
	queryCommand.PersistentFlags().StringArrayVarP(&fields, "fields", "f", []string{}, "Field for selection field=value")
	queryCommand.PersistentFlags().StringArrayVar(
		&fieldsOps, "fields-condition", []string{}, "Field Ops for selection field=value (match, exists, wildcard, regex)",
	)
	queryCommand.PersistentFlags().StringVar(
		&regex, "fields-regex", "",
		"Regex to extract field from log text, using named group \".*(?P<Level>INFO|WARN|ERROR).*\"")

	// LIVE DATA OPTIONS
	queryLogCommand.PersistentFlags().StringVar(
		&duration, "refresh-rate", "", "If provide refresh log at the rate provide (ex: 30s)")
	queryLogCommand.PersistentFlags().BoolVar(&refresh, "refresh", false, "If provide activate live data")

	// OUTPUT FORMATTING
	queryLogCommand.PersistentFlags().StringVar(
		&template,
		"format",
		"[{{.Timestamp.Format \"15:04:05\" }}][{{.Level}}] {{.Message}}", "Format for the log entry")

	queryCommand.PersistentFlags().StringArrayVarP(&inherits, "inherits", "h", []string{}, "When using config , list of inherits to execute on top of the one configure for the search")

	queryCommand.AddCommand(queryLogCommand)
	queryCommand.AddCommand(queryTagCommand)

}
