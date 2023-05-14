package opensearch

import (
	"errors"
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

func GetSearchRequest(logSearch client.LogSearch) (SearchRequest, error) {

	conditions := make([]Map, len(logSearch.Fields)+1)

	index := 0

	var gte, lte string

	var fromDate time.Time
	var err error

	if logSearch.Size.Value == 0 {
		logSearch.Size.S(100)
	}

	if logSearch.Range.Lte.Value != "" {
		fromDate, err = time.Parse(ty.Format, logSearch.Range.Lte.Value)
		if err != nil {
			return SearchRequest{}, errors.New("can't parse lte date")
		}
		lte = logSearch.Range.Lte.Value
	} else {
		fromDate = time.Now()
		lte = fromDate.Format(ty.Format)
	}

	if logSearch.Range.Gte.Value != "" {
		gte = logSearch.Range.Gte.Value
	} else {
		if !logSearch.Range.Last.Valid {
			return SearchRequest{}, errors.New("if not Range.Gte provided must provied Range.Last")
		}
		if duration, err := time.ParseDuration(logSearch.Range.Last.Value); err == nil {
			gte = fromDate.Add(-duration).Format(ty.Format)
		} else {
			return SearchRequest{}, errors.New("can't parse duration for last : " + logSearch.Range.Last.Value)
		}
	}

	for k, v := range logSearch.Fields {

		op, b := logSearch.FieldsCondition[k]
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
				"gte":    gte,
				"lte":    lte,
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
		Size:  logSearch.Size.Value,
	}, nil
}
