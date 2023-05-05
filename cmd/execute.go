package cmd

import (
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
	"log"
	"custom-controller/notification"
	"custom-controller/connection"
)

func SecretInformer() {
	k := connection.ConnectToK8s()
	factory := informers.NewSharedInformerFactory(k.Client, 0)
	informer := factory.Apps().V1().Deployments().Informer()
	stopper := make(chan struct{})
	defer close(stopper)

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				mObj, ok := obj.(*appsv1.Deployment)
				if !ok {
					log.Println("Deployment name is missing")
				}
				log.Println(mObj)
				log.Println(mObj.ObjectMeta.Name)
				notification.SendSlack(mObj.ObjectMeta.Name)
			},
		},
	)

	informer.Run(stopper)
	if !cache.WaitForCacheSync(stopper, informer.HasSynced) {
		runtime.HandleError(fmt.Errorf("timed out waiting for caches to sync"))
		return
	}
	<-stopper
}
