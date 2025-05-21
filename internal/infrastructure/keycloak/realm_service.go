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
	"github.com/dpe27/es-krake/pkg/nethttp"
	"github.com/dpe27/es-krake/pkg/utils"
)

type KcRealmService interface {
	GetRealm(ctx context.Context, realm, token string) (kcdto.KcRealm, error)
	PostRealm(ctx context.Context, body map[string]interface{}, token string) error
	PutRealm(ctx context.Context, body map[string]interface{}, realm, token string) error
}

type realmService struct {
	BaseKcService
	logger *log.Logger
}

func NewKcRealmService(base BaseKcService) KcRealmService {
	return &realmService{
		BaseKcService: base,
		logger:        log.With("service", "keycloak_realm_service"),
	}
}

// GetRealm implements RealmService.
func (r *realmService) GetRealm(
	ctx context.Context,
	realm string,
	token string,
) (kcdto.KcRealm, error) {
	url := r.AdminRealmUrl(realm)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		r.logger.Error(ctx, utils.ErrorCreateReq, "error", err.Error())
		return kcdto.KcRealm{}, err
	}

	req.Header.Add(nethttp.HeaderAuthorization, nethttp.AuthSchemeBearer+token)
	query := req.URL.Query()
	req.URL.RawQuery = query.Encode()

	opts := httpclient.ReqOptBuidler().
		Log().
		LogReqBodyOnlyError().
		LoggedReqBody([]string{}).
		LogResBody().
		LoggedResBody([]string{}).
		Build()

	res, err := r.Client().Do(req, opts)
	if err != nil {
		return kcdto.KcRealm{}, err
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			r.logger.Error(ctx, utils.ErrorCloseResponseBody, "error", err.Error())
		}
	}()

	if res.StatusCode == http.StatusNotFound {
		return kcdto.KcRealm{}, domainerr.ErrorNotFound
	}

	if res.StatusCode == http.StatusOK {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return kcdto.KcRealm{}, err
		}

		realm := kcdto.KcRealm{}
		err = json.Unmarshal(body, &realm)
		if err != nil {
			return kcdto.KcRealm{}, err
		}

		return realm, nil
	}

	return kcdto.KcRealm{}, fmt.Errorf("call api get realm status: %s", res.Status)
}

// PostRealm implements RealmService.
func (r *realmService) PostRealm(
	ctx context.Context,
	body map[string]interface{},
	token string,
) error {
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		r.logger.Error(ctx, utils.ErrorMarshalFailed, "error", err.Error())
		return err
	}

	reqBody := strings.NewReader(string(bodyBytes))
	url := r.AdminApiBaseUrl()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, reqBody)
	if err != nil {
		r.logger.Error(ctx, utils.ErrorCreateReq, "error", err.Error())
		return err
	}

	req.Header.Add(nethttp.HeaderContentType, nethttp.MIMEApplicationJSON)
	req.Header.Add(nethttp.HeaderAuthorization, nethttp.AuthSchemeBearer+token)
	opts := httpclient.ReqOptBuidler().
		Log().
		LogReqBodyOnlyError().
		LoggedReqBody([]string{}).
		LogResBody().
		LoggedResBody([]string{}).
		Build()

	res, err := r.Client().Do(req, opts)
	if err != nil {
		return err
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			r.logger.Error(ctx, utils.ErrorCloseResponseBody, "error", err.Error())
		}
	}()

	if res.StatusCode == http.StatusCreated {
		return nil
	}

	return fmt.Errorf("call api post realm status: %s", res.Status)
}

// PutRealm implements RealmService.
func (r *realmService) PutRealm(
	ctx context.Context,
	body map[string]interface{},
	realm string,
	token string,
) error {
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		r.logger.Error(ctx, utils.ErrorMarshalFailed, "error", err.Error())
		return err
	}

	reqBody := strings.NewReader(string(bodyBytes))
	url := r.AdminRealmUrl(realm)
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, reqBody)
	if err != nil {
		r.logger.Error(ctx, utils.ErrorCreateReq, "error", err.Error())
		return err
	}

	req.Header.Add(nethttp.HeaderContentType, nethttp.MIMEApplicationJSON)
	req.Header.Add(nethttp.HeaderAuthorization, nethttp.AuthSchemeBearer+token)
	opts := httpclient.ReqOptBuidler().
		Log().
		LogReqBodyOnlyError().
		LoggedReqBody([]string{}).
		LogResBody().
		LoggedResBody([]string{}).
		Build()

	res, err := r.Client().Do(req, opts)
	if err != nil {
		return err
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			r.logger.Error(ctx, utils.ErrorCloseResponseBody, "error", err.Error())
		}
	}()

	if res.StatusCode == http.StatusOK || res.StatusCode == http.StatusNoContent {
		return nil
	}

	return fmt.Errorf("call api put realm status: %s", res.Status)
}
