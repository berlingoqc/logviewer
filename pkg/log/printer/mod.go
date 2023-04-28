package printer

import "github.com/berlingoqc/logexplorer/pkg/log/client"

type LogPrinter interface {
	Append(entry []client.LogEntry) error
}
