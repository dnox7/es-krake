package service

import (
	"strings"

	domainService "github.com/dpe27/es-krake/internal/domain/buyer/service"
	"github.com/dpe27/es-krake/pkg/log"
)

type buyerService struct {
	logger *log.Logger
}

func NewBuyerService() domainService.BuyerService {
	return &buyerService{
		logger: log.With("service", "buyer_service"),
	}
}

// GenerateMailSignUpRedirectURI implements service.BuyerService.
func (b *buyerService) GenerateMailSignUpRedirectURI(
	isNewUser, verified, isRegisteredPassword bool,
	clientRootURL, otpToken, userKCID, mailAddress string,
) string {
	if isNewUser || !verified || !isRegisteredPassword {
		return strings.Trim(clientRootURL, "/") + "/join/password" + "?token=" + otpToken + "&user_id=" + userKCID + "&mail_address=" + mailAddress
	}
	return strings.Trim(clientRootURL, "/") + "/login" + "?token=" + otpToken
}
