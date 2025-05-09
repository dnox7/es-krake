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

	domainerr "github.com/dpe27/es-krake/internal/domain/shared/errors"
	"github.com/dpe27/es-krake/internal/infrastructure/httpclient"
	kcdto "github.com/dpe27/es-krake/internal/infrastructure/keycloak/dto"
	"github.com/dpe27/es-krake/pkg/log"
	"github.com/dpe27/es-krake/pkg/utils"
)

const (
	userPath         = "/users"
	countPath        = "/count"
	resetPaswordPath = "/reset-password"
	credentialsPath  = "/credentials"
	logoutAllPath    = "/logout-all"
)

type KcUserService interface {
	GetUserByID(ctx context.Context, realm, userID, token string) (kcdto.KcUser, error)
	GetUserList(ctx context.Context, conditions map[string]interface{}, realm, token string) ([]kcdto.KcUser, error)
	PostUser(ctx context.Context, body map[string]interface{}, realm, token string) (kcdto.KcUser, error)
	UpdateUser(ctx context.Context, body map[string]interface{}, realm, userID, token string) error
	DeleteUser(ctx context.Context, realm, userID, token string) error
	CountUsers(ctx context.Context, conditions map[string]interface{}, realm, token string) (int, error)

	LogoutOIDC(ctx context.Context, realm, clientID, refreshToken string) error
	LogoutAll(ctx context.Context, realm, userID, token string) error
	ResetPassword(ctx context.Context, body map[string]interface{}, realm, userID, token string) error
	CheckPasswordExist(ctx context.Context, realm, userID, token string) (bool, error)
}

type userService struct {
	BaseKcService
	logger *log.Logger
}

func NewKcUserService(base BaseKcService) KcUserService {
	return &userService{
		BaseKcService: base,
		logger:        log.With("service", "keycloak_user_service"),
	}
}

// GetUserByID implements KcUserService.
func (u *userService) GetUserByID(
	ctx context.Context,
	realm string,
	userID string,
	token string,
) (kcdto.KcUser, error) {
	url := u.AdminRealmUrl(realm) + userPath + "/" + userID
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		u.logger.Error(ctx, utils.ErrorCreateReq, "error", err.Error())
		return kcdto.KcUser{}, err
	}

	q := req.URL.Query()
	req.URL.RawQuery = q.Encode()
	req.Header.Add(httpclient.HeaderAuthorization, httpclient.AuthSchemeBearer+token)

	opts := httpclient.ReqOptBuidler().
		Log().
		LogReqBodyOnlyError().
		LoggedReqBody([]string{}).
		LogResBody().
		LoggedResBody([]string{}).
		Build()

	res, err := u.Client().Do(req, opts)
	if err != nil {
		return kcdto.KcUser{}, err
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			u.logger.Error(ctx, utils.ErrorCloseResponseBody, "error", err.Error())
		}
	}()

	if res.StatusCode == http.StatusNotFound || res.StatusCode == http.StatusUnauthorized {
		return kcdto.KcUser{}, domainerr.ErrorNotFound
	}

	if res.StatusCode == http.StatusOK {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return kcdto.KcUser{}, err
		}

		user := kcdto.KcUser{}
		err = json.Unmarshal(body, &user)
		if err != nil {
			return kcdto.KcUser{}, err
		}

		return user, nil
	}

	return kcdto.KcUser{}, fmt.Errorf("call api get user by id status: %s", res.Status)
}

// GetUserList implements KcUserService.
func (u *userService) GetUserList(
	ctx context.Context,
	conditions map[string]interface{},
	realm string,
	token string,
) ([]kcdto.KcUser, error) {
	url := u.AdminRealmUrl(realm) + userPath
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		u.logger.Error(ctx, utils.ErrorCreateReq, "error", err.Error())
		return nil, err
	}

	q := req.URL.Query()
	var input string
	for k, v := range conditions {
		input, err = utils.ToString(v)
		if err != nil {
			u.logger.Error(ctx, err.Error())
			return nil, err
		}
		q.Add(k, input)
	}
	req.URL.RawQuery = q.Encode()
	req.Header.Add(httpclient.HeaderAuthorization, httpclient.AuthSchemeBearer+token)

	opts := httpclient.ReqOptBuidler().
		Log().
		LogReqBodyOnlyError().
		LoggedReqBody([]string{}).
		LogResBody().
		LoggedResBody([]string{}).
		Build()

	res, err := u.Client().Do(req, opts)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			u.logger.Error(ctx, utils.ErrorCloseResponseBody, "error", err.Error())
		}
	}()

	if res.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		var users []kcdto.KcUser
		err = json.Unmarshal(bodyBytes, &users)
		if err != nil {
			return nil, err
		}
		return users, nil
	}

	if res.StatusCode == http.StatusNotFound {
		return nil, domainerr.ErrorNotFound
	}

	return nil, fmt.Errorf("call api get users list status: %s", res.Status)
}

// PostUser implements KcUserService.
func (u *userService) PostUser(
	ctx context.Context,
	body map[string]interface{},
	realm string,
	token string,
) (kcdto.KcUser, error) {
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		u.logger.Error(ctx, utils.ErrorMarshalFailed, "error", err.Error())
		return kcdto.KcUser{}, err
	}

	url := u.AdminRealmUrl(realm) + userPath
	reqBody := strings.NewReader(string(bodyBytes))
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, reqBody)
	if err != nil {
		u.logger.Error(ctx, utils.ErrorCreateReq, "error", err.Error())
		return kcdto.KcUser{}, err
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

	res, err := u.Client().Do(req, opts)
	if err != nil {
		return kcdto.KcUser{}, err
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			u.logger.Error(ctx, utils.ErrorCloseResponseBody, "error", err.Error())
		}
	}()

	if res.StatusCode == http.StatusCreated {
		locations := strings.Split(res.Header.Get(httpclient.HeaderLocation), "/")
		return kcdto.KcUser{
			ID: locations[len(locations)-1],
		}, nil
	}
	return kcdto.KcUser{}, fmt.Errorf("call api create new user status: %s", res.Status)
}

// UpdateUser implements KcUserService.
func (u *userService) UpdateUser(
	ctx context.Context,
	body map[string]interface{},
	realm string,
	userID string,
	token string,
) error {
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		u.logger.Error(ctx, utils.ErrorMarshalFailed, "error", err.Error())
		return err
	}

	url := u.AdminRealmUrl(realm) + userPath + "/" + userID
	bodyReq := strings.NewReader(string(bodyBytes))
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, bodyReq)
	if err != nil {
		u.logger.Error(ctx, utils.ErrorCreateReq, "error", err.Error())
		return err
	}

	req.Header.Add(httpclient.HeaderAuthorization, httpclient.AuthSchemeBearer+token)
	req.Header.Add(httpclient.HeaderContentType, httpclient.MIMEApplicationJSON)
	opts := httpclient.ReqOptBuidler().
		Log().
		LogReqBodyOnlyError().
		LoggedReqBody([]string{}).
		LogResBody().
		LoggedResBody([]string{}).
		Build()

	res, err := u.Client().Do(req, opts)
	if err != nil {
		return err
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			u.logger.Error(ctx, utils.ErrorCloseResponseBody, "error", err.Error())
		}
	}()

	if res.StatusCode == http.StatusNoContent {
		return nil
	}

	return fmt.Errorf("call api update user satus: %s", res.Status)
}

// DeleteUser implements KcUserService.
func (u *userService) DeleteUser(
	ctx context.Context,
	realm string,
	userID string,
	token string,
) error {
	url := u.AdminRealmUrl(realm) + userPath + "/" + userID
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		u.logger.Error(ctx, utils.ErrorCreateReq, "error", err.Error())
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

	res, err := u.Client().Do(req, opts)
	if err != nil {
		return err
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			u.logger.Error(ctx, utils.ErrorCloseResponseBody, "error", err.Error())
		}
	}()

	if res.StatusCode == http.StatusNoContent {
		return nil
	}

	return fmt.Errorf("call api delete user status: %s", res.Status)
}

// CountUsers implements KcUserService.
func (u *userService) CountUsers(
	ctx context.Context,
	conditions map[string]interface{},
	realm string,
	token string,
) (int, error) {
	url := u.AdminRealmUrl(realm) + userPath + countPath
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		u.logger.Error(ctx, utils.ErrorCreateReq, "error", err.Error())
		return 0, err
	}

	query := req.URL.Query()
	var valStr string
	for k, v := range conditions {
		valStr, err = utils.ToString(v)
		if err != nil {
			return 0, err
		}
		query.Add(k, valStr)
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

	res, err := u.Client().Do(req, opts)
	if err != nil {
		return 0, err
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			u.logger.Error(ctx, utils.ErrorCloseResponseBody, "error", err.Error())
		}
	}()

	if res.StatusCode == http.StatusOK {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return 0, err
		}
		var count int
		err = json.Unmarshal(body, &count)
		if err != nil {
			return 0, err
		}
		return count, nil
	}

	return 0, fmt.Errorf("call api users count status: %s", res.Status)
}

// Logout implements KcUserService.
func (u *userService) LogoutOIDC(
	ctx context.Context,
	realm string,
	clientID string,
	refreshToken string,
) error {
	endpoint := u.LogoutUrl(realm)
	data := url.Values{}
	data.Set("client_id", clientID)
	data.Set("refresh_token", refreshToken)

	parsedUrl, err := url.ParseRequestURI(endpoint)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, parsedUrl.String(), strings.NewReader(data.Encode()))
	if err != nil {
		u.logger.Error(ctx, utils.ErrorCreateReq, "error", err.Error())
		return err
	}

	req.Header.Add(httpclient.HeaderContentType, httpclient.MIMEApplicationForm)
	req.Header.Add(httpclient.HeaderContentLength, strconv.Itoa(len(data.Encode())))
	opts := httpclient.ReqOptBuidler().
		Log().
		LogReqBodyOnlyError().
		LoggedReqBody([]string{}).
		LogResBody().
		LoggedResBody([]string{}).
		Build()

	res, err := u.Client().Do(req, opts)
	if err != nil {
		return err
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			u.logger.Error(ctx, utils.ErrorCloseResponseBody, "error", err.Error())
		}
	}()

	if res.StatusCode == http.StatusNoContent {
		return nil
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var respData map[string]interface{}
	if err := json.Unmarshal(body, &respData); err != nil {
		return err
	}

	errRes, errResIsStr := respData["error"].(string)
	errDesc, errDescIsStr := respData["error_description"].(string)

	if errResIsStr && errDescIsStr {
		return errors.New(errRes + ". " + errDesc)
	}

	return fmt.Errorf("call api logout OIDC status: %s", res.Status)
}

// LogoutAll implements KcUserService.
func (u *userService) LogoutAll(
	ctx context.Context,
	realm string,
	userID string,
	token string,
) error {
	url := u.AdminRealmUrl(realm) + userPath + "/" + userID + logoutPath
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
	if err != nil {
		u.logger.Error(ctx, utils.ErrorCreateReq, "error", err.Error())
		return err
	}

	req.Header.Add(httpclient.HeaderAuthorization, httpclient.AuthSchemeBearer+token)
	opts := httpclient.ReqOptBuidler().
		Log().
		LoggedReqBody([]string{}).
		LoggedResBody([]string{}).
		LogReqBodyOnlyError().
		LogResBody().
		Build()

	res, err := u.Client().Do(req, opts)
	if err != nil {
		return err
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			u.logger.Error(ctx, utils.ErrorCloseResponseBody, "error", err.Error())
		}
	}()

	if res.StatusCode == http.StatusNoContent {
		return nil
	}

	return fmt.Errorf("call api logout all status: %s", res.Status)
}

// ResetPassword implements KcUserService.
func (u *userService) ResetPassword(
	ctx context.Context,
	body map[string]interface{},
	realm string,
	userID string,
	token string,
) error {
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		u.logger.Error(ctx, utils.ErrorMarshalFailed, "error", err.Error())
		return err
	}

	reqBody := strings.NewReader(string(bodyBytes))
	url := u.AdminRealmUrl(realm) + userPath + "/" + userID + resetPaswordPath
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, reqBody)
	if err != nil {
		u.logger.Error(ctx, utils.ErrorCreateReq, "error", err.Error())
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

	res, err := u.Client().Do(req, opts)
	if err != nil {
		return err
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			u.logger.Error(ctx, utils.ErrorCloseResponseBody, "error", err.Error())
		}
	}()

	if res.StatusCode == http.StatusNoContent {
		return nil
	}

	return fmt.Errorf("call api reset password status: %s", res.Status)
}

// CheckPasswordExist implements KcUserService.
func (u *userService) CheckPasswordExist(
	ctx context.Context,
	realm string,
	userID string,
	token string,
) (bool, error) {
	url := u.AdminRealmUrl(realm) + userPath + "/" + userID + credentialsPath
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		u.logger.Error(ctx, utils.ErrorCreateReq, "error", err.Error())
		return false, err
	}

	req.Header.Add(httpclient.HeaderAuthorization, httpclient.AuthSchemeBearer+token)
	opts := httpclient.ReqOptBuidler().
		Log().
		LogReqBodyOnlyError().
		LoggedReqBody([]string{}).
		LogResBody().
		LoggedResBody([]string{}).
		Build()

	res, err := u.Client().Do(req, opts)
	if err != nil {
		return false, err
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			u.logger.Error(ctx, utils.ErrorCloseResponseBody, "error", err.Error())
		}
	}()

	if res.StatusCode == http.StatusOK {
		body, err := io.ReadAll(req.Body)
		if err != nil {
			return false, err
		}

		data := []map[string]interface{}{}
		err = json.Unmarshal(body, &data)
		if err != nil {
			return false, err
		}

		for _, m := range data {
			if m["type"] == "password" {
				return true, nil
			}
		}
		return false, nil
	}

	return false, fmt.Errorf("call api check exist password: %s", res.Status)
}
