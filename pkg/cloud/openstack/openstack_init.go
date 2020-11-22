package openstackinit

import (
	"github.com/Chathuru/kubernetes-cluster-autoscaler/pkg/common/datastructures"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"time"
)

// FlavorsList, and other list og global variables
var (
	FlavorsList         datastructures.FlavorList
	NetworkUUID         string
	CoolDownTime        time.Duration
	IgnoreNamespaceList map[string]bool
	MinNodeCount        int
	MaxNodeCount        int
	ImageName           string
	SecurityGroupName   string
	IdentityEndpoint    string
	Username            string
	Password            string
	TenantID            string
	DomainName          string
	ProjectName         string
	ClientSecret        string
	ClientID            string
	AWSRegion           string
	AuthFile            string
)

// ConfigYaml used to decode the configuration file
type ConfigYaml struct {
	CloudType          string            `yaml:"CloudType"`
	AuthOptions        AuthOptions       `yaml:"AuthOptions"`
	Network            Network           `yaml:"Network"`
	WorkerImageName    string            `yaml:"WorkerImageName"`
	CoolDownTime       int               `yaml:"CoolDownTime"`
	MinNodeCount       int               `yaml:"MinNodeCount"`
	MaxNodeCount       int               `yaml:"MaxNodeCount"`
	OpenStackFlavours  OpenStackFlavours `yaml:"OpenStackFlavours"`
	PassConfigToPlugin bool              `yaml:"PassConfigToPlugin"`
}

// AuthOptions list of credentials to authenticate cloud infrastructure
type AuthOptions struct {
	IdentityEndpoint string `yaml:"IdentityEndpoint"`
	Username         string `yaml:"Username"`
	Password         string `yaml:"Password"`
	TenantID         string `yaml:"TenantID"`
	DomainName       string `yaml:"DomainName"`
	ProjectName      string `yaml:"ProjectName"`
	ClientSecret     string `yaml:"ClientSecret"`
	ClientID         string `yaml:"ClientId"`
	AWSRegion        string `yaml:"AWSRegion"`
	AuthFile         string `yaml:"AuthFile"`
}

// Network OpenStack network configuration to used
// when creating worker nodes
type Network struct {
	SecurityGroupName string `yaml:"SecurityGroupName"`
	NetworkUUID       string `yaml:"NetworkUUID"`
}

// OpenStackFlavours user configured Open Stack Flavours in the config file.
type OpenStackFlavours struct {
	DefaultFlavour string     `yaml:"DefaultFlavour"`
	Flavours       []Flavours `yaml:"Flavours"`
}

// Flavours configured in config.yml
type Flavours struct {
	Name   string `yaml:"Name"`
	VCPU   int64  `yaml:"VCPU"`
	Memory int64  `yaml:"Memory"`
}

// ReadConfig read and configure starup variables from the config.yml
func ReadConfig() string {
	ConfigFile, err := ioutil.ReadFile("conf.yml")
	if err != nil {
		log.Fatalf("[ERROR] Error reading Config YAML file: %s\n", err)
	}

	conf := ConfigYaml{}
	err = yaml.Unmarshal(ConfigFile, &conf)
	if err != nil {
		log.Fatalf("[ERROR] Error decording Config YAML file: %s\n", err)
	}

	if conf.CloudType == "" {
		log.Fatal("[ERROR] \"CloudType\" must be set to one of OpenStack, GCP, AWS, libvirt, Other value.")
	}
	IdentityEndpoint = conf.AuthOptions.IdentityEndpoint
	Username = conf.AuthOptions.Username
	Password = conf.AuthOptions.Password
	TenantID = conf.AuthOptions.TenantID
	DomainName = conf.AuthOptions.DomainName
	if conf.CloudType == "OpenStack" && (IdentityEndpoint == "" || Username == "" || Password == "" || TenantID == "" || DomainName == "") {
		log.Fatal("[ERROR] Authentication details should not be empty.")
	}

	if conf.CloudType == "AWS" && conf.AuthOptions.AWSRegion == "" {
		log.Fatal("[ERROR] AWS Region should be a valid value")
	}

	if conf.CloudType == "GCP" && conf.AuthOptions.ProjectName == "" {
		log.Fatal("[ERROR] Project name should not be empty")
	}

	CoolDownTime = time.Duration(conf.CoolDownTime)
	MinNodeCount = conf.MinNodeCount
	MaxNodeCount = conf.MinNodeCount
	ImageName = conf.WorkerImageName
	SecurityGroupName = conf.Network.SecurityGroupName
	NetworkUUID = conf.Network.NetworkUUID

	FlavorDetails := []datastructures.FlavorDetails{}
	for _, Flavor := range conf.OpenStackFlavours.Flavours {
		FlavorDetails = append(FlavorDetails, datastructures.FlavorDetails{Flavor.Name, Flavor.VCPU, Flavor.Memory})
	}

	FlavorsList = datastructures.FlavorList{len(conf.OpenStackFlavours.Flavours), FlavorDetails, conf.OpenStackFlavours.DefaultFlavour}
	IgnoreNamespaceList = map[string]bool{"ingress-nginx": true, "kube-node-lease": true, "kube-public": true, "kube-system": true}

	return conf.CloudType
}

// GetOpenstackToken authenticate OpenStack cloud
func GetOpenstackToken() *gophercloud.ServiceClient {
	opts := gophercloud.AuthOptions{
		IdentityEndpoint: IdentityEndpoint,
		Username:         Username,
		Password:         Password,
		TenantID:         TenantID,
		DomainName:       DomainName,
	}

	provider, err := openstack.AuthenticatedClient(opts)
	if err != nil {
		panic(err)
	}
	client, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{Region: "LK"})
	if err != nil {
		panic(err)
	}

	return client
}
