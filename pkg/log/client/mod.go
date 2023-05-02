package client

import (
	"context"
	"time"

	"github.com/berlingoqc/logexplorer/pkg/ty"
)

type SearchRange struct {
	Lte string `json="lte"`
	Gte string `json="gte"`
    Last string `json="last"`
}

type RefreshOptions struct {
	Duration string `json="duration"`
}

type TagExtraction struct {
	Regex string `json="regex"`
}

type LogSearch struct {
	// Current filterring tags
	Tags ty.MS `json="tags"`
	// Extra rules for filtering tags
	TagsCondition ty.MS `json="tagsCondition"`

	// Range of the log query to do , depends of the system for full availability
	Range SearchRange `json="range"`

	// Max size of the request
	Size int `json="size"`

	// Refresh options for live data
	Refresh RefreshOptions `json="refresh"`

	// Options to configure the implementation with specific configuration for the search
	Options ty.MI `json="options"`

	// Extra fields for tag extraction for system without tagging of log entry
	TagExtraction TagExtraction `json="tagExtraction"`
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
