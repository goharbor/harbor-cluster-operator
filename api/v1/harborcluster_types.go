/*

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1

import (
	cmmeta "github.com/jetstack/cert-manager/pkg/apis/meta/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// the name of component used harbor cluster.
type Component string

var (
	HarborClusterGVK = schema.GroupVersionKind{
		Group:   GroupVersion.Group,
		Version: GroupVersion.Version,
		Kind:    "HarborCluster",
	}
)

// all Component used in harbor cluster full stack.
const (
	ComponentHarbor   Component = "harbor"
	ComponentRedis    Component = "redis"
	ComponentStorage  Component = "storage"
	ComponentDatabase Component = "database"
)

// HarborClusterSpec defines the desired state of HarborCluster
type HarborClusterSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// harbor version to be deployed, this version determines the image tags of harbor service components
	// +kubebuilder:validation:Required
	// https://semver.org/#is-there-a-suggested-regular-expression-regex-to-check-a-semver-string
	// +kubebuilder:validation:Pattern="^(?P<major>0|[1-9]\\d*)\\.(?P<minor>0|[1-9]\\d*)\\.(?P<patch>0|[1-9]\\d*)(?:-(?P<prerelease>(?:0|[1-9]\\d*|\\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\\.(?:0|[1-9]\\d*|\\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\\+(?P<buildmetadata>[0-9a-zA-Z-]+(?:\\.[0-9a-zA-Z-]+)*))?$"
	Version string `json:"version"`

	// The url exposed to clients to access harbor
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern="^https?://.*$"
	PublicURL string `json:"publicURL"`

	// Password for the root admin
	// +kubebuilder:validation:Required
	AdminPasswordSecret string `json:"adminPasswordSecret"`

	// Secret reference for the TLS certs
	// +optional
	TLSSecret string `json:"tlsSecret,omitempty"`

	// The issuer for Harbor certificates.
	// If the 'kind' field is not set, or set to 'Issuer', an Issuer resource
	// with the given name in the same namespace as the Certificate will be used.
	// If the 'kind' field is set to 'ClusterIssuer', a ClusterIssuer with the
	// provided name will be used.
	// The 'name' field in this stanza is required at all times.
	CertificateIssuerRef cmmeta.ObjectReference `json:"certificateIssuerRef,omitempty"`

	// Indicates that the harbor is paused.
	// +optional
	Paused bool `json:"paused,omitempty"`

	// The Maximum priority. Deployments may be created with priority in interval ] priority - 100 ; priority ]
	// +kubebuilder:validation:Optional
	Priority *int32 `json:"priority,omitempty"`

	// Pod instance number
	// +kubebuilder:validation:Required
	Replicas int `json:"replicas"`

	// Source registry of images, the default is dockerhub
	// +kubebuilder:default=&ImageSource{registry: "docker.io";}
	ImageSource *ImageSource `json:"imageSource,omitempty"`

	// Extra configuration options for jobservices
	// +optional
	JobService *JobService `json:"jobService,omitempty"`

	// Extra configuration options for clair scanner
	// +optional
	Clair *Clair `json:"clair,omitempty"`

	// Extra configuration options for trivy scanner
	// +kubebuilder:validation:Optional
	Trivy *Trivy `json:"trivy,omitempty"`

	// Extra configuration options for chartmeseum
	// +kubebuilder:validation:Optional
	ChartMuseum *ChartMuseum `json:"chartMuseum,omitempty"`

	// Cache service(Redis) configurations might be external redis services or inCluster redis services
	// +kubebuilder:validation:Required
	Redis *Redis `json:"redis"`

	// database service (PostgresSQL) configuration
	// +kubebuilder:validation:Required
	Database *Database `json:"database"`

	// Storage service configurations. Might be external cloud storage services or inCluster storage (minIO)
	// +kubebuilder:validation:Required
	Stroage *Storage `json:"storage"`
}

// +kubebuilder:validation:Enum=Lion;Wolf;Dragon
type Storage struct {
	// set the kind of which storage service to be used. Set the kind as "azure", "gcs", "s3", "oss", "swift" or "inCluster", and fill the information.
	// in the options section. inCluster indicates the local storage service of harbor-cluster. We use minIO as a default built-in object storage service.
	// +kubebuilder:validation:Enum=inCLuster;azure;gcs;s3;oss;swift
	Kind string `json:"kind"`

	// inCLuster options.
	InCluster *InCluster `json:"options,omitempty"`

	// Azure options.
	Azure *Azure `json:"azure,omitempty"`

	// Gcs options.
	Gcs *Gcs `json:"gcs,omitempty"`

	// S3 options.
	S3 *S3 `json:"s3,omitempty"`

	// Swift options.
	Swift *Swift `json:"swift,omitempty"`

	// Oss options.
	Oss *Oss `json:"oss,omitempty"`
}

type Oss struct {
	// +kubebuilder:validation:Required
	Accesskeyid string `json:"accesskeyid"`
	// +kubebuilder:validation:Required
	Accesskeysecret string `json:"accesskeysecret"`
	// +kubebuilder:validation:Required
	Region string `json:"region"`
	// +kubebuilder:validation:Required
	Bucket string `json:"bucket"`
	// +kubebuilder:validation:Required
	Endpoint      string `json:"endpoint"`
	Internal      string `json:"internal,omitempty"`
	Encrypt       string `json:"encrypt,omitempty"`
	Secure        string `json:"secure,omitempty"`
	Chunksize     string `json:"chunksize,omitempty"`
	Rootdirectory string `json:"rootdirectory,omitempty"`
}

type Swift struct {
	// +kubebuilder:validation:Required
	Authurl string `json:"authurl"`
	// +kubebuilder:validation:Required
	Username string `json:"username"`
	// +kubebuilder:validation:Required
	Password string `json:"password"`
	// +kubebuilder:validation:Required
	Container string `json:"container"`
	// +kubebuilder:validation:Required
	Region string `json:"region"`
	// +kubebuilder:validation:Required
	Tenant              string   `json:"tenant"`
	Tenantid            string   `json:"tenantid,omitempty"`
	Domain              string   `json:"domain,omitempty"`
	Domainid            string   `json:"domainid,omitempty"`
	Trustid             string   `json:"trustid,omitempty"`
	Insecureskipverify  bool     `json:"insecureskipverify,omitempty"`
	Chunksize           string   `json:"chunksize,omitempty"`
	Prefix              string   `json:"prefix,omitempty"`
	Secretkey           string   `json:"secretkey,omitempty"`
	Authversion         int      `json:"authversion,omitempty"`
	Endpointtype        string   `json:"endpointtype,omitempty"`
	Tempurlcontainerkey bool     `json:"tempurlcontainerkey,omitempty"`
	Tempurlmethods      []string `json:"tempurlmethods,omitempty"`
}

type S3 struct {
	// +kubebuilder:validation:Required
	Region string `json:"region"`
	// +kubebuilder:validation:Required
	Bucket string `json:"bucket"`
	// +kubebuilder:validation:Required
	Accesskey string `json:"accesskey"`
	// +kubebuilder:validation:Required
	Secretkey string `json:"secretkey"`
	// +kubebuilder:validation:Required
	Regionendpoint string `json:"regionendpoint"`
	Encrypt        bool   `json:"encrypt,omitempty"`
	Keyid          string `json:"keyid,omitempty"`
	Secure         bool   `json:"secure,omitempty"`
	V4auth         bool   `json:"v4auth,omitempty"`
	Chunksize      string `json:"chunksize,omitempty"`
	Rootdirectory  string `json:"rootdirectory,omitempty"`
	Storageclass   string `json:"storageclass,omitempty"`
}

type Gcs struct {
	// +kubebuilder:validation:Required
	Bucket string `json:"bucket"`
	// The base64 encoded json file which contains the key
	Encodedkey string `json:"encodedkey"`
	// +kubebuilder:validation:Required
	Rootdirectory string `json:"rootdirectory"`
	Chunksize     string `json:"chunksize,omitempty"`
}

type Azure struct {
	// +kubebuilder:validation:Required
	Accountname string `json:"accountname"`
	// +kubebuilder:validation:Required
	Accountkey string `json:"accountkey"`
	// +kubebuilder:validation:Required
	Container string `json:"container"`
	Realm     string `json:"realm,omitempty"`
}

// InCluster of storage.
type InCluster struct {
	// inCluster Provider, just support minIO now.
	Provider string     `json:"provider,omitempty"`
	Spec     *MinIOSpec `json:"spec,omitempty"`
}

type MinIOSpec struct {
	// Supply number of replicas.
	// For standalone mode, supply 1. For distributed mode, supply 4 or more (should be even).
	// Note that the operator does not support upgrading from standalone to distributed mode.
	// +kubebuilder:validation:Required
	Replicas int32 `json:"replicas"`
	// Version defines the MinIO Client (mc) Docker image version.
	Version string `json:"version,omitempty"`
	// VolumeClaimTemplate allows a user to specify how volumes inside a MinIOInstance
	// +optional
	VolumeClaimTemplate corev1.PersistentVolumeClaim `json:"volumeClaimTemplate,omitempty"`
	// If provided, use these requests and limit for cpu/memory resource allocation
	// +optional
	Resources corev1.ResourceRequirements `json:"resources,omitempty"`
}

type PostgresSQL struct {
	Storage          string              `json:"storage,omitempty"`
	Replicas         int                 `json:"replicas,omitempty"`
	Version          string              `json:"version,omitempty"`
	StorageClassName string              `json:"storageClassName,omitempty"`
	Resources        corev1.ResourceList `json:"resources,omitempty"`

	// External params following.
	// The secret must contains "address:port","usernane" and "password".
	SecretName     string `json:"secretName,omitempty"`
	SslConfig      string `json:"sslConfig,omitempty"`
	ConnectTimeout int    `json:"connectTimeout,omitempty"`
}

type Database struct {
	// Set the kind of which redis service to be used, inCluster or external.
	// +kubebuilder:validation:Enum=inCluster;external
	Kind string `json:"kind"`

	PostgresSQL *PostgresSQL `json:"spec,omitempty"`
}

type Redis struct {
	// Set the kind of which redis service to be used, inCluster or external. Setting up a harbor-cluster with external redis service should provide client params to communicate. The difference between inCluster redis and external redis is that the inCluster redis installed automatically.
	// +kubebuilder:validation:Enum=inCluster;external
	Kind string `json:"kind"`

	// +kubebuilder:validation:Required
	Spec *RedisSpec `json:"spec"`
}

type RedisSpec struct {
	Server   *RedisServer `json:"server,omitempty"`
	Sentinel *Sentinel    `json:"sentinel,omitempty"`

	// External params following.
	// The secret must contains "address:port","usernane" and "password".
	SecretName string `json:"secretName,omitempty"`
	// Maximum number of socket connections.
	// Default is 10 connections per every CPU as reported by runtime.NumCPU.
	PoolSize int `json:"poolSize,omitempty"`
	// TLS Config to use. When set TLS will be negotiated.
	// set the secret which type of Opaque, and contains "tls.key","tls.crt","ca.crt".
	TlsConfig string `json:"tlsConfig,omitempty"`
	GroupName string `json:"groupName,omitempty"`
	// +kubebuilder:validation:Enum=sentinel;redis
	Schema string  `json:"schema,omitempty"`
	Hosts  []Hosts `json:"hosts,omitempty"`
}

type Hosts struct {
	Host string `json:"host,omitempty"`
	Port string `json:"port,omitempty"`
}

type Sentinel struct {
	Replicas int `json:"replicas,omitempty"`
}

type RedisServer struct {
	Replicas         int                 `json:"replicas,omitempty"`
	Resources        corev1.ResourceList `json:"resources,omitempty"`
	StorageClassName string              `json:"storageClassName,omitempty"`
	// the size of storage used in redis.
	Storage string `json:"storage,omitempty"`
}

type ChartMuseum struct {
	AbsoluteURL bool `json:"absoluteURL,omitempty"`
}

type Trivy struct {
	GithubToken string `json:"githubToken,omitempty"`
}

type ImageSource struct {
	Registry        string `json:"registry,omitempty"`
	ImagePullSecret string `json:"imagePullSecret,omitempty"`
}

type Clair struct {
	UpdateInterval       int      `json:"updateInterval,omitempty"`
	VulnerabilitySources []string `json:"vulnerabilitySources,omitempty"`
}

type JobService struct {
	// +kubebuilder:validation:Required
	Replicas string `json:"replicas"`

	// +optional
	WorkerCount int32 `json:"workerCount"`
}

// HarborClusterStatus defines the observed state of HarborCluster
type HarborClusterStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Conditions []HarborClusterCondition `json:"conditions,omitempty"`

	ComponentsStatus map[Component]*ComponentsStatus `json:"ComponentsStatus,omitempty"`
}

type ComponentsStatus struct {
	Replicas      int `json:"replicas,omitempty"`
	ReadyReplicas int `json:"readyReplicas,omitempty"`
}

// HarborClusterConditionType is a valid value for HarborClusterConditionType.Type
type HarborClusterConditionType string

// HarborClusterCondition contains details for the current condition of this pod.
type HarborClusterCondition struct {
	// Type is the type of the condition.
	Type HarborClusterConditionType `json:"type"`
	// Status is the status of the condition.
	// Can be True, False, Unknown.
	Status corev1.ConditionStatus `json:"status"`
	// Last time the condition transitioned from one status to another.
	// +optional
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty"`
	// Unique, one-word, CamelCase reason for the condition's last transition.
	// +optional
	Reason string `json:"reason,omitempty"`
	// Human-readable message indicating details about last transition.
	// +optional
	Message string `json:"message,omitempty"`
}

// +kubebuilder:object:root=true

// HarborCluster is the Schema for the harborclusters API
type HarborCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   HarborClusterSpec   `json:"spec,omitempty"`
	Status HarborClusterStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// HarborClusterList contains a list of HarborCluster
type HarborClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []HarborCluster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&HarborCluster{}, &HarborClusterList{})
}
