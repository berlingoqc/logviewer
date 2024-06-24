package logclient

import (
	"errors"
	"time"

	"github.com/berlingoqc/logviewer/pkg/log/client"
	"github.com/berlingoqc/logviewer/pkg/log/splunk/restapi"
	"github.com/berlingoqc/logviewer/pkg/ty"
)

const maxRetryDoneJob = 3

type SplunkLogSearchClientOptions struct {
	Url string `json:"url"`

	Headers    ty.MS `json:"headers"`
	SearchBody ty.MS `json:"searchBody"`
}

type SplunkLogSearchClient struct {
	client restapi.SplunkRestClient

	options SplunkLogSearchClientOptions
}

func (s SplunkLogSearchClient) Get(search *client.LogSearch) (client.LogSearchResult, error) {

	// initiate the things and wait for query to be done

	if s.options.Headers == nil {
		s.options.Headers = ty.MS{}
	}

	if s.options.SearchBody == nil {
		s.options.SearchBody = ty.MS{}
	}

	searchRequest, err := getSearchRequest(search)
	if err != nil {
		return nil, err
	}

	searchJobResponse, err := s.client.CreateSearchJob(searchRequest["search"], searchRequest["earliest_time"], searchRequest["latest_time"], s.options.Headers, s.options.SearchBody)
	if err != nil {
		return nil, err
	}

	isDone := false
	tryCount := 0

	for !isDone || tryCount >= maxRetryDoneJob {
		time.Sleep(1 * time.Second)
		status, err := s.client.GetSearchStatus(searchJobResponse.Sid)

		if err != nil {
			panic(err)
		}

		isDone = status.Entry[0].Content.IsDone

		tryCount += 1
	}

	if tryCount >= maxRetryDoneJob {
		return nil, errors.New("number of retry for splunk job failed")
	}

	firstResult, err := s.client.GetSearchResult(searchJobResponse.Sid, 0, search.Size.Value)

	if err != nil {
		return nil, err
	}

	return SplunkLogSearchResult{
		logClient: &s,
		search:    search,
		results:   []restapi.SearchResultsResponse{firstResult},
	}, nil
}

func GetClient(options SplunkLogSearchClientOptions) (client.LogClient, error) {

	target := restapi.SplunkTarget{
		Endpoint: options.Url,
	}

	restClient, err := restapi.GetSplunkRestClient(target)
	if err != nil {
		return nil, err
	}

	client := SplunkLogSearchClient{
		client:  restClient,
		options: options,
	}

	return client, nil
}
