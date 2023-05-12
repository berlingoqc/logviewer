package factory

import (
	"errors"

	"github.com/berlingoqc/logexplorer/pkg/log/client"
	"github.com/berlingoqc/logexplorer/pkg/log/config"
	"github.com/berlingoqc/logexplorer/pkg/log/k8s"
	"github.com/berlingoqc/logexplorer/pkg/log/local"
	"github.com/berlingoqc/logexplorer/pkg/log/opensearch"
	"github.com/berlingoqc/logexplorer/pkg/log/ssh"
	"github.com/berlingoqc/logexplorer/pkg/ty"
)

type logClientFactory struct {
	clients ty.LazyMap[string, client.LogClient]
}

func GetLogClientFactory(clients config.Clients) (*logClientFactory, error) {

	logClientFactory := new(logClientFactory)
	logClientFactory.clients = make(ty.LazyMap[string, client.LogClient])

	for k, v := range clients {
		switch v.Type {
		case "opensearch":
			options := v.Options
			logClientFactory.clients[k] = ty.GetLazy(func() (*client.LogClient, error) {
				vv, err := opensearch.GetClient(opensearch.OpenSearchTarget{
					Endpoint: options["Endpoint"],
				})
				if err != nil {
					return nil, err
				}

				return &vv, nil
			})
		case "local":
			logClientFactory.clients[k] = ty.GetLazy(func() (*client.LogClient, error) {
				vv, err := local.GetLogClient()
				if err != nil {
					return nil, err
				}

				return &vv, nil
			})
		case "k8s":
			logClientFactory.clients[k] = ty.GetLazy(func() (*client.LogClient, error) {
				vv, err := k8s.GetLogClient(k8s.K8sLogClientOptions{
					KubeConfig: v.Options["KubeConfig"],
				})
				if err != nil {
					return nil, err
				}

				return &vv, nil
			})
		case "ssh":
			logClientFactory.clients[k] = ty.GetLazy(func() (*client.LogClient, error) {
				vv, err := ssh.GetLogClient(ssh.SSHLogClientOptions{
					User:       v.Options["User"],
					Addr:       v.Options["Addr"],
					PrivateKey: v.Options["PrivateKey"],
				})
				if err != nil {
					return nil, err
				}

				return &vv, nil
			})
		default:
			return nil, errors.New("invalid type for client : " + v.Type)
		}
	}

	return logClientFactory, nil
}
