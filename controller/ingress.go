package controller

import (
	"fmt"
	"log"

	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
	utils "pingdom_controller/general"
	extensions "k8s.io/api/extensions/v1beta1"
)

type IngressMetadata struct {
	Metadata   struct {
		Name        string `yaml:"name"`
		Annotations struct {
			AnnotationsMap map[string]string
		} `yaml:"annotations"`
	}
}

func IngressInformerFactory(pc *PingdomEngine) {
	fmt.Printf("Inside StreamDeployments func\n")
	kubeclient := GetKubeClient()

	factory := informers.NewSharedInformerFactory(kubeclient, 0)
	ingressInformer := factory.Extensions().V1beta1().Ingresses()
	stopper := make(chan struct{})
	defer close(stopper)
	ingressInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			ing := obj.(*extensions.Ingress)
			createPingdomCheck, _ := utils.GetAnnotationValue(ing.Annotations, "create-check")
			if createPingdomCheck == "true"{
				pc.incomingIngress <- ing
			}
		},
		UpdateFunc: func(old, new interface{}) {
			ing := new.(*extensions.Ingress)
			createPingdomCheck, _ := utils.GetAnnotationValue(ing.Annotations, "create-check")
			if createPingdomCheck == "true"{
				pc.incomingIngress <- ing
			}
		},
		DeleteFunc: func(obj interface{}) {log.Printf("\nIn DeleteFunc")},
	})

	ingressInformer.Informer().Run(stopper)
}
