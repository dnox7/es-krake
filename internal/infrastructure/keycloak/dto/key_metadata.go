package kcdto

type KcKeysMetadata struct {
	Active map[string]interface{} `json:"active"`
	Keys   []KcKeyMetadata        `json:"keys"`
}

type KcKeyMetadata struct {
	ProviderID       string `json:"providerI"`
	ProviderPriority int64  `json:"providerPriority"`
	Kid              string `json:"kid"`
	Status           string `json:"status"`
	Type             string `json:"type"`
	Algo             string `json:"algorithm"`
	PublicKey        string `json:"publicKey"`
	Certificate      string `json:"certificate"`
	ValidTo          int64  `json:"validTo"`
}
