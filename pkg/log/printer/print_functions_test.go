package printer_test

import (
	"testing"

	"github.com/berlingoqc/logviewer/pkg/log/printer"
	"github.com/stretchr/testify/assert"
)

func TestExpandJson(t *testing.T) {

	logEntries := []string{
		"get data from json : {\"dadaad\": 2244 }",
	}

	for _, v := range logEntries {
		expandedJson := printer.ExpandJson(v)

		assert.NotEqual(t, "", expandedJson)
	}

}
