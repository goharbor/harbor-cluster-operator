package status

type CRStatus struct {
	Phase      Phase       `json:"phase"`
	Properties []*Property `json:"properties"`
}

type Phase string

const (
	PendingPhase    Phase = "Pending"
	DeployingPhase  Phase = "Deploying"
	ReadyPhase      Phase = "Ready"
	UpgradingPhase  Phase = "Upgrading"
	DestroyingPhase Phase = "Destroying"
)

const (
	ProperConn  = "Connection"
	ProperPort  = "Port"
	ProperUser  = "Username"
	ProperPass  = "Password"
	ProperNodes = "AvailableNodes"
)

type Property struct {
	Name  string
	Type  string
	Value interface{}
}
