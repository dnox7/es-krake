package keycloak

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	domainerr "github.com/dpe27/es-krake/internal/domain/shared/errors"
	"github.com/dpe27/es-krake/internal/infrastructure/httpclient"
	kcdto "github.com/dpe27/es-krake/internal/infrastructure/keycloak/dto"
	"github.com/dpe27/es-krake/pkg/log"
	"github.com/dpe27/es-krake/pkg/utils"
)

const identityProviderPath = "/identity-provider/instances"

type KcIdentityProviderService interface {
	GetIdPList(ctx context.Context, realm, token string) ([]kcdto.KcIdentityProvider, error)
	PostIdP(ctx context.Context, body map[string]interface{}, realm, token string) error
	PutIdP(ctx context.Context, body map[string]interface{}, realm, alias, token string) error
	DeleteIdP(ctx context.Context, realm, alias, token string) error
}

type identityProvider struct {
	BaseKcService
	logger *log.Logger
}

func NewKcIdentityProviderService(base BaseKcService) KcIdentityProviderService {
	return &identityProvider{
		BaseKcService: base,
		logger:        log.With("service", "keycloak_identity_provider_service"),
	}
}

// GetIdPList implements KcIdentityProviderService.
func (i *identityProvider) GetIdPList(
	ctx context.Context,
	realm string,
	token string,
) ([]kcdto.KcIdentityProvider, error) {
	url := i.AdminRealmUrl(realm) + identityProviderPath
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		i.logger.Error(ctx, utils.ErrorCreateReq, "error", err.Error())
		return nil, err
	}

	query := req.URL.Query()
	req.URL.RawQuery = query.Encode()
	req.Header.Add(httpclient.HeaderAuthorization, httpclient.AuthSchemeBearer+token)
	opts := httpclient.ReqOptBuidler().
		Log().
		LogReqBodyOnlyError().
		LoggedReqBody([]string{}).
		LogResBody().
		LoggedResBody([]string{}).
		Build()

	res, err := i.Client().Do(req, opts)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			i.logger.Error(ctx, utils.ErrorCloseResponseBody, "error", err.Error())
		}
	}()

	if res.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		var idpSlice []kcdto.KcIdentityProvider
		err = json.Unmarshal(bodyBytes, &idpSlice)
		if err != nil {
			return nil, err
		}
		return idpSlice, nil
	}

	if res.StatusCode == http.StatusNotFound {
		return nil, domainerr.ErrorNotFound
	}

	return nil, fmt.Errorf("call api get identity provider list status: %s", res.Status)
}

// PostIdP implements KcIdentityProviderService.
func (i *identityProvider) PostIdP(
	ctx context.Context,
	body map[string]interface{},
	realm string,
	token string,
) error {
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		i.logger.Error(ctx, utils.ErrorMarshalFailed, "error", err.Error())
		return err
	}

	url := i.AdminRealmUrl(realm) + identityProviderPath
	reqBody := strings.NewReader(string(bodyBytes))
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, reqBody)
	if err != nil {
		i.logger.Error(ctx, utils.ErrorCreateReq, "error", err.Error())
		return err
	}

	req.Header.Add(httpclient.HeaderContentType, httpclient.MIMEApplicationJSON)
	req.Header.Add(httpclient.HeaderAuthorization, httpclient.AuthSchemeBearer+token)
	opts := httpclient.ReqOptBuidler().
		Log().
		LogReqBodyOnlyError().
		LoggedReqBody([]string{}).
		LogResBody().
		LoggedResBody([]string{}).
		Build()

	res, err := i.Client().Do(req, opts)
	if err != nil {
		return err
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			i.logger.Error(ctx, utils.ErrorCloseResponseBody, "error", err.Error())
		}
	}()

	if res.StatusCode == http.StatusCreated {
		return nil
	}

	return fmt.Errorf("call api create identity provider status: %s", res.Status)
}

// PutIdP implements KcIdentityProviderService.
func (i *identityProvider) PutIdP(
	ctx context.Context,
	body map[string]interface{},
	realm string,
	alias string,
	token string,
) error {
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		i.logger.Error(ctx, utils.ErrorMarshalFailed, "error", err.Error())
		return err
	}

	reqBody := strings.NewReader(string(bodyBytes))
	url := i.AdminRealmUrl(realm) + identityProviderPath + "/" + alias
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, reqBody)
	if err != nil {
		i.logger.Error(ctx, utils.ErrorCreateReq, "error", err.Error())
		return err
	}

	req.Header.Add(httpclient.HeaderContentType, httpclient.MIMEApplicationJSON)
	req.Header.Add(httpclient.HeaderAuthorization, httpclient.AuthSchemeBearer+token)
	opts := httpclient.ReqOptBuidler().
		Log().
		LogReqBodyOnlyError().
		LoggedReqBody([]string{}).
		LogResBody().
		LoggedResBody([]string{}).
		Build()

	res, err := i.Client().Do(req, opts)
	if err != nil {
		return err
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			i.logger.Error(ctx, utils.ErrorCloseResponseBody, "error", err.Error())
		}
	}()

	if res.StatusCode == http.StatusOK || res.StatusCode == http.StatusNoContent {
		return nil
	}

	return fmt.Errorf("call api update identity provider status: %s", res.Status)
}

// DeleteIdP implements KcIdentityProviderService.
func (i *identityProvider) DeleteIdP(
	ctx context.Context,
	realm string,
	alias string,
	token string,
) error {
	url := i.AdminRealmUrl(realm) + identityProviderPath + "/" + alias
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		i.logger.Error(ctx, utils.ErrorCreateReq, "error", err.Error())
		return err
	}

	req.Header.Add(httpclient.HeaderAuthorization, httpclient.AuthSchemeBearer+token)
	opts := httpclient.ReqOptBuidler().
		Log().
		LogReqBodyOnlyError().
		LoggedReqBody([]string{}).
		LogResBody().
		LoggedResBody([]string{}).
		Build()

	res, err := i.Client().Do(req, opts)
	if err != nil {
		return err
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			i.logger.Error(ctx, utils.ErrorCloseResponseBody, "error", err.Error())
		}
	}()

	if res.StatusCode == http.StatusOK {
		return nil
	}

	return fmt.Errorf("call api delete identity provider status: %s", res.Status)
}
