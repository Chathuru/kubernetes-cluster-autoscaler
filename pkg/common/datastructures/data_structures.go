package datastructures

import v1 "k8s.io/api/core/v1"

// Event Type for Kubernetes Event decoder struct
type Event struct {
	Type   string `json:"Type"`
	Object v1.Pod `json:"Object"`
}

// FlavorList type is the list of OpenStack flavors
// to use when creating Kubernetes worker node
type FlavorList struct {
	FlavorNum     int
	Flavor        []FlavorDetails
	FlavorDefault string
}

// FlavorDetails is the OpenStack flavor details
type FlavorDetails struct {
	Name           string
	RequestsCPU    int64
	RequestsMemory int64
}
