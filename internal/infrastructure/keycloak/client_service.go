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

const clientPath = "/clients"

type KcClientService interface {
	GetClientList(ctx context.Context, conditions map[string]string, realm, token string) ([]kcdto.KcClient, error)
	PostClient(ctx context.Context, body map[string]interface{}, realm, token string) error
	PutClient(ctx context.Context, body map[string]interface{}, realm, uuid, token string) error
}

type clientService struct {
	BaseKcService
	logger *log.Logger
}

func NewKcClientService(base BaseKcService) KcClientService {
	return &clientService{
		BaseKcService: base,
		logger:        log.With("service", "keycloak_client_service"),
	}
}

// GetClientList implements KcClientService.
func (c *clientService) GetClientList(
	ctx context.Context,
	conditions map[string]string,
	realm string,
	token string,
) ([]kcdto.KcClient, error) {
	url := c.AdminRealmUrl(realm) + clientPath
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		c.logger.Error(ctx, utils.ErrorCreateReq, "error", err.Error())
		return nil, err
	}

	query := req.URL.Query()
	for k, v := range conditions {
		query.Add(k, v)
	}
	req.URL.RawQuery = query.Encode()
	req.Header.Add(httpclient.HeaderAuthorization, httpclient.AuthSchemeBearer+token)
	opts := httpclient.ReqOptBuidler().
		Log().
		LogReqBodyOnlyError().
		LoggedReqBody([]string{}).
		LogResBody().
		LoggedResBody([]string{}).
		Build()

	res, err := c.Client().Do(req, opts)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			c.logger.Error(ctx, utils.ErrorCloseResponseBody, "error", err.Error())
		}
	}()

	if res.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		var clients []kcdto.KcClient
		err = json.Unmarshal(bodyBytes, &clients)
		if err != nil {
			return nil, err
		}
		return clients, nil
	}

	if res.StatusCode == http.StatusNotFound {
		return nil, domainerr.ErrorNotFound
	}

	return nil, fmt.Errorf("call api get client list status: %s", res.Status)
}

// PostClient implements KcClientService.
func (c *clientService) PostClient(
	ctx context.Context,
	body map[string]interface{},
	realm string,
	token string,
) error {
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		c.logger.Error(ctx, utils.ErrorMarshalFailed, "error", err.Error())
		return err
	}

	url := c.AdminRealmUrl(realm) + clientPath
	reqBody := strings.NewReader(string(bodyBytes))
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, reqBody)
	if err != nil {
		c.logger.Error(ctx, utils.ErrorCreateReq, "error", err.Error())
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

	res, err := c.Client().Do(req, opts)
	if err != nil {
		return err
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			c.logger.Error(ctx, utils.ErrorCloseResponseBody, "error", err.Error())
		}
	}()

	if res.StatusCode == http.StatusCreated {
		return nil
	}

	return fmt.Errorf("call api create client status: %s", res.Status)
}

// PutClient implements KcClientService.
func (c *clientService) PutClient(
	ctx context.Context,
	body map[string]interface{},
	realm string,
	uuid string,
	token string,
) error {
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		c.logger.Error(ctx, utils.ErrorMarshalFailed, "error", err.Error())
		return err
	}

	reqBody := strings.NewReader(string(bodyBytes))
	url := c.AdminRealmUrl(realm) + clientPath + "/" + uuid
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, reqBody)
	if err != nil {
		c.logger.Error(ctx, utils.ErrorCreateReq, "error", err.Error())
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

	res, err := c.Client().Do(req, opts)
	if err != nil {
		return err
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			c.logger.Error(ctx, utils.ErrorCloseResponseBody, "error", err.Error())
		}
	}()

	if res.StatusCode == http.StatusOK || res.StatusCode == http.StatusNoContent {
		return nil
	}

	return fmt.Errorf("call api update identity provider status: %s", res.Status)
}
