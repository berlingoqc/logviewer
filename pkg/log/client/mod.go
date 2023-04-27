package client

import (
	"context"
	"time"

	"git.tmaws.io/tmconnect/logexplorer/pkg/ty"
)

type SearchRange struct {
    Lte string
    Gte string
}


type RefreshOptions struct {
    Duration string
}

type LogSearch struct {
    Tags ty.MS    
    TagsCondition ty.MS

    Range SearchRange
    
    Size int

    RefreshOptions RefreshOptions
}


type AvailableTags map[string][]string


type LogEntry struct {
    Timestamp time.Time
    Message string
    Level string

    Fields ty.MI
}



// Result of the search , may be used to get more log
// or keep updated
type LogSearchResult interface {
    GetEntries() ([]LogEntry, error)
    OnChange(ctx context.Context) (chan LogSearchResult, error)
    GetTags() (AvailableTags, error)
}

// Client to start a log search
type LogClient interface {
    Get(search LogSearch) (LogSearchResult, error)
}
