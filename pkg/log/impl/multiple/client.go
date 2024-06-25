package multiple

import (
	"github.com/berlingoqc/logviewer/pkg/log/client"
)

type MultipleLogClient struct {
	results []client.LogSearchResult
}

func (lc *MultipleLogClient) Add(result client.LogSearchResult) error {
	lc.results = append(lc.results, result)
	return nil
}

func (lc MultipleLogClient) Get(search *client.LogSearch) (client.LogSearchResult, error) {
	return MultipleLogSearchResult{
		search:  search,
		results: lc.results,
	}, nil
}
