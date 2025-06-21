package validator

import "regexp"

const (
	AllCharactersHalfWidth string = `^[\w\sｧ-ﾝﾞﾟ!@#$%^&*()-_=+{}|;:'",<.>/?]*$`
	AlphaNumericRegex      string = `^[a-zA-Z0-9]+$`
	Domain                 string = `^https:\/\/[a-z0-9-]+([\-\.]{1}[a-z0-9-]+)*\.[a-z]{2,5}(\/[\.a-z0-9_-]+)*$`
	AsciiCharacterPassword string = `^[\x00-\x7F]{8,32}$`
	ExceptSpecialCharacter string = `^[^¥\\ ]+$`
	AtLeastOneNumber       string = `[0-9]`
	AtLeastOneLowerCase    string = `[a-z]`
	AtLeastOneUpperCase    string = `[A-Z]`
	AtLeastOneSpecialChar  string = `[!@#~$%^&*()+|_]{1}`
	Url                    string = `^https?:\/\/[\w/:%#$@&?()~.=+-]+$`
	IDSnS                  string = `^[A-Za-z0-9-.]*$`
)

var (
	domain                 = regexp.MustCompile(Domain)
	asciiCharacterPassword = regexp.MustCompile(AsciiCharacterPassword)
	exceptSpecialCharacter = regexp.MustCompile(ExceptSpecialCharacter)
	atLeastOneNumber       = regexp.MustCompile(AtLeastOneNumber)
	atLeastOneLowerCase    = regexp.MustCompile(AtLeastOneLowerCase)
	atLeastOneUpperCase    = regexp.MustCompile(AtLeastOneUpperCase)
	regexUrl               = regexp.MustCompile(Url)
	regexIDSnS             = regexp.MustCompile(IDSnS)
)
