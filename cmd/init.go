package cmd

import (
	"github.com/berlingoqc/logexplorer/pkg/log"
	"github.com/berlingoqc/logexplorer/pkg/log/opensearch"
	"github.com/berlingoqc/logexplorer/pkg/log/ssh"
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

	logger log.MyLoggerOptions

	myLog bool
)

func onCommandStart(cmd *cobra.Command, args []string) {
	log.ConfigureMyLogger(&logger)
}

func init() {
	target = opensearch.OpenSearchTarget{}

	// CONFIG

	queryCommand.PersistentFlags().StringVarP(&contextPath, "config", "c", "", "Config for preconfigure context for search")
	queryCommand.PersistentFlags().StringVarP(&contextId, "id", "i", "", "Context id to execute")

	// IMPL SPECIFIQUE

	// ME
	queryCommand.PersistentFlags().BoolVar(&myLog, "mylog", false, "read from logviewer logs file")

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

	queryCommand.PersistentFlags().StringArrayVar(&inherits, "inherits", []string{}, "When using config , list of inherits to execute on top of the one configure for the search")

	queryCommand.AddCommand(queryLogCommand)
	queryCommand.AddCommand(queryFieldCommand)

}
