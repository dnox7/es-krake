package service

type BuyerService interface {
	GenerateMailSignUpRedirectURI(
		isNewUser, verified, isRegisteredPassword bool,
		clientRootURL, otpToken, userKCID, mailAddress string,
	) string
}
