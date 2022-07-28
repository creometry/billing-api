package utils

import (
	"path/filepath"

	"github.com/k0kubun/pp"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

var Clientset *kubernetes.Clientset

func GetKubernetesClient() *kubernetes.Clientset {
	// TODO: implement in cluster client

	home := homedir.HomeDir()
	config, err := clientcmd.BuildConfigFromFlags("", filepath.Join(home, ".kube", "config"))
	if err != nil {
		pp.Println("error", err)
		// panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		pp.Println("clientset error", err)
		// panic(err.Error())
	}

	return clientset
}
