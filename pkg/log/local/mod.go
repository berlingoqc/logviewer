package local

import (
	"bufio"
	"errors"
	"os/exec"
	"strings"

	"github.com/berlingoqc/logviewer/pkg/log/client"
	"github.com/berlingoqc/logviewer/pkg/log/reader"
)

const (
	OptionsCmd = "Cmd"
)

type localLogClient struct{}

func (lc localLogClient) Get(search *client.LogSearch) (client.LogSearchResult, error) {

	cmd := search.Options.GetString(OptionsCmd)

	if cmd == "" {
		panic(errors.New("cmd is missing for localLogClient"))
	}

	splits := strings.Split(cmd, " ")

	ecmd := exec.Command(splits[0], splits[1:]...)

	stdout, err := ecmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	if err = ecmd.Start(); err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(stdout)

	return reader.GetLogResult(search, scanner, stdout), nil
}

func GetLogClient() (client.LogClient, error) {
	return localLogClient{}, nil
}
