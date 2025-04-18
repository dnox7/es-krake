package dto

type KeycloakAuthenticationFlow struct {
	ID                       string
	Alias                    string
	AuthenticationExecutions []KeycloakAuthenticationExecution
}

type KeycloakAuthenticationExecution struct {
	AuthenticatorFlow bool   `json:"authenticatorFlow"`
	FlowAlias         string `json:"FlowAlias"`
}
