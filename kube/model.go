package kube

import "time"

type GetPods struct {
	APIVersion string   `json:"apiVersion"`
	Items      []Items  `json:"items"`
	Kind       string   `json:"kind"`
	Metadata   Metadata `json:"metadata"`
}
type Annotations struct {
	KubernetesIoPsp    string `json:"kubernetes.io/psp"`
	PrometheusIoPath   string `json:"prometheus.io/path"`
	PrometheusIoPort   string `json:"prometheus.io/port"`
	PrometheusIoScrape string `json:"prometheus.io/scrape"`
}
type Labels struct {
	App             string `json:"app"`
	PodTemplateHash string `json:"pod-template-hash"`
	Team            string `json:"team"`
}
type OwnerReferences struct {
	APIVersion         string `json:"apiVersion"`
	BlockOwnerDeletion bool   `json:"blockOwnerDeletion"`
	Controller         bool   `json:"controller"`
	Kind               string `json:"kind"`
	Name               string `json:"name"`
	UID                string `json:"uid"`
}
type Metadata struct {
	Annotations       Annotations       `json:"annotations"`
	CreationTimestamp time.Time         `json:"creationTimestamp"`
	GenerateName      string            `json:"generateName"`
	Labels            Labels            `json:"labels"`
	Name              string            `json:"name"`
	Namespace         string            `json:"namespace"`
	OwnerReferences   []OwnerReferences `json:"ownerReferences"`
	ResourceVersion   string            `json:"resourceVersion"`
	SelfLink          string            `json:"selfLink"`
	UID               string            `json:"uid"`
}

type Env struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
type Exec struct {
	Command []string `json:"command"`
}
type PreStop struct {
	Exec Exec `json:"exec"`
}
type Lifecycle struct {
	PreStop PreStop `json:"preStop"`
}
type HTTPGet struct {
	Path   string      `json:"path"`
	Port   interface{} `json:"port"`
	Scheme string      `json:"scheme"`
}
type LivenessProbe struct {
	FailureThreshold    int     `json:"failureThreshold"`
	HTTPGet             HTTPGet `json:"httpGet"`
	InitialDelaySeconds int     `json:"initialDelaySeconds"`
	PeriodSeconds       int     `json:"periodSeconds"`
	SuccessThreshold    int     `json:"successThreshold"`
	TimeoutSeconds      int     `json:"timeoutSeconds"`
}
type Ports struct {
	ContainerPort int    `json:"containerPort"`
	Name          string `json:"name"`
	Protocol      string `json:"protocol"`
}
type ReadinessProbe struct {
	FailureThreshold    int     `json:"failureThreshold"`
	HTTPGet             HTTPGet `json:"httpGet"`
	InitialDelaySeconds int     `json:"initialDelaySeconds"`
	PeriodSeconds       int     `json:"periodSeconds"`
	SuccessThreshold    int     `json:"successThreshold"`
	TimeoutSeconds      int     `json:"timeoutSeconds"`
}
type Limits struct {
	CPU    string `json:"cpu"`
	Memory string `json:"memory"`
}
type Requests struct {
	CPU    string `json:"cpu"`
	Memory string `json:"memory"`
}

type SecurityContext struct {
	AllowPrivilegeEscalation bool `json:"allowPrivilegeEscalation"`
}
type VolumeMounts struct {
	MountPath string `json:"mountPath"`
	Name      string `json:"name"`
	ReadOnly  bool   `json:"readOnly,omitempty"`
	SubPath   string `json:"subPath,omitempty"`
}
type Containers struct {
	Env                      []Env           `json:"env"`
	Image                    string          `json:"image"`
	ImagePullPolicy          string          `json:"imagePullPolicy"`
	Lifecycle                Lifecycle       `json:"lifecycle"`
	LivenessProbe            LivenessProbe   `json:"livenessProbe"`
	Name                     string          `json:"name"`
	Ports                    []Ports         `json:"ports"`
	ReadinessProbe           ReadinessProbe  `json:"readinessProbe"`
	Resources                Resources       `json:"resources"`
	SecurityContext          SecurityContext `json:"securityContext"`
	TerminationMessagePath   string          `json:"terminationMessagePath"`
	TerminationMessagePolicy string          `json:"terminationMessagePolicy"`
	VolumeMounts             []VolumeMounts  `json:"volumeMounts"`
}
type ImagePullSecrets struct {
	Name string `json:"name"`
}
type Resources struct {
	Limits   Limits   `json:"limits"`
	Requests Requests `json:"requests"`
}

type InitContainers struct {
	Args                     []string        `json:"args"`
	Env                      []Env           `json:"env"`
	Image                    string          `json:"image"`
	ImagePullPolicy          string          `json:"imagePullPolicy"`
	Name                     string          `json:"name"`
	Resources                Resources       `json:"resources"`
	SecurityContext          SecurityContext `json:"securityContext"`
	TerminationMessagePath   string          `json:"terminationMessagePath"`
	TerminationMessagePolicy string          `json:"terminationMessagePolicy"`
	VolumeMounts             []VolumeMounts  `json:"volumeMounts"`
}
type Tolerations struct {
	Effect            string `json:"effect"`
	Key               string `json:"key"`
	Operator          string `json:"operator"`
	TolerationSeconds int    `json:"tolerationSeconds"`
}
type ConfigMap struct {
	DefaultMode int    `json:"defaultMode"`
	Name        string `json:"name"`
}
type EmptyDir struct {
	Medium string `json:"medium"`
}
type Secret struct {
	DefaultMode int    `json:"defaultMode"`
	SecretName  string `json:"secretName"`
}
type Volumes struct {
	ConfigMap ConfigMap `json:"configMap,omitempty"`
	Name      string    `json:"name"`
	EmptyDir  EmptyDir  `json:"emptyDir,omitempty"`
	Secret    Secret    `json:"secret,omitempty"`
}
type Spec struct {
	Containers                    []Containers       `json:"containers"`
	DNSPolicy                     string             `json:"dnsPolicy"`
	EnableServiceLinks            bool               `json:"enableServiceLinks"`
	ImagePullSecrets              []ImagePullSecrets `json:"imagePullSecrets"`
	InitContainers                []InitContainers   `json:"initContainers"`
	NodeName                      string             `json:"nodeName"`
	Priority                      int                `json:"priority"`
	RestartPolicy                 string             `json:"restartPolicy"`
	SchedulerName                 string             `json:"schedulerName"`
	SecurityContext               SecurityContext    `json:"securityContext"`
	ServiceAccount                string             `json:"serviceAccount"`
	ServiceAccountName            string             `json:"serviceAccountName"`
	TerminationGracePeriodSeconds int                `json:"terminationGracePeriodSeconds"`
	Tolerations                   []Tolerations      `json:"tolerations"`
	Volumes                       []Volumes          `json:"volumes"`
}
type Conditions struct {
	LastProbeTime      interface{} `json:"lastProbeTime"`
	LastTransitionTime time.Time   `json:"lastTransitionTime"`
	Status             string      `json:"status"`
	Type               string      `json:"type"`
}
type Terminated struct {
	ContainerID string    `json:"containerID"`
	ExitCode    int       `json:"exitCode"`
	FinishedAt  time.Time `json:"finishedAt"`
	Reason      string    `json:"reason"`
	StartedAt   time.Time `json:"startedAt"`
}
type LastState struct {
	Terminated Terminated `json:"terminated"`
}
type Running struct {
	StartedAt time.Time `json:"startedAt"`
}
type State struct {
	Running    Running    `json:"running"`
	Terminated Terminated `json:"terminated"`
}
type ContainerStatuses struct {
	ContainerID  string    `json:"containerID"`
	Image        string    `json:"image"`
	ImageID      string    `json:"imageID"`
	LastState    LastState `json:"lastState"`
	Name         string    `json:"name"`
	Ready        bool      `json:"ready"`
	RestartCount int       `json:"restartCount"`
	State        State     `json:"state"`
}
type InitContainerStatuses struct {
	ContainerID  string    `json:"containerID"`
	Image        string    `json:"image"`
	ImageID      string    `json:"imageID"`
	LastState    LastState `json:"lastState"`
	Name         string    `json:"name"`
	Ready        bool      `json:"ready"`
	RestartCount int       `json:"restartCount"`
	State        State     `json:"state"`
}
type Status struct {
	Conditions            []Conditions            `json:"conditions"`
	ContainerStatuses     []ContainerStatuses     `json:"containerStatuses"`
	HostIP                string                  `json:"hostIP"`
	InitContainerStatuses []InitContainerStatuses `json:"initContainerStatuses"`
	Phase                 string                  `json:"phase"`
	PodIP                 string                  `json:"podIP"`
	QosClass              string                  `json:"qosClass"`
	StartTime             time.Time               `json:"startTime"`
}
type Items struct {
	APIVersion string   `json:"apiVersion"`
	Kind       string   `json:"kind"`
	Metadata   Metadata `json:"metadata"`
	Spec       Spec     `json:"spec"`
	Status     Status   `json:"status"`
}
