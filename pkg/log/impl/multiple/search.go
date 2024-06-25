package multiple

import (
	"context"
	"sort"

	"github.com/berlingoqc/logviewer/pkg/log/client"
	"github.com/berlingoqc/logviewer/pkg/ty"
)

type MultipleLogSearchResult struct {
	search *client.LogSearch

	results []client.LogSearchResult
}

func (ls MultipleLogSearchResult) GetSearch() *client.LogSearch {
	return ls.search
}

func (ls MultipleLogSearchResult) GetEntries(context context.Context) ([]client.LogEntry, chan []client.LogEntry, error) {
	// TODO: add channel support for live update
	var err error
	size := 0
	entries := make([][]client.LogEntry, len(ls.results))
	for i, v := range ls.results {
		entries[i], _, err = v.GetEntries(context)
		if err != nil {
			return nil, nil, err
		}

		for _, vv := range entries[i] {
			vv.Fields["logSearchName"] = v.GetSearch().Name
			size += 1
		}
	}

	concatEntries := make([]client.LogEntry, size)
	i := 0
	for _, v := range entries {
		for _, vv := range v {
			concatEntries[i] = vv
			i += 1
		}
	}

	// order by time
	sort.Slice(concatEntries, func(i, j int) bool {
		return concatEntries[i].Timestamp.Before(concatEntries[j].Timestamp)
	})

	return concatEntries, nil, nil
}

func (ls MultipleLogSearchResult) GetFields() (ty.UniSet[string], chan ty.UniSet[string], error) {
	return nil, nil, nil
}
