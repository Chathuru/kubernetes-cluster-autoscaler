package functions

import (
	"flag"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"os"
	"path/filepath"
)

// LoadKubeConfig search and load Kubernetes configuration
func LoadKubeConfig() *rest.Config {
	kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")
	flag.StringVar(&kubeconfig, "kubeconfig", kubeconfig, "kubeconfig file")
	flag.Parse()

	configIncluster, err := rest.InClusterConfig()
	if err != nil {
		log.Println("[INFO] In Cluster config not found", err)
		configExternal, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("[INFO] Out of the cluste config found.")
		return configExternal
	} else {
		return configIncluster
	}
}
