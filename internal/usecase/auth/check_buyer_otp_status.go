package usecase

import (
	"context"
	"time"

	"github.com/dpe27/es-krake/internal/domain/auth/repository"
	"github.com/dpe27/es-krake/pkg/utils"
	"gorm.io/gorm"
)

type OtpStatus struct {
	Status               string
	EmailVerified        bool
	IsRegisteredPassword bool
	IsRegisteredSns      bool
}

func (u *AuthUsecase) CheckBuyerOtpStatus(ctx context.Context, otp string) (OtpStatus, error) {
	otpDetail, err := u.deps.OtpRepo.TakeByConditions(ctx, map[string]interface{}{
		"token": otp,
	}, nil)
	if err != nil && err != gorm.ErrRecordNotFound {
		return OtpStatus{}, err
	}

	if err == gorm.ErrRecordNotFound {
		return OtpStatus{
			Status: string(repository.OtpStatusInvalid),
		}, nil
	}

	now := time.Now()
	expiredAt, err := time.Parse(utils.FormatDateTimeISO, otpDetail.ExpiredAt)
	if err != nil {
		return OtpStatus{}, err
	}

	if now.After(expiredAt) {
		return OtpStatus{
			Status: string(repository.OtpStatusExpired),
		}, nil
	}

	buyerAccount, err := u.deps.BuyerAccountRepo.TakeByConditions(ctx, map[string]interface{}{
		"kc_user_id": otpDetail.KcUserID,
	}, nil)
	if err != nil {
		return OtpStatus{}, err
	}

	masterRealmToken, err := u.deps.KcTokenService.GetMasterRealmToken(ctx)
	if err != nil {
		return OtpStatus{}, err
	}

	buyerClient := u.deps.KcClientService.GetPlatformClient()
	realmName := buyerClient["realm_name"]

	existedPassword, err := u.deps.KcUserService.CheckPasswordExist(ctx, realmName, buyerAccount.KcUserID, masterRealmToken.AccessToken)
	if err != nil {
		return OtpStatus{}, err
	}

	if !existedPassword {
		err = u.deps.KcUserService.UpdateUser(ctx, map[string]interface{}{
			"emailVerified": true,
		}, realmName, buyerAccount.KcUserID, masterRealmToken.AccessToken)
		if err != nil {
			return OtpStatus{}, err
		}

		_, err := u.deps.BuyerAccountRepo.Update(ctx, buyerAccount, map[string]interface{}{
			"mail_verified": true,
		})
		if err != nil {
			return OtpStatus{}, err
		}
	}

	kcUser, err := u.deps.KcUserService.GetUserByID(ctx, realmName, buyerAccount.KcUserID, masterRealmToken.AccessToken)
	if err != nil {
		return OtpStatus{}, err
	}

	return OtpStatus{
		Status:               string(repository.OtpStatusValid),
		EmailVerified:        true,
		IsRegisteredPassword: existedPassword,
		IsRegisteredSns:      len(kcUser.FederatedIdentities) > 0,
	}, nil
}
