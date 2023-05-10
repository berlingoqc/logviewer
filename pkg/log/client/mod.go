package client

import (
	"context"
	"time"

	"github.com/berlingoqc/logexplorer/pkg/ty"
)

type AvailableFields map[string][]string

func (at *AvailableFields) AddField(key, value string) string {
	if _, ok := (*at)[key]; ok {
		for _, v := range (*at)[key] {
			if v == value {
				return ""
			}
		}
		(*at)[key] = append((*at)[key], value)
	} else {
		(*at)[key] = []string{value}
	}
	return value
}

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
	GetFields() (AvailableFields, error)
}

// Client to start a log search
type LogClient interface {
	Get(search LogSearch) (LogSearchResult, error)
}
