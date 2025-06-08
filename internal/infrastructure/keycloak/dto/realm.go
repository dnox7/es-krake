package kcdto

//nolint:tagliatelle
type KcRealm struct {
	ID                  string `json:"id"`
	Realm               string `json:"realm"`
	DisplayName         string `json:"displayName"`
	Enabled             bool   `json:"enabled"`
	EditUsernameAllowed bool   `json:"editUsernameAllowed"`
}
