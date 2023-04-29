package factory

import (
	"errors"

	"github.com/berlingoqc/logexplorer/pkg/log/client"
	"github.com/berlingoqc/logexplorer/pkg/log/config"
	"github.com/berlingoqc/logexplorer/pkg/log/k8s"
	"github.com/berlingoqc/logexplorer/pkg/log/local"
	"github.com/berlingoqc/logexplorer/pkg/log/opensearch"
	"github.com/berlingoqc/logexplorer/pkg/log/ssh"
)

type clientMaps map[string]client.LogClient

type logClientFactory struct {
	clients clientMaps
}


func GetLogClientFactory(clients config.Clients) (*logClientFactory, error) {

	logClientFactory := new(logClientFactory)
    logClientFactory.clients = make(clientMaps)

    var err error

	for k, v := range clients {
		switch v.Type {
		case "opensearch":
			logClientFactory.clients[k] = opensearch.GetClient(opensearch.OpenSearchTarget{})
        case "local":
            if logClientFactory.clients[k], err = local.GetLogClient(); err != nil {
                return nil, err
            }
        case "k8s":
            if logClientFactory.clients[k], err = k8s.GetLogClient(k8s.K8sLogClientOptions{}); err != nil {
                return nil, err
            }
        case "ssh":
            if logClientFactory.clients[k], err = ssh.GetLogClient(ssh.SSHLogClientOptions{}); err != nil {
                return nil, err
            }
		default:
			return nil, errors.New("invalid type for client : " + v.Type)
		}
	}

	return logClientFactory, nil
}
