package mock

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"

	"github.com/berlingoqc/logviewer/pkg/log/elk/opensearch"
	"github.com/berlingoqc/logviewer/pkg/ty"
)

type OpenSearchMock struct {
	cancel chan int
	server *httptest.Server
}

func (osm *OpenSearchMock) Start(mockFile string) {

	mockRes := opensearch.SearchResult{}

	err := ty.ReadJsonFile(mockFile, &mockRes)
	if err != nil {
		panic(err)
	}

	originalHits := mockRes.Hits.Hits[:]

	firstResponse := true

	index := 0

	batchSize := 10

	osm.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var bytes []byte

		if firstResponse {
			firstResponse = false
		} else {
			requestBatchSize := rand.Intn(batchSize)

			startIndex := index
			endIndex := index + requestBatchSize

			if len(originalHits) <= index+requestBatchSize {
				endIndex = len(originalHits) - 1
				index = 0
			} else {
				index = endIndex
			}

			fmt.Printf(" start %d endd %d requestBatchSize %d len original hits %d \n", startIndex, endIndex, requestBatchSize, len(originalHits))

			hits := originalHits[startIndex:endIndex]

			mockRes.Hits.Hits = hits
		}

		bytes, _ = json.Marshal(&mockRes)

		w.Write(bytes)

	}))

	fmt.Println(osm.server.URL)

	c := make(chan int, 1)

	osm.cancel = c
}

func (osm *OpenSearchMock) Stop() {
	osm.server.Close()
}
