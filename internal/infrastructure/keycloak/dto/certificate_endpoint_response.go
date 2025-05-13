package kcdto

// CertEndpointResp represents the response from a certificate endpoint
// It contains a list of public keys (PubKey) that are used to verify JWTs
// or perform other cryptographic operations.
type CertEndpointResp struct {
	// Keys is a list of JWKs in the form of PubKey
	// These keys are used by clients to validate the signature of JWTs or
	// other cryptographic tasks.
	Keys []PubKey `json:"keys"`
}

type PubKey struct {
	// The "kty" (Key Type) parameter identifies the cryptographic algorithm
	// family used with the key, such as "RSA" or "EC"
	Kty string `json:"kty"`
	// The "use" (public key use) parameter identifies the intended use of
	// the public key.  The "use" parameter is employed to indicate whether
	// a public key is used for encrypting data (enc) or verifying the
	// signature on data (sig).
	Use string `json:"use"`
	// The "alg" (algorithm) parameter identifies the algorithm intended for
	// use with the key.
	Alg string `json:"alg"`
	// The "kid" (key ID) parameter is used to match a specific key.  This
	// is used, for instance, to choose among a set of keys within a JWK Set
	// during key rollover.
	Kid string `json:"kid"`
	//  The "x5c" (X.509 certificate chain) parameter contains a chain of one
	// or more PKIX certificates. The certificate chain is represented as a
	// JSON array of certificate value strings. Each string in the array is
	// a base64-encoded DER PKIX certificate value.
	X5c []string `json:"x5c"`
	// The "x5t" (X.509 certificate SHA-1 thumbprint) parameter is a
	// base64url-encoded SHA-1 thumbprint (a.k.a. digest) of the DER
	// encoding of an X.509 certificate
	X5t string `json:"x5t"`
	// The "x5t#S256" (X.509 certificate SHA-256 thumbprint) parameter is a
	// base64url-encoded SHA-256 thumbprint (a.k.a. digest) of the DER
	// encoding of an X.509 certificate
	X5tS2566 string `json:"x5t#S256"`
}
