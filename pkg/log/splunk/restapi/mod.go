package restapi

import (
	"fmt"
	"strconv"

	"github.com/berlingoqc/logviewer/pkg/http"
	"github.com/berlingoqc/logviewer/pkg/ty"
)

// Struct to hold search job response
type SearchJobResponse struct {
	Sid string `json:"sid"`
}

// Struct to hold search job status response
type JobStatusResponse struct {
	Entry []struct {
		Content struct {
			IsDone bool `json:"isDone"`
		} `json:"content"`
	} `json:"entry"`
}

// Struct to hold search results response
type SearchResultsResponse struct {
	Results []ty.MI `json:"results"`
}

type SplunkTarget struct {
	Endpoint string `json:"endpoint"`
	Auth     http.Auth
}

type SplunkRestClient struct {
	target SplunkTarget
	client http.HttpClient
}

func (src SplunkRestClient) CreateSearchJob(
	searchQuery string,
	earliestTime string,
	latestTime string,

	headers ty.MS,

	data ty.MS,
) (SearchJobResponse, error) {
	var searchJobResponse SearchJobResponse

	searchPath := fmt.Sprintf("/search/jobs")

	data["custom.dispatch.latest_time"] = latestTime
	data["search"] = `search+` + searchQuery
	data["earliest_time"] = earliestTime
	data["latest_time"] = latestTime
	data["custom.search"] = searchQuery
	data["custom.dispatch.earliest_time"] = earliestTime

	err := src.client.PostData(searchPath, headers, data, &searchJobResponse, src.target.Auth)

	return searchJobResponse, err

}

func (src SplunkRestClient) GetSearchStatus(
	sid string,
) (JobStatusResponse, error) {
	var response JobStatusResponse

	searchPath := fmt.Sprintf("/search/jobs/%s", sid)

	queryParams := ty.MS{
		"output_mode": "json",
	}

	err := src.client.Get(searchPath, queryParams, nil, &response, src.target.Auth)
	return response, err
}

func (src SplunkRestClient) GetSearchResult(
	sid string,
	offset int,
	count int,
) (SearchResultsResponse, error) {
	var response SearchResultsResponse

	searchPath := fmt.Sprintf("/search/jobs/%s/events", sid)

	queryParams := ty.MS{
		"output_mode": "json",
		"offset":      strconv.Itoa(offset),
		"count":       strconv.Itoa(count),
	}

	err := src.client.Get(searchPath, queryParams, nil, &response, src.target.Auth)
	return response, err

}

func GetSplunkRestClient(
	target SplunkTarget,
) (SplunkRestClient, error) {
	return SplunkRestClient{
		target: target,
		client: http.GetClient(target.Endpoint),
	}, nil
}
