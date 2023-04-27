package printer

import "git.tmaws.io/tmconnect/logexplorer/pkg/log/client"


type LogPrinter interface{
    Append(entry []client.LogEntry) error
}
