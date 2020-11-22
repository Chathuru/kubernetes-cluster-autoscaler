package main

import (
	"context"
	"encoding/json"
	"github.com/Chathuru/kubernetes-cluster-autoscaler/pkg/cloud/openstack"
	"github.com/Chathuru/kubernetes-cluster-autoscaler/pkg/cloud/openstack/handel-node-delete"
	"github.com/Chathuru/kubernetes-cluster-autoscaler/pkg/cloud/openstack/handle-node-add"
	"github.com/Chathuru/kubernetes-cluster-autoscaler/pkg/common/datastructures"
	"github.com/Chathuru/kubernetes-cluster-autoscaler/pkg/common/functions"
	_ "github.com/mattn/go-sqlite3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"log"
	"plugin"
	"sync"
)

var (
	wg sync.WaitGroup
	CloudType string
	modifyEventAnalyzer func(string, string, string, string, string)
	deleteEventAnalyzer func(string, string, string, string, string)
)

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	CloudType = openstackinit.ReadConfig()

	if CloudType != "OpenStack" {
		pluginName := CloudType+".so"
		plugIn, err := plugin.Open("plugin/" + pluginName)
		if err != nil {
			log.Fatalf("[ERROR] Cloud not load the plugin %s in plugin directory %v",pluginName,err)
		}

		var ok bool
		modifyEventAnalyzerSymbol, err := plugIn.Lookup("ModifyEventAnalyzer")
		modifyEventAnalyzer, ok = modifyEventAnalyzerSymbol.(func(string, string, string, string, string))
		if err != nil || !ok {
			log.Fatalf("[ERROR] Something went wrong while loading plugin %v", err)
		}

		deleteEventAnalyzerSymbol, err := plugIn.Lookup("DeleteEventAnalyzer")
		deleteEventAnalyzer, ok = deleteEventAnalyzerSymbol.(func(string, string, string, string, string))
		if err != nil || !ok {
			log.Fatalf("[ERROR] Something went wrong while loading plugin %v", err)
		}
		log.Printf("[INFO] %s plugin loaded succesfully", CloudType)
	}

	config := functions.LoadKubeConfig()
	dynamicClient, err := dynamic.NewForConfig(config)
	checkErr(err)

	resource := schema.GroupVersionResource{Version: "v1", Resource: "pods"}
	w, err := dynamicClient.Resource(resource).Namespace("").Watch(context.TODO(), metav1.ListOptions{})
	checkErr(err)

	defer w.Stop()
	wg.Add(1)

	log.Println("[INFO] K8s cluster auto scalar started")
	for {
		wCh := w.ResultChan()

		for event := range wCh {
			eventFilter(event, config)
		}
	}
}

func eventFilter(event watch.Event, config *rest.Config) {
	defer handlenodeadd.PanicRecovery()

	b, _ := json.Marshal(event)
	var EventList datastructures.Event
	err := json.Unmarshal(b, &EventList)
	if err != nil {
		panic("Need to write the code. Stop containing")
	}

	switch EventList.Type {
	case "MODIFIED":
		if CloudType == "OpenStack" {
			handlenodeadd.ModifyEventAnalyzer(EventList)
		} else {
			modifyEventAnalyzer(openstackinit.ProjectName, openstackinit.ClientSecret, openstackinit.ClientId, openstackinit.AWSRegion, openstackinit.AuthFile)
		}
	case "DELETED":
		if CloudType == "OpenStack" {
			handelnodedelete.DeleteEventAnalyzer(EventList, config)
		} else {
			deleteEventAnalyzer(openstackinit.ProjectName, openstackinit.ClientSecret, openstackinit.ClientId, openstackinit.AWSRegion, openstackinit.AuthFile)
		}
	}
}
