package kcdto

// KcAuthFlow: AuthenticationFlowRepresentation
type KcAuthFlow struct {
	ID         string `json:"id"`
	Alias      string `json:"alias"`
	Desc       string `json:"description"`
	ProviderID string `json:"providerID"`
	TopLevel   bool   `json:"topLevel"`
	BuiltIn    bool   `json:"builtIn"`
	AuthExecs  []KcAuthExecExport
}

// KcAuthExecExport: AuthenticationExecutionExportRepresentation
type KcAuthExecExport struct {
	FlowAlias         string `json:"flowAlias"`
	AuthenticatorFlow string `json:"authenticatorFlow"`
	Priority          string `json:"priority"`
}
