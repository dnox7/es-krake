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

const (
	authPath          = "/authentication"
	flowPath          = "/flows"
	execPath          = "/executions"
	copyPath          = "/copy"
	lowerPriorityPath = "/lower-priority"
	raisePriorityPath = "/raise-priority"
	authFlowPath      = authPath + flowPath
	authExecPath      = authPath + execPath
)

type KcAuthenticationService interface {
	// Get authentication executions for a flow
	GetExecs(ctx context.Context, realm, flowAlias, token string) ([]kcdto.KcAuthExecInfo, error)
	// Update authentication executions of a flow
	PutExec(ctx context.Context, body map[string]interface{}, realm, flowAlias, token string) error
	// Delete execution by ID
	DeleteExec(ctx context.Context, realm, execID, token string) error
	// Raise execution’s priority
	PromoteExec(ctx context.Context, realm, execID, token string) error
	// Lower execution’s priority
	DemoteExec(ctx context.Context, realm, execID, token string) error
	// Get authentication flows. Returns a stream of authentication flows.
	GetFlows(ctx context.Context, realm, token string) ([]kcdto.KcAuthFlow, error)
	// Copy existing authentication flow under a new name.
	// The new name is given as 'newName' attribute of the passed JSON object
	PostFlowCopy(ctx context.Context, body map[string]interface{}, realm, flowAlias, token string) error
	// Update an authentication flow
	PutFlow(ctx context.Context, body map[string]interface{}, realm, flowID, token string) error
	// Delete an authentication flow
	DeleteFlow(ctx context.Context, realm, flowID, token string) error
}

type authenService struct {
	BaseKcService
	logger *log.Logger
}

func NewKcAuthenticationService(base BaseKcService) KcAuthenticationService {
	return &authenService{
		BaseKcService: base,
		logger:        log.With("service", "keycloak_authentication_service"),
	}
}

// GetExecs implements KcAuthenticationService.
func (a *authenService) GetExecs(
	ctx context.Context,
	realm string,
	flowAlias string,
	token string,
) ([]kcdto.KcAuthExecInfo, error) {
	url := a.AdminRealmUrl(realm) + authExecPath + "/" + flowAlias + execPath
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		a.logger.Error(ctx, utils.ErrorCreateReq, "error", err.Error())
		return nil, err
	}

	req.Header.Add(nethttp.HeaderAuthorization, nethttp.AuthSchemeBearer+token)
	query := req.URL.Query()
	req.URL.RawQuery = query.Encode()
	opts := httpclient.ReqOptBuidler().
		Log().LogReqBodyOnlyError().
		LogResBody().
		LoggedResBody([]string{}).
		LoggedReqBody([]string{}).
		Build()

	res, err := a.Client().Do(req, opts)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			a.logger.Error(ctx, utils.ErrorCloseResponseBody, "error", err.Error())
		}
	}()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("call api get AuthenticationExecutions with flowAlias: %s, status: %s", flowAlias, res.Status)
	}

	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	execs := []kcdto.KcAuthExecInfo{}
	err = json.Unmarshal(bytes, &execs)
	return execs, err
}

// PutExec implements KcAuthenticationService.
func (a *authenService) PutExec(
	ctx context.Context,
	body map[string]interface{},
	realm string,
	flowAlias string,
	token string,
) error {
	url := a.AdminRealmUrl(realm) + authFlowPath + "/" + flowAlias + execPath
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		a.logger.Error(ctx, utils.ErrorMarshalFailed, "error", err.Error())
		return err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPut,
		url,
		strings.NewReader(string(bodyBytes)),
	)
	if err != nil {
		a.logger.Error(ctx, utils.ErrorCreateReq, "error", err.Error())
		return err
	}

	req.Header.Add(nethttp.HeaderContentType, nethttp.MIMEApplicationJSON)
	req.Header.Add(nethttp.HeaderAuthorization, nethttp.AuthSchemeBearer+token)
	opts := httpclient.ReqOptBuidler().
		LogReqBodyOnlyError().
		LogResBody().
		LoggedReqBody([]string{}).
		LoggedResBody([]string{}).
		Log().Build()

	res, err := a.Client().Do(req, opts)
	if err != nil {
		return err
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			a.logger.Error(ctx, utils.ErrorCloseResponseBody, "error", err.Error())
		}
	}()

	if res.StatusCode == http.StatusAccepted {
		return nil
	}
	return fmt.Errorf("call api update authentication executions with flowAlias: %s, status: %s", flowAlias, res.Status)
}

// DeleteExec implements KcAuthenticationService.
func (a *authenService) DeleteExec(
	ctx context.Context,
	realm string,
	execID string,
	token string,
) error {
	url := a.AdminRealmUrl(realm) + authExecPath + "/" + execID
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		a.logger.Error(ctx, utils.ErrorCreateReq, "error", err.Error())
		return err
	}

	req.Header.Add(nethttp.HeaderAuthorization, nethttp.AuthSchemeBearer+token)
	opts := httpclient.ReqOptBuidler().
		LogReqBodyOnlyError().
		LogResBody().
		LoggedReqBody([]string{}).
		LoggedResBody([]string{}).
		Log().Build()

	res, err := a.Client().Do(req, opts)
	if err != nil {
		return nil
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			a.logger.Error(ctx, utils.ErrorCloseResponseBody, "error", err.Error())
		}
	}()

	if res.StatusCode == http.StatusNoContent {
		return nil
	}

	return fmt.Errorf("call api Delete execution with ID: %s, status: %s", execID, res.Status)
}

// PromoteExec implements KcAuthenticationService.
func (a *authenService) PromoteExec(ctx context.Context, realm string, execID string, token string) error {
	url := a.AdminRealmUrl(realm) + authExecPath + "/" + execID + raisePriorityPath
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
	if err != nil {
		a.logger.Error(ctx, utils.ErrorCreateReq, "error", err.Error())
		return err
	}

	req.Header.Add(nethttp.HeaderAuthorization, nethttp.AuthSchemeBearer+token)
	opts := httpclient.ReqOptBuidler().
		Log().
		LogReqBodyOnlyError().
		LoggedReqBody([]string{}).
		LogResBody().
		LoggedResBody([]string{}).
		Build()

	res, err := a.Client().Do(req, opts)
	if err != nil {
		return err
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			a.logger.Error(ctx, utils.ErrorCloseResponseBody, "error", err.Error())
		}
	}()

	if res.StatusCode == http.StatusOK {
		return nil
	}

	return fmt.Errorf("call api Raise execution’s priority witd executionID: %s, status: %s", execID, res.Status)
}

// DemoteExec implements KcAuthenticationService.
func (a *authenService) DemoteExec(ctx context.Context, realm string, execID string, token string) error {
	url := a.AdminRealmUrl(realm) + authExecPath + "/" + execID + lowerPriorityPath
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
	if err != nil {
		a.logger.Error(ctx, utils.ErrorCreateReq, "error", err.Error())
		return err
	}

	req.Header.Add(nethttp.HeaderAuthorization, nethttp.AuthSchemeBearer+token)
	opts := httpclient.ReqOptBuidler().
		Log().
		LogReqBodyOnlyError().
		LoggedReqBody([]string{}).
		LogResBody().
		LoggedResBody([]string{}).
		Build()

	res, err := a.Client().Do(req, opts)
	if err != nil {
		return err
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			a.logger.Error(ctx, utils.ErrorCloseResponseBody, "error", err.Error())
		}
	}()

	if res.StatusCode == http.StatusOK {
		return nil
	}

	return fmt.Errorf("call api Lower execution’s priority witd executionID: %s, status: %s", execID, res.Status)
}

// GetFlows implements KcAuthenticationService.
func (a *authenService) GetFlows(ctx context.Context, realm string, token string) ([]kcdto.KcAuthFlow, error) {
	url := a.AdminRealmUrl(realm) + authFlowPath
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		a.logger.Error(ctx, utils.ErrorCreateReq, "error", err.Error())
		return nil, err
	}

	req.Header.Add(nethttp.HeaderAuthorization, nethttp.AuthSchemeBearer+token)
	query := req.URL.Query()
	req.URL.RawQuery = query.Encode()
	opts := httpclient.ReqOptBuidler().
		Log().LogReqBodyOnlyError().
		LogResBody().
		LoggedResBody([]string{}).
		LoggedReqBody([]string{}).
		Build()

	res, err := a.Client().Do(req, opts)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			a.logger.Error(ctx, utils.ErrorCloseResponseBody, "error", err.Error())
		}
	}()

	if res.StatusCode == http.StatusNotFound {
		return nil, domainerr.ErrorNotFound
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("call api Get authentication flows with status: %s", res.Status)
	}

	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	flows := []kcdto.KcAuthFlow{}
	err = json.Unmarshal(bytes, &flows)
	return flows, err
}

// PostFlowCopy implements KcAuthenticationService.
func (a *authenService) PostFlowCopy(
	ctx context.Context,
	body map[string]interface{},
	realm string,
	flowAlias string,
	token string,
) error {
	url := a.AdminRealmUrl(realm) + authFlowPath + "/" + flowAlias + copyPath
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		a.logger.Error(ctx, utils.ErrorMarshalFailed, "error", err.Error())
		return err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		url,
		strings.NewReader(string(bodyBytes)),
	)
	if err != nil {
		a.logger.Error(ctx, utils.ErrorCreateReq, "error", err.Error())
		return err
	}

	req.Header.Add(nethttp.HeaderContentType, nethttp.MIMEApplicationJSON)
	req.Header.Add(nethttp.HeaderAuthorization, nethttp.AuthSchemeBearer+token)
	opts := httpclient.ReqOptBuidler().
		Log().LogReqBodyOnlyError().
		LogResBody().
		LoggedResBody([]string{}).
		LoggedReqBody([]string{}).
		Build()

	res, err := a.Client().Do(req, opts)
	if err != nil {
		return nil
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			a.logger.Error(ctx, utils.ErrorCloseResponseBody, "error", err.Error())
		}
	}()

	if res.StatusCode == http.StatusCreated {
		return nil
	}

	return fmt.Errorf("call api Copy existing authentication flow under a new name with status: %s", res.Status)
}

// PutFlow implements KcAuthenticationService.
func (a *authenService) PutFlow(
	ctx context.Context,
	body map[string]interface{},
	realm string,
	flowID string,
	token string,
) error {
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		a.logger.Error(ctx, utils.ErrorMarshalFailed, "error", err.Error())
		return err
	}

	url := a.AdminRealmUrl(realm) + authFlowPath + "/" + flowID
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPut,
		url,
		strings.NewReader(string(bodyBytes)),
	)
	if err != nil {
		a.logger.Error(ctx, utils.ErrorCreateReq, "error", err.Error())
		return err
	}

	req.Header.Add(nethttp.HeaderContentType, nethttp.MIMEApplicationJSON)
	req.Header.Add(nethttp.HeaderAuthorization, nethttp.AuthSchemeBearer+token)
	opts := httpclient.ReqOptBuidler().
		LogReqBodyOnlyError().
		LogResBody().
		LoggedReqBody([]string{}).
		LoggedResBody([]string{}).
		Log().Build()

	res, err := a.Client().Do(req, opts)
	if err != nil {
		return err
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			a.logger.Error(ctx, utils.ErrorCloseResponseBody, "error", err.Error())
		}
	}()

	if res.StatusCode != http.StatusNoContent {
		return nil
	}
	return fmt.Errorf("call api Update authentication flow with id: %s, status: %s", flowID, res.Status)
}

// DeleteFlow implements KcAuthenticationService.
func (a *authenService) DeleteFlow(ctx context.Context, realm string, flowID string, token string) error {
	url := a.AdminRealmUrl(realm) + authFlowPath + "/" + flowID
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		a.logger.Error(ctx, utils.ErrorCreateReq, "error", err.Error())
		return err
	}

	req.Header.Add(nethttp.HeaderAuthorization, nethttp.AuthSchemeBearer+token)
	opts := httpclient.ReqOptBuidler().
		LogReqBodyOnlyError().
		LogResBody().
		LoggedReqBody([]string{}).
		LoggedResBody([]string{}).
		Log().Build()

	res, err := a.Client().Do(req, opts)
	if err != nil {
		return nil
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			a.logger.Error(ctx, utils.ErrorCloseResponseBody, "error", err.Error())
		}
	}()

	if res.StatusCode == http.StatusNoContent {
		return nil
	}

	return fmt.Errorf("call api delete authentication flow with status: %s", res.Status)
}
