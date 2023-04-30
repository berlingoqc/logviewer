package factory

import (
	"errors"

	"github.com/berlingoqc/logexplorer/pkg/log/client"
	"github.com/berlingoqc/logexplorer/pkg/log/config"
	"github.com/berlingoqc/logexplorer/pkg/log/printer"
)


type logSearchFactory struct {
    clientsFactory *logClientFactory
    searchesContext config.Contexts

    config config.ContextConfig
}


func (sf *logSearchFactory) GetSearchResult(contextId string) (client.LogSearchResult, printer.PrinterOptions, error) {

    searchContext, err := sf.config.GetSearchContext(contextId)
    if err != nil {
        return nil, printer.PrinterOptions{}, err
    }

    logClient := sf.clientsFactory.clients[searchContext.Client]
    if logClient == nil {
        return nil, printer.PrinterOptions{}, errors.New("cant find client : " + searchContext.Client)
    }
    
    sr, err := logClient.Get(searchContext.Search)

    return sr, searchContext.PrinterOptions, err
}

func GetLogSearchFactory(
    f *logClientFactory,
    c config.ContextConfig,
) (*logSearchFactory, error) {

    factory := new (logSearchFactory)
    factory.searchesContext = make(config.Contexts)
    factory.clientsFactory = f
    factory.config = c

    return factory, nil
}
