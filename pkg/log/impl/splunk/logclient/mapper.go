package logclient

import (
	"strings"

	"github.com/berlingoqc/logviewer/pkg/log/client"
	"github.com/berlingoqc/logviewer/pkg/ty"
)

func getSearchRequest(logSearch *client.LogSearch) (ty.MS, error) {
	ms := ty.MS{
		"search":        "index=" + logSearch.Options.GetString("index") + "+",
		"earliest_time": logSearch.Range.Gte.Value,
		"latest_time":   logSearch.Range.Lte.Value,
	}

	searchItem := make([]string, len(logSearch.Fields))

	i := 0
	for k, v := range logSearch.Fields {
		searchItem[i] = k + "=\"" + v + "\""
		i += 1
	}

	ms["search"] += strings.Join(searchItem, "&")

	return ms, nil
}
