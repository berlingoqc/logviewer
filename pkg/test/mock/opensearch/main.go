package main

import (
	"time"

	"github.com/berlingoqc/logexplorer/pkg/test/mock"
)

func main() {

	reqFile := "./pkg/test/mock/opensearch/req.json"

	mock := mock.OpenSearchMock{}

	mock.Start(reqFile)

	duration, _ := time.ParseDuration("10m")

	<-time.After(duration)

	mock.Stop()

}
