package controller

import (
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
	utils "pingdom_controller/general"
	extensions "k8s.io/api/extensions/v1beta1"
)

var(
	createEvent = "create"
	updateEvent = "update"
	deleteEvent = "delete"
)

func IngressInformerFactory(pc *PingdomEngine) {
	kubeclient := GetKubeClient()

	factory := informers.NewSharedInformerFactory(kubeclient, 0)
	ingressInformer := factory.Extensions().V1beta1().Ingresses()
	stopper := make(chan struct{})
	defer close(stopper)
	ingressInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			ingressEvent(obj, pc, createEvent)
		},
		UpdateFunc: func(old, new interface{}) {
			ingressEvent(new, pc, updateEvent)
		},
		DeleteFunc: func(obj interface{}) {
			ingressEvent(obj, pc, deleteEvent)
		},
	})
	ingressInformer.Informer().Run(stopper)
}

func ingressEvent(obj interface{}, pe *PingdomEngine, event string){
	ing := obj.(*extensions.Ingress)
	createPingdomCheck := utils.GetAnnotationValue(ing.Annotations, "apply")
	if createPingdomCheck == "true"{
		switch event {
		case createEvent:
			pe.addIngress <- ing
		case updateEvent:
			pe.updateIngress <- ing
		case deleteEvent:
			pe.deleteIngress <- ing.Name
		}
	}
}
