package client

import (
	"context"
	"time"
)

type SearchRange struct {
    Lte string
    Gte string
}


type RefreshOptions struct {
    Duration string
}

type LogSearch struct {
    Tags map[string]string

    Range SearchRange
    
    Size int

    RefreshOptions RefreshOptions
}


type AvailableTags map[string][]string


type LogEntry struct {
    Timestamp time.Time
    Message string
    Level string
}



// Result of the search , may be used to get more log
// or keep updated
type LogSearchResult interface {
    GetEntries() ([]LogEntry, error)
    OnChange(ctx context.Context) (chan []LogEntry, error)
    GetTags() (AvailableTags, error)
}

// Client to start a log search
type LogClient interface {
    Get(search LogSearch) (LogSearchResult, error)
}
