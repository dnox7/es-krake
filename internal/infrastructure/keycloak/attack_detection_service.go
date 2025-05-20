package keycloak

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	domainerr "github.com/dpe27/es-krake/internal/domain/shared/errors"
	"github.com/dpe27/es-krake/internal/infrastructure/httpclient"
	kcdto "github.com/dpe27/es-krake/internal/infrastructure/keycloak/dto"
	"github.com/dpe27/es-krake/pkg/log"
	"github.com/dpe27/es-krake/pkg/utils"
)

const bruteForceDetectionPath = "/attack-detection/brute-force/users"

type KcAttackDetectionService interface {
	GetUserBruteForceStatus(ctx context.Context, realm, userID, token string) (kcdto.BruteForceStatus, error)
	ClearAllBruteForceStates(ctx context.Context, realm, token string) error
	ResetBruteForceForUser(ctx context.Context, realm, userID, token string) error
}

type attackDectectionService struct {
	BaseKcService
	logger *log.Logger
}

func NewKcAttackDetectionService(base BaseKcService) KcAttackDetectionService {
	return &attackDectectionService{
		BaseKcService: base,
		logger:        log.With("service", "keycloak_attack_detection_service"),
	}
}

// GetUserBruteForceStatus implements KcAttackDetectionService.
func (a *attackDectectionService) GetUserBruteForceStatus(
	ctx context.Context,
	realm string,
	userID string,
	token string,
) (kcdto.BruteForceStatus, error) {
	url := a.AdminRealmUrl(realm) + bruteForceDetectionPath + "/" + userID
	req, err := http.NewRequestWithContext(ctx, url, http.MethodGet, nil)
	if err != nil {
		a.logger.Error(ctx, utils.ErrorCreateReq, "error", err.Error())
		return kcdto.BruteForceStatus{}, err
	}

	req.Header.Add(httpclient.HeaderAuthorization, httpclient.AuthSchemeBearer+token)
	req.Header.Add(httpclient.HeaderContentType, httpclient.MIMEApplicationJSON)
	opts := httpclient.ReqOptBuidler().
		Log().
		LogReqBodyOnlyError().
		LogResBody().
		LoggedReqBody([]string{}).
		LoggedResBody([]string{}).
		Build()

	res, err := a.Client().Do(req, opts)
	if err != nil {
		return kcdto.BruteForceStatus{}, err
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			a.logger.Error(ctx, utils.ErrorCloseResponseBody, "error", err.Error())
		}
	}()

	if res.StatusCode == http.StatusNotFound {
		return kcdto.BruteForceStatus{}, domainerr.ErrorNotFound
	}

	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return kcdto.BruteForceStatus{}, err
	}

	if res.StatusCode == http.StatusOK {
		status := kcdto.BruteForceStatus{}
		err = json.Unmarshal(bytes, &status)
		if err != nil {
			return kcdto.BruteForceStatus{}, err
		}
		return status, nil
	}

	var errMsg interface{}
	err = json.Unmarshal(bytes, &errMsg)
	if err != nil {
		return kcdto.BruteForceStatus{}, err
	}

	err = fmt.Errorf(
		"call api get status of a username in brute force detection, status: %s, error message: %v",
		res.Status, errMsg,
	)
	return kcdto.BruteForceStatus{}, err
}

// ClearAllBruteForceStates implements KcAttackDetectionService.
func (a *attackDectectionService) ClearAllBruteForceStates(
	ctx context.Context,
	realm string,
	token string,
) error {
	url := a.AdminRealmUrl(realm) + bruteForceDetectionPath
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		a.logger.Error(ctx, utils.ErrorCreateReq, "error", err.Error())
		return err
	}

	req.Header.Add(httpclient.HeaderAuthorization, httpclient.AuthSchemeBearer+token)
	reqOpts := httpclient.ReqOptBuidler().
		LogReqBodyOnlyError().
		LogResBody().
		LoggedReqBody([]string{}).
		LoggedResBody([]string{}).
		Log().
		Build()

	res, err := a.Client().Do(req, reqOpts)
	if err != nil {
		return err
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			a.logger.Error(ctx, utils.ErrorCloseResponseBody, "error", err.Error())
		}
	}()

	if res.StatusCode == http.StatusNoContent {
		return nil
	}

	return fmt.Errorf("call api Clear any user login failures for all users, status: %s", res.Status)
}

// ResetBruteForceForUser implements KcAttackDetectionService.
func (a *attackDectectionService) ResetBruteForceForUser(
	ctx context.Context,
	realm string,
	userID string,
	token string,
) error {
	url := a.AdminRealmUrl(realm) + bruteForceDetectionPath + "/" + userID
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		a.logger.Error(ctx, utils.ErrorCreateReq, "error", err.Error())
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

	res, err := a.Client().Do(req, opts)
	if err != nil {
		return err
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			a.logger.Error(ctx, utils.ErrorCloseResponseBody, "error", err.Error())
		}
	}()

	if res.StatusCode == http.StatusNoContent {
		return nil
	}

	return fmt.Errorf("Clear any user login failures for the use, status: %s", res.Status)
}
