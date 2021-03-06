/*
Copyright © 2020 The k3d Author(s)

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package types

import (
	"fmt"
	"time"
)

// DefaultClusterName specifies the default name used for newly created clusters
const DefaultClusterName = "k3s-default"

// DefaultClusterNameMaxLength specifies the maximal length of a passed in cluster name
// This restriction allows us to construct an name consisting of
// <DefaultObjectNamePrefix[3]>-<ClusterName>-<TypeSuffix[5-10]>-<Counter[1-3]>
// ... and still stay within the 64 character limit (e.g. of docker)
const DefaultClusterNameMaxLength = 32

// DefaultK3sImageRepo specifies the default image repository for the used k3s image
const DefaultK3sImageRepo = "docker.io/rancher/k3s"

// DefaultLBImage defines the default cluster load balancer image
const DefaultLBImage = "docker.io/iwilltry42/k3d-proxy:v0.0.2"

// DefaultObjectNamePrefix defines the name prefix for every object created by k3d
const DefaultObjectNamePrefix = "k3d"

// ReadyLogMessageMaster defines the log messages we wait for until a master node is considered ready
var ReadyLogMessageByRole = map[Role]string{
	MasterRole:       "Wrote kubeconfig",
	WorkerRole:       "Successfully registered node",
	LoadBalancerRole: "start worker processes",
}

// Role defines a k3d node role
type Role string

// existing k3d node roles
const (
	MasterRole       Role = "master"
	WorkerRole       Role = "worker"
	NoRole           Role = "noRole"
	LoadBalancerRole Role = "loadbalancer"
)

// NodeRoles defines the roles available for nodes
var NodeRoles = map[string]Role{
	string(MasterRole):       MasterRole,
	string(WorkerRole):       WorkerRole,
	string(LoadBalancerRole): LoadBalancerRole,
}

// DefaultObjectLabels specifies a set of labels that will be attached to k3d objects by default
var DefaultObjectLabels = map[string]string{
	"app": "k3d",
}

// List of k3d technical label name
const (
	LabelToken           string = "k3d.cluster.token"
	LabelImageVolume     string = "k3d.cluster.imageVolume"
	LabelNetworkExternal string = "k3d.cluster.network.external"
	LabelNetwork         string = "k3d.cluster.network"
)

// DefaultRoleCmds maps the node roles to their respective default commands
var DefaultRoleCmds = map[Role][]string{
	MasterRole: {"server"},
	WorkerRole: {"agent"},
}

// DefaultTmpfsMounts specifies tmpfs mounts that are required for all k3d nodes
var DefaultTmpfsMounts = []string{
	"/run",
	"/var/run",
}

// DefaultNodeEnv defines some default environment variables that should be set on every node
var DefaultNodeEnv = []string{
	"K3S_KUBECONFIG_OUTPUT=/output/kubeconfig.yaml",
}

// DefaultToolsContainerImage defines the default image used for the tools container
const DefaultToolsContainerImage = "docker.io/iwilltry42/k3d-tools:v0.0.3" // TODO: get version dynamically or at build time

// DefaultImageVolumeMountPath defines the mount path inside k3d nodes where we will mount the shared image volume by default
const DefaultImageVolumeMountPath = "/k3d/images"

// DefaultConfigDirName defines the name of the config directory (where we'll e.g. put the kubeconfigs)
const DefaultConfigDirName = ".k3d" // should end up in $HOME/

// DefaultKubeconfigPrefix defines the default prefix for kubeconfig files
const DefaultKubeconfigPrefix = DefaultObjectNamePrefix + "-kubeconfig"

// DefaultAPIPort defines the default Kubernetes API Port
const DefaultAPIPort = "6443"

// DefaultAPIHost defines the default host (IP) for the Kubernetes API
const DefaultAPIHost = "0.0.0.0"

// DoNotCopyMasterFlags defines a list of commands/args that shouldn't be copied from an existing node when adding a similar node to a cluster
var DoNotCopyMasterFlags = []string{
	"--cluster-init",
}

// CreateClusterOpts describe a set of options one can set when creating a cluster
type CreateClusterOpts struct {
	DisableImageVolume  bool
	WaitForMaster       bool
	Timeout             time.Duration
	DisableLoadBalancer bool
	K3sServerArgs       []string
	K3sAgentArgs        []string
}

// StartClusterOpts describe a set of options one can set when (re-)starting a cluster
type StartClusterOpts struct {
	WaitForMaster bool
	Timeout       time.Duration
}

// CreateNodeOpts describes a set of options one can set when creating a new node
type CreateNodeOpts struct {
	Wait    bool
	Timeout time.Duration
}

// StartNodeOpts describes a set of options one can set when (re-)starting a node
type StartNodeOpts struct {
	Wait    bool
	Timeout time.Duration
}

// LoadImageOpts describes a set of options one can set for loading image(s) into cluster(s)
type LoadImageOpts struct {
	KeepTar bool
}

// ClusterNetwork describes a network which a cluster is running in
type ClusterNetwork struct {
	Name     string `yaml:"name" json:"name,omitempty"`
	External bool   `yaml:"external" json:"isExternal,omitempty"`
}

// Cluster describes a k3d cluster
type Cluster struct {
	Name               string             `yaml:"name" json:"name,omitempty"`
	Network            ClusterNetwork     `yaml:"network" json:"network,omitempty"`
	Token              string             `yaml:"cluster_token" json:"clusterToken,omitempty"`
	Nodes              []*Node            `yaml:"nodes" json:"nodes,omitempty"`
	InitNode           *Node              // init master node
	ExternalDatastore  ExternalDatastore  `yaml:"external_datastore" json:"externalDatastore,omitempty"`
	CreateClusterOpts  *CreateClusterOpts `yaml:"options" json:"options,omitempty"`
	ExposeAPI          ExposeAPI          `yaml:"expose_api" json:"exposeAPI,omitempty"`
	MasterLoadBalancer *Node              `yaml:"master_loadbalancer" json:"masterLoadBalancer,omitempty"`
	ImageVolume        string             `yaml:"image_volume" json:"imageVolume,omitempty"`
}

// MasterCount return number of master node into cluster
func (c *Cluster) MasterCount() int {
	masterCount := 0
	for _, node := range c.Nodes {
		if node.Role == MasterRole {
			masterCount++
		}
	}
	return masterCount
}

// WorkerCount return number of worker node into cluster
func (c *Cluster) WorkerCount() int {
	workerCount := 0
	for _, node := range c.Nodes {
		if node.Role == WorkerRole {
			workerCount++
		}
	}
	return workerCount
}

// Node describes a k3d node
type Node struct {
	Name       string            `yaml:"name" json:"name,omitempty"`
	Role       Role              `yaml:"role" json:"role,omitempty"`
	Image      string            `yaml:"image" json:"image,omitempty"`
	Volumes    []string          `yaml:"volumes" json:"volumes,omitempty"`
	Env        []string          `yaml:"env" json:"env,omitempty"`
	Cmd        []string          // filled automatically based on role
	Args       []string          `yaml:"extra_args" json:"extraArgs,omitempty"`
	Ports      []string          `yaml:"port_mappings" json:"portMappings,omitempty"`
	Restart    bool              `yaml:"restart" json:"restart,omitempty"`
	Labels     map[string]string // filled automatically
	Network    string            // filled automatically
	MasterOpts MasterOpts        `yaml:"master_opts" json:"masterOpts,omitempty"`
	WorkerOpts WorkerOpts        `yaml:"worker_opts" json:"workerOpts,omitempty"`
}

// MasterOpts describes some additional master role specific opts
type MasterOpts struct {
	IsInit    bool      `yaml:"is_initializing_master" json:"isInitializingMaster,omitempty"`
	ExposeAPI ExposeAPI // filled automatically
}

// ExternalDatastore describes an external datastore used for HA/multi-master clusters
type ExternalDatastore struct {
	Endpoint string `yaml:"endpoint" json:"endpoint,omitempty"`
	CAFile   string `yaml:"ca_file" json:"caFile,omitempty"`
	CertFile string `yaml:"cert_file" json:"certFile,omitempty"`
	KeyFile  string `yaml:"key_file" json:"keyFile,omitempty"`
	Network  string `yaml:"network" json:"network,omitempty"`
}

// ExposeAPI describes specs needed to expose the API-Server
type ExposeAPI struct {
	Host   string `yaml:"host" json:"host,omitempty"`
	HostIP string `yaml:"host_ip" json:"hostIP,omitempty"`
	Port   string `yaml:"port" json:"port"`
}

// WorkerOpts describes some additional worker role specific opts
type WorkerOpts struct{}

// GetDefaultObjectName prefixes the passed name with the default prefix
func GetDefaultObjectName(name string) string {
	return fmt.Sprintf("%s-%s", DefaultObjectNamePrefix, name)
}
