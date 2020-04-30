package lcm

const (
	//ProperConn represents the connection info of the component.
	ProperConn = "Connection"
	//ProperPort represents the connection port of the component.
	ProperPort = "Port"
	//ProperUser represents the connection user of the component.
	ProperUser = "Username"
	//ProperPass represents the connection password of the component.
	ProperPass = "Password"
	//ProperNodes represent the available nodes of the component.
	ProperNodes = "AvailableNodes"
)

//Property is the current property of component.
type Property struct {
	//Property name, e.p: Connection,Port.
	Name string
	//Property value, e.p: "rfs-harborcluster-sample.svc"
	Value interface{}
}

type Properties []*Property

//New returns a new properties
func (ps Properties) New(Name string, Value interface{}) *Property {
	p := &Property{
		Name:  Name,
		Value: Value,
	}

	ps = append(ps, p)

	return p
}

//Update updates properties according to the given arguments
func (ps Properties) Update(Name string, Value interface{}) {
	for _, p := range ps {
		if p.Name == Name {
			p.Value = Value
			return
		}
	}
}

//ToInt parse properties value to int type
func (p *Property) ToInt() int {
	if p.Value != nil {
		if v, ok := p.Value.(int); ok {
			return v
		}
	}

	return 0
}

//ToString parse properties value to string type
func (p *Property) ToString() string {
	if p.Value != nil {
		if v, ok := p.Value.(string); ok {
			return v
		}
	}

	return ""
}

//ToFloat64 parse properties value to float64 type
func (p *Property) ToFloat64() float64 {
	if p.Value != nil {
		if v, ok := p.Value.(float64); ok {
			return v
		}
	}

	return 0
}
