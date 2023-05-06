package opensearch

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/berlingoqc/logexplorer/pkg/http"
	"github.com/berlingoqc/logexplorer/pkg/log/client"
	"github.com/berlingoqc/logexplorer/pkg/ty"
)

type OpenSearchTarget struct {
	Endpoint string `json:"endpoint"`
}

type kibanaClient struct {
	target OpenSearchTarget
	client http.JsonGetClient
}

func (kc kibanaClient) Get(search client.LogSearch) (client.LogSearchResult, error) {
	var searchResult SearchResult

	index := search.Options.GetString("Index")

	if index == "" {
		return nil, errors.New("index is not provided for opensearch log client")
	}

	request, err := GetSearchRequest(search)
	if err != nil {
		return nil, err
	}

	err = kc.client.Get(fmt.Sprintf("/%s/_search", index), ty.MS{}, &request, &searchResult)
	if err != nil {
		return nil, err
	}

	return logSearchResult{
		client: &kc,
		result: searchResult,
		search: search,
	}, nil
}

type logSearchResult struct {
	client *kibanaClient

	search client.LogSearch
	result SearchResult

	// store loaded entries

	// store extracted tags
}

func (sr logSearchResult) GetSearch() *client.LogSearch {
	return &sr.search
}

func (sr logSearchResult) GetEntries(context context.Context) ([]client.LogEntry, chan []client.LogEntry, error) {

	entries := sr.parseResults()

	c, err := sr.onChange(context)

	return entries, c, err
}

// TODO: tags not being updated from live data
func (sr logSearchResult) GetTags() (client.AvailableTags, error) {

	tags := client.AvailableTags{}

	for _, h := range sr.result.Hits.Hits {
	SOURCE:
		for k, v := range h.Source {
			if k == "message" || k == "@timestamp" {
				continue
			}

			// TODO handle object
			if vString, valid := v.(string); valid {
				if tags[k] == nil {
					tags[k] = make([]string, 1)
					tags[k][0] = vString
				}
				for _, vv := range tags[k] {
					if vv == v {
						continue SOURCE
					}
				}
				tags[k] = append(tags[k], vString)
			}
		}
	}
	return tags, nil
}

func (sr logSearchResult) parseResults() []client.LogEntry {
	size := len(sr.result.Hits.Hits)

	entries := make([]client.LogEntry, size)

	for i, h := range sr.result.Hits.Hits {
		message, b := h.Source["message"].(string)
		if !b {
			fmt.Printf("message is not string : %+v \n", h.Source["message"])
			entries[size-i-1] = client.LogEntry{}
			continue
		}
		if timestamp, b1 := h.Source["@timestamp"].(string); b1 {
			date, _ := time.Parse(ty.Format, timestamp)

			var level string
			if h.Source["level"] != nil {
				level, _ = h.Source["level"].(string)
			}

			entries[size-i-1] = client.LogEntry{
				Message:   message,
				Timestamp: date,
				Level:     level, Fields: h.Source}
		} else {
			fmt.Printf("timestamp is not string : %+v \n", h.Source["@timestamp"])
		}
	}

	return entries
}

func (sr logSearchResult) onChange(ctx context.Context) (chan []client.LogEntry, error) {
	if sr.search.Refresh.Duration.Value == "" {
		return nil, nil
	}

	duration, err := time.ParseDuration(sr.search.Refresh.Duration.Value)
	if err != nil {
		return nil, err
	}

	c := make(chan []client.LogEntry, 5)
	go func() {
		for {
			select {
			case <-time.After(duration):
				{
					sr.search.Range.Gte.Value = sr.search.Range.Lte.Value
					sr.search.Range.Lte.Value = time.Now().Format(time.RFC3339)
					result, err1 := sr.client.Get(sr.search)
					if err1 != nil {
						fmt.Println("failed to get new logs " + err1.Error())
					}
					c <- result.(logSearchResult).parseResults()
				}
			case <-ctx.Done():
				close(c)
				return
			}
		}
	}()
	return c, nil
}

func GetClient(target OpenSearchTarget) client.LogClient {
	return kibanaClient{
		target: target,
		client: http.GetClient(target.Endpoint),
	}
}
