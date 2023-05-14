package elk

import (
	"context"
	"fmt"
	"time"

	"github.com/berlingoqc/logexplorer/pkg/log/client"
	"github.com/berlingoqc/logexplorer/pkg/ty"
)

type Hit struct {
	Index  string `json:"_index"`
	Type   string `json:"_type"`
	Id     string `json:"_id"`
	Score  int32  `json:"_score"`
	Source ty.MI  `json:"_source"`
}

type Hits struct {
	// total
	// max_score
	Hits []Hit `json:"hits"`
}

type ElkSearchResult struct {
	client client.LogClient
	search client.LogSearch
	result Hits

	// store loaded entries

	// store extracted fields
}

func GetSearchResult(client client.LogClient, search client.LogSearch, hits Hits) ElkSearchResult {
	return ElkSearchResult{
		client: client,
		search: search,
		result: hits,
	}
}

func (sr ElkSearchResult) GetSearch() *client.LogSearch {
	return &sr.search
}

func (sr ElkSearchResult) GetEntries(context context.Context) ([]client.LogEntry, chan []client.LogEntry, error) {

	entries := sr.parseResults()

	c, err := sr.onChange(context)

	return entries, c, err
}

// TODO: fields not being updated from live data
func (sr ElkSearchResult) GetFields() (client.AvailableFields, error) {

	fields := client.AvailableFields{}

	for _, h := range sr.result.Hits {
	SOURCE:
		for k, v := range h.Source {
			if k == "message" || k == "@timestamp" {
				continue
			}

			// TODO handle object
			if vString, valid := v.(string); valid {
				if fields[k] == nil {
					fields[k] = make([]string, 1)
					fields[k][0] = vString
				}
				for _, vv := range fields[k] {
					if vv == v {
						continue SOURCE
					}
				}
				fields[k] = append(fields[k], vString)
			}
		}
	}
	return fields, nil
}

func (sr ElkSearchResult) parseResults() []client.LogEntry {
	size := len(sr.result.Hits)

	entries := make([]client.LogEntry, size)

	for i, h := range sr.result.Hits {
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

func (sr ElkSearchResult) onChange(ctx context.Context) (chan []client.LogEntry, error) {
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
					c <- result.(ElkSearchResult).parseResults()
				}
			case <-ctx.Done():
				close(c)
				return
			}
		}
	}()
	return c, nil
}
