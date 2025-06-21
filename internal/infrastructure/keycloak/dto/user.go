package kcdto

//nolint:tagliatelle
type KcUser struct {
	ID                  string                `json:"id"`
	Username            string                `json:"username"`
	FirstName           string                `json:"firstName"`
	LastName            string                `json:"lastName"`
	Email               string                `json:"email"`
	EmailVerified       bool                  `json:"emailVerified"`
	Credentials         []KcCredential        `json:"credentials"`
	FederatedIdentities []KcFederatedIdentity `json:"federatedIdentities"`
	Enabled             bool                  `json:"enabled"`
}

type KcCredential struct {
	ID        string `json:"id"`
	Type      string `json:"type"`
	Value     string `json:"value"`
	Temporary bool   `json:"temporary"`
}

//nolint:tagliatelle
type KcFederatedIdentity struct {
	IdentityProvider string `json:"identityProvider"`
	UserID           string `json:"userId"`
	Username         string `json:"userName"`
}
