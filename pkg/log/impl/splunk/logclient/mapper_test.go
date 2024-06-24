package logclient

import (
	"testing"

	"github.com/berlingoqc/logviewer/pkg/log/client"
	"github.com/berlingoqc/logviewer/pkg/ty"
	"github.com/stretchr/testify/assert"
)

func TestSearchRequest(t *testing.T) {

	logSearch := client.LogSearch{
		Fields:  ty.MS{},
		Options: ty.MI{},
	}
	logSearch.Range.Gte.S("24h@h")
	logSearch.Range.Lte.S("now")

	logSearch.Fields["application_name"] = "wq.services.pet"
	logSearch.Fields["trace_id"] = "1234"
	logSearch.Options["index"] = "nonprod"

	requestBodyFields, err := getSearchRequest(&logSearch)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, `index=nonprod+application_name="wq.services.pet"&trace_id="1234"`, requestBodyFields["search"])
}
