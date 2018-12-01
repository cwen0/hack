package types

// NetworkConfig is network config
type NetworkConfig struct {
	Ingress []string `json:"ingress"`
	Egress  []string `json:"egress"`
}

// FailpointConfig is failpoint config
type FailpointConfig struct {
	Path  string `json:"path"`
	Value string `json:"value"`
}

// FailpointFe is failpoint fe
type FailpointFe struct {
	Type string `json:"type"`
}