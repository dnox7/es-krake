package keycloak

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/dpe27/es-krake/internal/infrastructure/httpclient"
	kcdto "github.com/dpe27/es-krake/internal/infrastructure/keycloak/dto"
	"github.com/dpe27/es-krake/pkg/log"
	"github.com/dpe27/es-krake/pkg/utils"
)

var ErrPubKeyNotFound = errors.New("Public key not found")

const keyPath = "/keys"

type KcKeyService interface {
	// GetJWKSKeysFromEndpoint retrieves the current public keys from the OIDC JWKS endpoint.
	// These keys are used by clients to verify JWT signatures.
	GetJWKSKeysFromEndpoint(ctx context.Context, realm string) (kcdto.CertEndpointResp, error)
	// GetJWKSKeysFromCache retrieves the cached JWKS for the specified realm.
	// It returns the cached CertEndpointResp and a boolean indicating whether the cache was hit.
	// If the cache is not found or expired, the boolean will be false.
	GetJWKSKeysFromCache(realm string) (kcdto.CertEndpointResp, bool)
	// CacheJWKSKeyForRealm stores the provided JWKS (JSON Web Key Set) for the specified realm
	// into the local cache. This function allows preloading or updating the cached keys without
	// fetching them from the Keycloak server.
	CacheJWKSKeyForRealm(realm string, pubKeys kcdto.CertEndpointResp)
	// ExtractKey searches for a public key in the CertEndpointResp that matches the provided keyID.
	// It returns the first certificate (x5c[0]) associated with the matching key.
	// If no matching key is found or the key list is empty, it returns an error.
	ExtractKey(keyID string, resp kcdto.CertEndpointResp) (string, error)
	// GetRealmKeys retrieves all keys (including private and inactive ones) from the Keycloak Admin API.
	// This function requires admin privileges to access the endpoint.
	GetRealmKeys(ctx context.Context, realm, token string) (kcdto.KcKeysMetadata, error)
}

type keyService struct {
	BaseKcService
	logger *log.Logger
	cache  map[string]kcdto.CertEndpointResp
}

func NewKcKeyService(base BaseKcService) KcKeyService {
	return &keyService{
		logger:        log.With("service", "keycloak_key_service"),
		BaseKcService: base,
	}
}

// GetJWKSKeysFromEndpoint implements KcKeyService.
func (k *keyService) GetJWKSKeysFromEndpoint(
	ctx context.Context,
	realm string,
) (kcdto.CertEndpointResp, error) {
	url := k.PublicKeyUrl(realm)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		k.logger.Error(ctx, utils.ErrorCreateReq, "error", err.Error())
		return kcdto.CertEndpointResp{}, err
	}

	opts := httpclient.ReqOptBuidler().
		Log().LogResBody().LogReqBodyOnlyError().
		LoggedReqBody([]string{}).
		LoggedResBody([]string{}).
		Build()

	res, err := k.Client().Do(req, opts)
	if err != nil {
		return kcdto.CertEndpointResp{}, err
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			k.logger.Error(ctx, utils.ErrorCloseResponseBody, "error", err.Error())
		}
	}()

	if res.StatusCode != http.StatusOK {
		return kcdto.CertEndpointResp{}, fmt.Errorf("call api certificate endpoint status: %s", res.Status)
	}

	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return kcdto.CertEndpointResp{}, err
	}

	certResp := kcdto.CertEndpointResp{}
	err = json.Unmarshal(bytes, &certResp)
	if err != nil {
		return kcdto.CertEndpointResp{}, err
	}

	return certResp, nil
}

// GetJWKSKeysFromCache implements KcKeyService.
func (k *keyService) GetJWKSKeysFromCache(realm string) (kcdto.CertEndpointResp, bool) {
	pubKeys, ok := k.cache[realm]
	return pubKeys, ok
}

// CacheJWKSKeyForRealm implements KcKeyService.
func (k *keyService) CacheJWKSKeyForRealm(realm string, pubKeys kcdto.CertEndpointResp) {
	mu := &sync.Mutex{}
	mu.Lock()
	k.cache[realm] = pubKeys
	mu.Unlock()
}

// ExtractKey implements KcKeyService.
func (k *keyService) ExtractKey(keyID string, resp kcdto.CertEndpointResp) (string, error) {
	if len(resp.Keys) == 0 {
		return "", ErrPubKeyNotFound
	}
	for _, pk := range resp.Keys {
		if pk.Kid == keyID {
			return pk.X5c[0], nil
		}
	}
	return "", ErrPubKeyNotFound
}

// GetRealmKeys implements KcKeyService.
func (k *keyService) GetRealmKeys(
	ctx context.Context,
	realm string,
	token string,
) (kcdto.KcKeysMetadata, error) {
	url := k.AdminRealmUrl(realm) + keyPath
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		k.logger.Error(ctx, utils.ErrorCreateReq, "error", err.Error())
		return kcdto.KcKeysMetadata{}, err
	}

	req.Header.Add(httpclient.HeaderAuthorization, httpclient.AuthSchemeBearer+token)
	opts := httpclient.ReqOptBuidler().
		Log().LogResBody().LogReqBodyOnlyError().
		LoggedReqBody([]string{}).
		LoggedResBody([]string{}).
		Build()

	res, err := k.Client().Do(req, opts)
	if err != nil {
		return kcdto.KcKeysMetadata{}, err
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			k.logger.Error(ctx, utils.ErrorCloseResponseBody, "error", err.Error())
		}
	}()

	if res.StatusCode != http.StatusOK {
		return kcdto.KcKeysMetadata{}, fmt.Errorf("call api get key metadata status: %s", res.Status)
	}

	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return kcdto.KcKeysMetadata{}, err
	}

	keys := kcdto.KcKeysMetadata{}
	err = json.Unmarshal(bytes, &keys)
	if err != nil {
		return kcdto.KcKeysMetadata{}, err
	}

	return keys, nil
}
