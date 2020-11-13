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
	"sync"
)

var (
	wg sync.WaitGroup
)

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	openstackinit.ReadConfig()

	config := functions.LoadKubeConfig()
	dynamicClient, err := dynamic.NewForConfig(config)
	checkErr(err)

	resource := schema.GroupVersionResource{Version: "v1", Resource: "pods"}
	w, err := dynamicClient.Resource(resource).Namespace("").Watch(context.TODO(), metav1.ListOptions{})
	checkErr(err)

	defer w.Stop()
	wg.Add(1)

	log.Println("[INFO] K8s cluster auto scaler started")
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
		panic("Need to write the code. Stop contineing")
	}

	switch EventList.Type {
	case "MODIFIED":
		handlenodeadd.ModifyEventAnalyzer(EventList)
	case "DELETED":
		handelnodedelete.DeleteEventAnalyzer(EventList, config)
	}
}
