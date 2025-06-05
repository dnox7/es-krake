package usecase

import (
	"github.com/dpe27/es-krake/internal/domain/auth/repository"
	"github.com/dpe27/es-krake/pkg/log"
)

type AuthUsecaseDeps struct {
	RoleTypeRepo repository.RoleTypeRepository
}

type AuthUsecase struct {
	logger *log.Logger
	deps   *AuthUsecaseDeps
}

func NewAuthUsecase(deps *AuthUsecaseDeps) AuthUsecase {

	return AuthUsecase{
		logger: log.With("object", "auth_usecase"),
		deps:   deps,
	}
}
