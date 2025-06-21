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

	"github.com/dpe27/es-krake/internal/domain/auth/repository"
	"github.com/dpe27/es-krake/internal/infrastructure/httpclient"
	kcdto "github.com/dpe27/es-krake/internal/infrastructure/keycloak/dto"
	"github.com/dpe27/es-krake/pkg/log"
	"github.com/dpe27/es-krake/pkg/nethttp"
	"github.com/dpe27/es-krake/pkg/utils"
)

const (
	GoogleTokenPath = "/token"

	GoogleUserInfoPath = "/v1/userinfo"
)

type SnsProviderService interface {
	GetTokenByAuthCode(ctx context.Context, providerID int, clientID, clientSecret, code, redirectURI string) (kcdto.TokenEndpointResp, error)

	GetUserInfo(ctx context.Context, providerID int, token kcdto.TokenEndpointResp) (kcdto.SnsUserInfo, error)
}

type snsProviderService struct {
	logger *log.Logger
	cli    map[int]httpclient.HttpClient
}

func NewSnsProviderService(cli map[int]httpclient.HttpClient) SnsProviderService {
	return &snsProviderService{
		logger: log.With("service", "sns_provider_service"),
		cli:    cli,
	}
}

func (s *snsProviderService) getTokenEndpoint(providerID int) string {
	switch providerID {
	case int(repository.GoogleProvider):
		return "https://oauth2.googleapis.com" + GoogleTokenPath
	}
	return ""
}

func (s *snsProviderService) getUserInfoEndpoint(providerID int) string {
	switch providerID {
	case int(repository.GoogleProvider):
		return "https://openidconnect.googleapis.com" + GoogleUserInfoPath
	}
	return ""
}

func (s *snsProviderService) GetTokenByAuthCode(
	ctx context.Context,
	providerID int,
	clientID string,
	clientSecret string,
	code string,
	redirectURI string,
) (kcdto.TokenEndpointResp, error) {
	endpoint := s.getTokenEndpoint(providerID)
	u, err := url.ParseRequestURI(endpoint)
	if err != nil {
		return kcdto.TokenEndpointResp{}, err
	}

	data := url.Values{}
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", redirectURI)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		u.String(),
		strings.NewReader(data.Encode()),
	)
	if err != nil {
		s.logger.Error(ctx, utils.ErrorCreateReq, "error", err.Error())
		return kcdto.TokenEndpointResp{}, err
	}

	req.Header.Add(nethttp.HeaderContentType, nethttp.MIMEApplicationForm)
	req.Header.Add(nethttp.HeaderContentLength, strconv.Itoa(len(data.Encode())))

	client, ok := s.cli[providerID]
	if !ok {
		return kcdto.TokenEndpointResp{}, errors.New("requested provider not found")
	}

	opts := httpclient.ReqOptBuidler().
		Log().
		LogReqBodyOnlyError().
		LogResBody().
		LoggedReqBody([]string{}).
		LoggedResBody([]string{}).
		Build()

	resp, err := client.Do(req, opts)
	if err != nil {
		return kcdto.TokenEndpointResp{}, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			s.logger.Error(ctx, utils.ErrorCloseResponseBody, "error", err.Error())
		}
	}()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return kcdto.TokenEndpointResp{}, err
	}

	var respData map[string]interface{}
	err = json.Unmarshal(bodyBytes, &respData)
	if err != nil {
		return kcdto.TokenEndpointResp{}, err
	}

	if resp.StatusCode != http.StatusOK {
		errMsg, ok := respData["error"].(string)
		if !ok {
			return kcdto.TokenEndpointResp{}, fmt.Errorf(
				"cannot handle the error because the 'error' field is not a string. Got: %+v",
				respData["error"],
			)
		}
		errDesc, ok := respData["error_description"].(string)
		if ok {
			errMsg += ". " + errDesc
		}
		return kcdto.TokenEndpointResp{}, errors.New(errMsg)
	}

	token := kcdto.TokenEndpointResp{}
	err = utils.MapToStruct(respData, &token)
	return token, err
}

// GetUserInfo implements SnsProviderService.
func (s *snsProviderService) GetUserInfo(
	ctx context.Context,
	providerID int,
	token kcdto.TokenEndpointResp,
) (kcdto.SnsUserInfo, error) {
	endpoint := s.getUserInfoEndpoint(providerID)
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		endpoint,
		nil,
	)
	if err != nil {
		s.logger.Error(ctx, utils.ErrorCreateReq, "error", err.Error())
		return kcdto.SnsUserInfo{}, err
	}

	req.Header.Add(nethttp.HeaderAuthorization, nethttp.AuthSchemeBearer+token.AccessToken)

	client, ok := s.cli[providerID]
	if !ok {
		return kcdto.SnsUserInfo{}, errors.New("requested provider not found")
	}

	opts := httpclient.ReqOptBuidler().
		Log().
		LogReqBodyOnlyError().
		LogResBody().
		LoggedReqBody([]string{}).
		LoggedResBody([]string{}).
		Build()
	resp, err := client.Do(req, opts)
	if err != nil {
		return kcdto.SnsUserInfo{}, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			s.logger.Error(ctx, utils.ErrorCloseResponseBody, "error", err.Error())
		}
	}()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return kcdto.SnsUserInfo{}, err
	}

	var respData map[string]interface{}
	err = json.Unmarshal(bodyBytes, &respData)
	if err != nil {
		return kcdto.SnsUserInfo{}, err
	}

	if resp.StatusCode != http.StatusOK {
		errMsg, ok := respData["error"].(string)
		if !ok {
			return kcdto.SnsUserInfo{}, fmt.Errorf(
				"cannot handle the error because the 'error' field is not a string. Got: %+v",
				respData["error"],
			)
		}

		errDesc, ok := respData["error_description"].(string)
		if ok {
			errMsg += ". " + errDesc
		}
		return kcdto.SnsUserInfo{}, errors.New(errMsg)
	}

	userInfo := kcdto.SnsUserInfo{}
	err = json.Unmarshal(bodyBytes, &userInfo)
	return userInfo, err
}
