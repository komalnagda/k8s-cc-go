package connection

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"os"
	"path/filepath"
)

type K8s struct {
	Client     *kubernetes.Clientset  // interface
	RestConfig *rest.Config
}

func ConnectToK8s() *K8s {
	home, exists := os.LookupEnv("HOME")
	if !exists {
		home = "/root"
	}
	config, err := rest.InClusterConfig()
	if err != nil {
		configPath := filepath.Join(home, ".kube", "config")
		log.Println("Using kube config from ", configPath)
		config, err = clientcmd.BuildConfigFromFlags("", configPath)
		if err != nil {
			log.Fatalln("failed to create K8s config")
		}
	} else {
		log.Println("Using in-cluster configuration")
	}

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Failed to create K8s client set. Error: %s\n", err)
	}

	return &K8s{
		clientSet,
		config,
	}
}
