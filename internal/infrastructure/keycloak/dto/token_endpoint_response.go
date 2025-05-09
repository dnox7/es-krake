package kcdto

type TokenEndpointResponse struct {
	TokenType        string `json:"token_type"`
	AccessToken      string `json:"access_token"`
	ExpiresIn        int    `json:"expires_in"`
	RefreshToken     string `json:"refresh_token"`
	RefreshExpiresIn int    `json:"refresh_expires_in"`
	IDToken          string `json:"id_token"`
	NotBeforePolicy  int    `json:"not-before-policy"`
	Scope            string `json:"scope"`
	SessionState     string `json:"session_state"`
}
