package opensearch

import (
	"encoding/json"
	"fmt"
	"testing"

	"git.tmaws.io/tmconnect/logexplorer/pkg/log/client"
)



func TestBody(t *testing.T) {

    logSearch := client.LogSearch{
        Tags: map[string]string{
            "instance": "pod-1234",
            "applicationName": "mfx.services.tsapi",
        },
        Range: client.SearchRange{
            Gte: "gte",
            Lte: "lte",
        },
        Size: 100,
    }


    request := GetSearchRequest(logSearch)

    b, _ := json.MarshalIndent(&request, "", "    ")

    fmt.Printf(string(b))
}
