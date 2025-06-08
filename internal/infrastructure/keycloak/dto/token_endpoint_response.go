package kcdto

// TokenEndpointResp represents the response returned from the token endpoint
// It contains various fields related to the access token, refresh token, and
// other related information for the authentication session.
type TokenEndpointResp struct {
	// TokenType indicates the type of the token (e.g., "Bearer").
	TokenType string `json:"token_type"`
	// AccessToken is the actual token used to access protected resources.
	AccessToken string `json:"access_token"`
	// ExpiresIn is the lifetime in seconds of the access token.
	ExpiresIn int `json:"expires_in"`
	// RefreshToken is used to obtain a new access token when the current
	// one expires.
	RefreshToken string `json:"refresh_token"`
	// RefreshExpiresIn indicates the lifetime in seconds of the refresh token.
	RefreshExpiresIn int `json:"refresh_expires_in"`
	// IDToken is the ID token that contains identity information about
	// the authenticated user.
	IDToken string `json:"id_token"`
	// NotBeforePolicy indicates the time before which the token should not be
	// accepted for validation.
	NotBeforePolicy int `json:"not-before-policy"` //nolint:tagliatelle
	// Scope defines the permissions granted by the access token.
	Scope string `json:"scope"`
	// SessionState provides the state of the session, typically used for
	// tracking the session on the server.
	SessionState string `json:"session_state"`
}
