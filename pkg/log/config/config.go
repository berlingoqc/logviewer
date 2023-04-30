package config

import (
	"errors"

	"github.com/berlingoqc/logexplorer/pkg/log/client"
	"github.com/berlingoqc/logexplorer/pkg/log/printer"
	"github.com/berlingoqc/logexplorer/pkg/ty"
)


type Client struct {
    Type string `json="type"`
    Options ty.MS `json="options"`
}

type SearchContext struct {
    Client string `json="client"`
    SearchInherit string `json="searchInherit"`
    Search client.LogSearch `json="search"`
    PrinterOptions printer.PrinterOptions `json="printerOptions"`
}

type Clients map[string]Client

type Searches map[string]client.LogSearch

type Contexts map[string]SearchContext

type ContextConfig struct {
    Clients
    Searches
    Contexts
}

func (cc ContextConfig) GetSearchContext(contextId string) (SearchContext, error) {
    if searchContext, b := cc.Contexts[contextId]; b {
        // TODO: inheritance
        return searchContext, nil
    } else {
        return SearchContext{}, errors.New("cant find context : " + contextId)
    }
}
