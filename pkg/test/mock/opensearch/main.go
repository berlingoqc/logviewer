package main

import (
	"time"

	"github.com/berlingoqc/logviewer/pkg/test/mock"
)

func main() {

	reqFile := "./pkg/test/mock/opensearch/req.json"

	mock := mock.OpenSearchMock{}

	mock.Start(reqFile)

	duration, _ := time.ParseDuration("10m")

	<-time.After(duration)

	mock.Stop()

}
