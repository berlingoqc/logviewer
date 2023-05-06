package factory

import (
	"errors"

	"github.com/berlingoqc/logexplorer/pkg/log/client"
	"github.com/berlingoqc/logexplorer/pkg/log/config"
)

type logSearchFactory struct {
	clientsFactory  *logClientFactory
	searchesContext config.Contexts

	config config.ContextConfig
}

func (sf *logSearchFactory) GetSearchResult(contextId string, inherits []string, logSearch client.LogSearch) (client.LogSearchResult, error) {

	searchContext, err := sf.config.GetSearchContext(contextId, inherits, logSearch)
	if err != nil {
		return nil, err
	}

	logClient := sf.clientsFactory.clients[searchContext.Client]
	if logClient == nil {
		return nil, errors.New("cant find client : " + searchContext.Client)
	}

	sr, err := logClient.Get(searchContext.Search)

	return sr, err
}

func GetLogSearchFactory(
	f *logClientFactory,
	c config.ContextConfig,
) (*logSearchFactory, error) {

	factory := new(logSearchFactory)
	factory.searchesContext = make(config.Contexts)
	factory.clientsFactory = f
	factory.config = c

	return factory, nil
}
