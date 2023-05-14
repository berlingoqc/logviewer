package client

import (
	"context"
	"time"

	"github.com/berlingoqc/logexplorer/pkg/ty"
)

type LogEntry struct {
	Timestamp time.Time
	Message   string
	Level     string

	Fields ty.MI
}

// Result of the search , may be used to get more log
// or keep updated
type LogSearchResult interface {
	GetSearch() *LogSearch
	GetEntries(context context.Context) ([]LogEntry, chan []LogEntry, error)
	GetFields() (ty.UniSet[string], error)
}

// Client to start a log search
type LogClient interface {
	Get(search LogSearch) (LogSearchResult, error)
}
