package usecase

import (
	"context"
	"errors"
	"time"

	netMail "net/mail"

	"github.com/dpe27/es-krake/internal/domain/auth/entity"
	kcdto "github.com/dpe27/es-krake/internal/infrastructure/keycloak/dto"
	gormtx "github.com/dpe27/es-krake/internal/infrastructure/rdb/gorm/transaction"
	mailService "github.com/dpe27/es-krake/internal/infrastructure/service/mail"
	baseUtils "github.com/dpe27/es-krake/pkg/utils"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

func (u *AuthUsecase) SignupBuyer(ctx context.Context, kcUserInput map[string]interface{}) error {
	masterRealmToken, err := u.deps.KcTokenService.GetMasterRealmToken(ctx)
	if err != nil {
		return err
	}

	buyerClient := u.deps.KcClientService.GetPlatformClient()
	realmName := buyerClient["realm_name"]

	email := kcUserInput["email"].(string)
	users, err := u.deps.KcUserService.GetUserList(ctx, map[string]interface{}{
		"email": email,
		"exact": true,
	}, realmName, masterRealmToken.AccessToken)
	if err != nil {
		return err
	}

	var kcUser kcdto.KcUser
	if len(users) == 0 {
		kcUserInput["username"] = email
		kcUserInput["enabled"] = true
		kcUser, err = u.deps.KcUserService.PostUser(
			ctx,
			kcUserInput,
			realmName,
			masterRealmToken.AccessToken,
		)
		if err != nil {
			return err
		}
	} else {
		kcUser = users[0]
		if !kcUser.Enabled {
			err := u.deps.KcUserService.UpdateUser(
				ctx,
				map[string]interface{}{
					"enabled": true,
				},
				realmName,
				kcUser.ID,
				masterRealmToken.AccessToken,
			)
			if err != nil {
				return err
			}
		}
	}

	buyerAccount, err := u.deps.BuyerAccountRepo.TakeByConditions(ctx, map[string]interface{}{
		"kc_user_id": kcUser.ID,
	}, nil)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	otp := entity.Otp{}
	isNewUser := false

	if err := gormtx.GormTransaction(ctx, u.logger, u.deps.DB.GormDB(), nil, func(tx *gorm.DB) error {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			kcUserInfo, err := u.deps.KcUserService.GetUserByID(
				ctx,
				kcUser.ID,
				realmName,
				masterRealmToken.AccessToken,
			)
			if err != nil {
				return err
			}

			buyerAccount, err = u.deps.BuyerAccountRepo.CreateWithTx(ctx, gormtx.NewGormTx(tx), map[string]interface{}{
				"kc_user_id":   kcUser.ID,
				"mail_address": kcUserInfo.Email,
			})
			if err != nil {
				pgErr, isPgErr := err.(*pq.Error)
				if !isPgErr || pgErr.Code != "23505" {
					return err
				}
			}
			isNewUser = true
		}

		expiredAt := time.Now().AddDate(0, 0, 1).Format(baseUtils.FormatDateTimeSQL)
		otp, err = u.deps.OtpRepo.CreateWithTx(ctx, gormtx.NewGormTx(tx), map[string]interface{}{
			"kc_user_id": kcUser.ID,
			"token":      uuid.NewString(),
			"expired_at": expiredAt,
		})
		return err
	}); err != nil {
		return err
	}

	isRegisteredPassword, err := u.deps.KcUserService.CheckPasswordExist(
		ctx,
		realmName,
		kcUser.ID,
		masterRealmToken.AccessToken,
	)
	if err != nil {
		return err
	}

	redirectURI := u.deps.BuyerService.GenerateMailSignUpRedirectURI(
		isNewUser,
		buyerAccount.MailVerified,
		isRegisteredPassword,
		"https://es-krake.com",
		otp.Token,
		kcUser.ID,
		buyerAccount.MailAddress,
	)

	data := struct {
		Subject     string
		RedirectURI string
	}{
		Subject:     "Signup Buyer",
		RedirectURI: redirectURI,
	}
	pathTemplate := mailService.GetPathFileTemplate("/buyer/post_signup.html")
	body, err := mailService.ReadMailTemplate(pathTemplate, data)
	if err != nil {
		return err
	}

	sendEmailParams, err := mailService.NewSendEmailParams(
		mailService.SendEmailTypeHtml,
		netMail.Address{
			Address: "noreply@es-krake.com",
		},
		buyerAccount.MailAddress,
		"Signup Buyer",
		body,
		"",
	)
	if err != nil {
		return err
	}

	return u.deps.MailService.SendEmail(ctx, sendEmailParams)
}
