package kcdto

type BruteForceStatus struct {
	Disabled      bool   `json:"disabled"`
	LastFailure   int    `json:"lastFailure"`
	LastIPFailure string `json:"lastIPFailure"`
	NumFailures   int    `json:"numFailures"`
}
