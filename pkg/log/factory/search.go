package factory

import (
	"github.com/berlingoqc/logviewer/pkg/log/client"
	"github.com/berlingoqc/logviewer/pkg/log/config"
	"github.com/berlingoqc/logviewer/pkg/log/impl/multiple"
	"github.com/berlingoqc/logviewer/pkg/ty"
)

type logSearchFactory struct {
	clientsFactory  *logClientFactory
	searchesContext config.Contexts

	config config.ContextConfig
}

func (sf *logSearchFactory) GetSearchResult(
	contextId string,
	inherits []string,
	logSearch client.LogSearch,
) (client.LogSearchResult, error) {

	searchContext, err := sf.config.GetSearchContext(contextId, inherits, logSearch)
	if err != nil {
		return nil, err
	}

	if searchContext.Search.Name == "" {
		searchContext.Search.Name = contextId //fmt.Sprintf("%s:%s", searchContext.Client, contextId)
	}

	logClient, err := sf.clientsFactory.clients.Get(searchContext.Client)
	if err != nil {
		return nil, err
	}

	sr, err := (*logClient).Get(&searchContext.Search)

	return sr, err
}

func ConcatSearchResult(results []client.LogSearchResult) (client.LogSearchResult, error) {
	logClient := multiple.MultipleLogClient{}

	for _, r := range results {
		if err := logClient.Add(r); err != nil {
			return nil, err
		}
	}

	return logClient.Get(&client.LogSearch{PrinterOptions: client.PrinterOptions{
		Template: ty.OptWrap[string]("[{{.Fields.logSearchName}}] {{.Message}}"),
	}})
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
