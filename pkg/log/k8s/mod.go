package k8s

import (
	"bufio"
	"context"
	"path/filepath"
	"time"

	"github.com/berlingoqc/logexplorer/pkg/log/client"
	"github.com/berlingoqc/logexplorer/pkg/log/reader"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

const (
	FieldNamespace = "namespace"
	FieldContainer = "container"
	FieldPrevious  = "previous"
	FieldPod       = "pod"

	OptionsTimestamp = "timestamp"
)

/*

* Need to support regex for pod name , to be able to get all pods from a deployment or someting
* similar and get them in the same log flow to maybe parse them afterwards
 */

type k8sLogClient struct {
	clientset *kubernetes.Clientset
}

func (lc k8sLogClient) Get(search client.LogSearch) (client.LogSearchResult, error) {

	namespace := search.Options.GetString(FieldNamespace)
	pod := search.Options.GetString(FieldPod)
	container := search.Options.GetString(FieldContainer)
	previous := search.Options.GetBool(FieldPrevious)
	timestamp := search.Options.GetBool(OptionsTimestamp)
	follow := search.RefreshOptions.Duration != ""

	tailLines := int64(search.Size)

	ipod := lc.clientset.CoreV1().Pods(namespace)

	logOptions := v1.PodLogOptions{
		TailLines:  &tailLines,
		Follow:     follow,
		Timestamps: timestamp,
		Container:  container,
		Previous:   previous,
	}

    if search.Range.Last != "" {
        lastDuration, err := time.ParseDuration(search.Range.Last)
        if err != nil { return nil, err }
        seconds := int64(lastDuration.Seconds())
        logOptions.SinceSeconds = &seconds
    } else if search.Range.Gte != "" {
		time, err := time.Parse(time.RFC3339, search.Range.Gte)
		if err != nil {
			return nil, err
		}
		metaTime := metav1.NewTime(time)
		logOptions.SinceTime = &metaTime
	}

	req := ipod.GetLogs(pod, &logOptions)

	ctx := context.Background()

	podLogs, err2 := req.Stream(ctx)
	if err2 != nil {
		return nil, err2
	}

	scanner := bufio.NewScanner(podLogs)

	return reader.GetLogResult(search, scanner, podLogs), nil
}

func GetLogClient() (client.LogClient, error) {

	kubeconfig := filepath.Join(homedir.HomeDir(), ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	logClient := k8sLogClient{
		clientset: clientset,
	}

	return logClient, nil
}
