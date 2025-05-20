package kcdto

type KcIdentityProvider struct {
	Alias       string `json:"alias"`
	DisplayName string `json:"displayName"`
	InternalID  string `json:"internalId"`
	ProviderID  string `json:"providerId"`
	Enabled     bool   `json:"enabled"`
}
