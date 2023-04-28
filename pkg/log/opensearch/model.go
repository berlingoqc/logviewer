package opensearch

import (
	"github.com/berlingoqc/logexplorer/pkg/log/client"
	"github.com/berlingoqc/logexplorer/pkg/ty"
)

type Hit struct {
	Index  string `json:"_index"`
	Type   string `json:"_type"`
	Id     string `json:"_id"`
	Score  int32  `json:"_score""`
	Source ty.MI  `json:"_source"`
}

type Hits struct {
	// total
	// max_score
	Hits []Hit `json:"hits"`
}

type SearchResult struct {
	Took int `json:"took"`
	//timeout_out
	//_shards
	Hits Hits `json:"hits"`
}

type SortItem map[string]map[string]string
type Map map[string]interface{}

type SearchRequest struct {
	Query Map        `json:"query"`
	Size  int        `json:"size"`
	Sort  []SortItem `json:"sort"`
}

func GetSearchRequest(logSearch client.LogSearch) SearchRequest {

	conditions := make([]Map, len(logSearch.Tags)+1)

	index := 0

	for k, v := range logSearch.Tags {

		op, b := logSearch.TagsCondition[k]
		if !b || op == "" {
			op = "match"
		}

		conditions[index] = Map{
			op: Map{
				k: v,
			},
		}

		index += 1
	}

	conditions[index] = Map{
		"range": Map{
			"@timestamp": Map{
				"format": "strict_date_optional_time",
				"gte":    logSearch.Range.Gte,
				"lte":    logSearch.Range.Lte,
			},
		},
	}

	query := Map{
		"bool": Map{
			"must": conditions,
		},
	}

	sortItem := SortItem{
		"@timestamp": map[string]string{
			"order":         "desc",
			"unmapped_type": "boolean",
		},
	}

	return SearchRequest{
		Query: query,
		Sort:  []SortItem{sortItem},
		Size:  logSearch.Size,
	}
}
