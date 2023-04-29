package factory

import (
	"errors"

	"github.com/berlingoqc/logexplorer/pkg/log/client"
	"github.com/berlingoqc/logexplorer/pkg/log/config"
)


type logSearchFactory struct {
    clientsFactory *logClientFactory
    searchesContext config.Contexts

    config config.ContextConfig
}


func (sf *logSearchFactory) GetSearchResult(contextId string) (client.LogSearchResult, error) {

    clientId, search, err := sf.config.GetSearchContext(contextId)
    if err != nil {
        return nil, err
    }

    logClient := sf.clientsFactory.clients[clientId]
    if logClient == nil {
        return nil, errors.New("cant find client : " + clientId)
    }
    

    return logClient.Get(search)
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
