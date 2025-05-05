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

type RealmService interface {
	GetRealm(ctx context.Context, realm, token string) (kcdto.KcRealm, error)
}

type realmService struct {
	BaseKcService
	logger *log.Logger
}

func NewRealmService(base BaseKcService) RealmService {
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
	endpoint := r.BaseKcService.AdminRealmUrl(realm)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		r.logger.Error(ctx, utils.ErrorCreateReq, "error", err.Error())
		return kcdto.KcRealm{}, err
	}

	req.Header.Add(utils.HeaderAuthorization, "Bearer "+token)
	query := req.URL.Query()
	req.URL.RawQuery = query.Encode()

	opt := httpclient.ReqOpt{}
	res, err := r.Client().Do(
		req,
		opt.LoggedRequestBodyOnlyError([]string{}),
		opt.LoggedResponseBody([]string{}),
	)
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
