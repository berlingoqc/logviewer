package opensearch

import (
	"context"
	"fmt"
	"time"

	"git.tmaws.io/tmconnect/logexplorer/pkg/http"
	"git.tmaws.io/tmconnect/logexplorer/pkg/log/client"
	"git.tmaws.io/tmconnect/logexplorer/pkg/ty"
)

type OpenSearchTarget struct {
    Endpoint string `json:"endpoint"`
    Index string `json:"index"`
}

type kibanaClient struct {
    target OpenSearchTarget
    client http.JsonGetClient
}


func (kc kibanaClient) Get(search client.LogSearch) (client.LogSearchResult, error) {
    var searchResult SearchResult

    request := GetSearchRequest(search)

    err := kc.client.Get("/_search", ty.MS{}, &request, &searchResult)
    if err != nil { return nil, err }

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

func (sr logSearchResult) GetEntries() ([]client.LogEntry, error) {


    size := len(sr.result.Hits.Hits)

    entries := make([]client.LogEntry, size)

    for i, h := range sr.result.Hits.Hits {
        message, b := h.Source["message"].(string)
        if !b { 
            fmt.Printf("message is not string : %+v \n", h.Source["message"])
            entries[size-i-1] = client.LogEntry{}
            continue;
        }
        if timestamp, b1 := h.Source["@timestamp"].(string); b1 {
            date, _ := time.Parse(ty.Format, timestamp)

            var level string
            if h.Source["level"] != nil {
                level, _ = h.Source["level"].(string)
            }

            entries[size-i-1] = client.LogEntry{Message: message, Timestamp: date, Level: level, Fields: h.Source }
        } else {
            fmt.Printf("timestamp is not string : %+v \n", h.Source["@timestamp"])
        }
    }

    return entries, nil
}

func (sr logSearchResult) GetTags() (client.AvailableTags, error) {

    tags := client.AvailableTags{}
    
    for _, h := range sr.result.Hits.Hits {
        SOURCE: for k, v := range h.Source {
            if k == "message" || k == "@timestamp" { continue }

            // TODO handle object
            if vString, valid := v.(string); valid {
                if tags[k] == nil { tags[k] = make([]string, 1); tags[k][0] = vString }
                for _, vv := range tags[k] {
                    if vv == v {
                        continue SOURCE
                    }
                }
                tags[k] = append(tags[k], vString)
            }
        }
    }
    return tags, nil;
}

func (sr logSearchResult) OnChange(ctx context.Context) (chan client.LogSearchResult, error) {
    if sr.search.RefreshOptions.Duration == "" {
        return nil, nil
    }

    duration, err := time.ParseDuration(sr.search.RefreshOptions.Duration)
    if err != nil { return nil, err }

    c := make(chan client.LogSearchResult, 5)
    go func () {
        for {
            select {
            case <- time.After(duration):
                {
                    sr.search.Range.Gte = sr.search.Range.Lte
                    sr.search.Range.Lte = time.Now().Format(time.RFC3339)
                    result, err1 := sr.client.Get(sr.search)
                    if err1 != nil { fmt.Println("failed to get new logs " + err1.Error()) }
                    c <- result
                }
            case <- ctx.Done():
                close(c)
                return;
            }
        }
    }()
    return c, nil
}


func GetClient(target OpenSearchTarget) client.LogClient {
    return kibanaClient{
        target: target,
        client: http.GetClient(target.Endpoint + "/" + target.Index),
    }
}
