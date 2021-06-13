package controller

import (
	"log"
	"os"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// GetKubeClient build a kube client
func GetKubeClient() *kubernetes.Clientset {
	var cfg *rest.Config
	var err error
	cfg, err = rest.InClusterConfig()
	home := homeDir()
	if err != nil && home != "" {
		cfg, err = clientcmd.BuildConfigFromFlags("", home+"/.kube/config")
		exitOnErr(err)
	} else {
		cfg, err = clientcmd.BuildConfigFromFlags("", "")
		exitOnErr(err)
	}

	kubeclient, err := kubernetes.NewForConfig(cfg)
	exitOnErr(err)

	return kubeclient
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

func exitOnErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
