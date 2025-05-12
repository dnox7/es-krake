package keycloak

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/dpe27/es-krake/internal/infrastructure/httpclient"
	kcdto "github.com/dpe27/es-krake/internal/infrastructure/keycloak/dto"
	"github.com/dpe27/es-krake/pkg/log"
	"github.com/dpe27/es-krake/pkg/utils"
)

type KcTokenService interface {
	GetTokenWithPassword(ctx context.Context, realm, clientID, username, password string) (kcdto.TokenEndpointResp, error)
	GetTokenWithCode(ctx context.Context, realm, clientID, code, redirectURI string) (kcdto.TokenEndpointResp, error)
	RefreshToken(ctx context.Context, realm, clientID, refreshToken string) (kcdto.TokenEndpointResp, error)
}

type tokenService struct {
	BaseKcService
	logger *log.Logger
}

func NewKcTokenService(base BaseKcService) KcTokenService {
	return &tokenService{
		BaseKcService: base,
		logger:        log.With("service", "keycloak_token_service"),
	}
}

// GetTokenWithPassword implements KcTokenService.
func (t *tokenService) GetTokenWithPassword(
	ctx context.Context,
	realm string,
	clientID string,
	username string,
	password string,
) (kcdto.TokenEndpointResp, error) {
	params := map[string]string{
		"client_id":  clientID,
		"grant_type": "password",
		"username":   username,
		"password":   password,
	}
	return t.requestToken(ctx, realm, params)
}

// GetTokenWithCode implements KcTokenService.
func (t *tokenService) GetTokenWithCode(
	ctx context.Context,
	realm string,
	clientID string,
	code string,
	redirectURI string,
) (kcdto.TokenEndpointResp, error) {
	params := map[string]string{
		"client_id":    clientID,
		"grant_type":   "authorization_code",
		"code":         code,
		"redirect_uri": redirectURI,
	}
	return t.requestToken(ctx, realm, params)
}

// RefreshToken implements KcTokenService.
func (t *tokenService) RefreshToken(
	ctx context.Context,
	realm string,
	clientID string,
	refreshToken string,
) (kcdto.TokenEndpointResp, error) {
	params := map[string]string{
		"client_id":     clientID,
		"grant_type":    "refresh_token",
		"refresh_token": refreshToken,
	}
	return t.requestToken(ctx, realm, params)
}

func (t *tokenService) requestToken(
	ctx context.Context,
	realm string,
	params map[string]string,
) (kcdto.TokenEndpointResp, error) {
	endpoint := t.TokenUrl(realm)
	parsedUrl, err := url.ParseRequestURI(endpoint)
	if err != nil {
		return kcdto.TokenEndpointResp{}, err
	}

	data := url.Values{}
	for k, v := range params {
		data.Set(k, v)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		parsedUrl.String(),
		strings.NewReader(data.Encode()),
	)
	if err != nil {
		t.logger.Error(ctx, utils.ErrorCreateReq, "error", err.Error())
		return kcdto.TokenEndpointResp{}, err
	}

	req.Header.Add(httpclient.HeaderContentType, httpclient.MIMEApplicationForm)
	req.Header.Add(httpclient.HeaderContentLength, strconv.Itoa(len(data.Encode())))
	opts := httpclient.ReqOptBuidler().
		Log().
		LogReqBodyOnlyError().
		LogResBody().
		LoggedReqBody([]string{}).
		LoggedResBody([]string{}).
		Build()

	res, err := t.Client().Do(req, opts)
	if err != nil {
		return kcdto.TokenEndpointResp{}, err
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			t.logger.Error(ctx, utils.ErrorCloseResponseBody, "error", err.Error())
		}
	}()

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return kcdto.TokenEndpointResp{}, err
	}

	if res.StatusCode == http.StatusOK {
		token := kcdto.TokenEndpointResp{}
		err = json.Unmarshal(bodyBytes, &token)
		if err != nil {
			return kcdto.TokenEndpointResp{}, err
		}
		return token, nil
	}

	var body map[string]interface{}
	err = json.Unmarshal(bodyBytes, &body)
	if err != nil {
		return kcdto.TokenEndpointResp{}, err
	}

	errMsg, ok := body["error"].(string)
	if !ok {
		return kcdto.TokenEndpointResp{}, fmt.Errorf(
			"call apt get token (grant_type: %s) status: %s",
			params["grant_type"],
			res.Status,
		)
	}

	errDesc, ok := body["error_description"].(string)
	if ok {
		errMsg += ". " + errDesc
	}

	return kcdto.TokenEndpointResp{}, errors.New(errMsg)
}
