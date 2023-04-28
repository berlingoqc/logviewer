package client

import (
	"context"
	"time"

	"github.com/berlingoqc/logexplorer/pkg/ty"
)

type SearchRange struct {
	Lte string
	Gte string
    Last string
}

type RefreshOptions struct {
	Duration string
}

type TagExtraction struct {
	Regex string
}

type LogSearch struct {
	// Current filterring tags
	Tags ty.MS
	// Extra rules for filtering tags
	TagsCondition ty.MS

	// Range of the log query to do , depends of the system for full availability
	Range SearchRange

	// Max size of the request
	Size int

	// Refresh options for live data
	RefreshOptions RefreshOptions

	// Options to configure the implementation with specific configuration for the search
	Options ty.MI

	// Extra fields for tag extraction for system without tagging of log entry
	TagExtraction TagExtraction
}

type AvailableTags map[string][]string

func (at *AvailableTags) AddTag(key, value string) string {
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
	GetEntries(context context.Context) ([]LogEntry, chan []LogEntry, error)
	GetTags() (AvailableTags, error)
}

// Client to start a log search
type LogClient interface {
	Get(search LogSearch) (LogSearchResult, error)
}
