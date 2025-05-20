package kcdto

type KcClient struct {
	ID       string `json:"id"`
	ClientID string `json:"clientId"`
	Name     string `json:"name"`
	RootUrl  string `json:"rootUrl"`
	Enabled  bool   `json:"enabled"`
}
