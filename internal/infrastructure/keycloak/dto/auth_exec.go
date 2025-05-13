package kcdto

// KcAuthExecInfo: AuthenticationExecutionInfoRepresentation
type KcAuthExecInfo struct {
	ID                 string `json:"id"`
	Alias              string `json:"alias"`
	DisplayName        string `json:"displayName"`
	Level              int    `json:"level"`
	Priority           int    `json:"priority"`
	FlowID             string `json:"flowID"`
	AuthenticationFlow bool   `json:"authenticationFlow"`
	ProviderID         string `json:"providerID"`
}
