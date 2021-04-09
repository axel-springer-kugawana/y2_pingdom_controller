
package controller

import (
	"log"
	"os"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func GetKubeClient() *kubernetes.Clientset{
	var cfg *rest.Config
	var err error
	cfg, err = rest.InClusterConfig()
	home := homeDir()
	if err != nil && home != "" {
		cfg, err = clientcmd.BuildConfigFromFlags("", home+"/.kube/k3s")
		ExitOnErr(err)
	} else {
		cfg, err = clientcmd.BuildConfigFromFlags("", "")
		ExitOnErr(err)
	}

	kubeclient, err := kubernetes.NewForConfig(cfg)
	ExitOnErr(err)

	return kubeclient
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

func ExitOnErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}