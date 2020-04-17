#Redis LCM

```gotemplate
type ReconcileCache struct {
	Client   client.Client
	Recorder record.EventRecorder
	Log      logr.Logger
	CTX      context.Context
	Request  controller.Request
	DClient  dynamic.Interface
	Harbor   *harborCluster.HarborCluster
	Scheme   *runtime.Scheme
	ExpectCR *unstructured.Unstructured
	ActualCR *unstructured.Unstructured
	Labels   map[string]string
}

type CRLCMController interface {
	Reconcile() (*CRStatus, *reconciler.Results)
}

type CRStatus struct {
	Phase           Phase  `json:"phase,omitempty"`
	ExternalService string `json:"service,omitempty"`
	AvailableNodes  int32  `json:"availableNodes,omitempty"`
}

type Phase string

const (
	PendingPhase Phase = "Pending"
	ReadyPhase   Phase = "Ready"
	UpgradingPhase Phase = "Upgrading"
	DestroyingPhase Phase = "Destroying"
)
```

##Phase

* Pending, The reconcile has created CR , but the CR has not been scheduled yet.
* Ready,The CR has been deployed, the number of nodes is as expected, and the service is available.
* Upgrading, The CR is processing upscale、downscale、rollingupdate.
* Destroying, The CR is destroying.

##ExternalService
* connection access endpoint, expose the connection information of each component to HarborCluster Reconcile, reconcile can tell the harbor operator how to access redis, pg, etc.