package validator

import (
	"time"

	"github.com/dpe27/es-krake/pkg/utils"
	"github.com/xeipuuv/gojsonschema"
)

type NonStandardDatetimeFormatChecker struct{}

func (NonStandardDatetimeFormatChecker) IsFormat(input interface{}) bool {
	if inputStr, ok := input.(string); ok {
		_, err := time.Parse(utils.FormatDateTimeSQL, inputStr)
		return err == nil
	}
	return false
}

type StrongPassswordChecker struct{}

func (StrongPassswordChecker) IsFormat(input interface{}) bool {
	if inputStr, ok := input.(string); ok {
		return len(inputStr) >= 12 &&
			asciiCharacterPassword.MatchString(inputStr) &&
			atLeastOneNumber.MatchString(inputStr) &&
			atLeastOneLowerCase.MatchString(inputStr) &&
			atLeastOneUpperCase.MatchString(inputStr) &&
			exceptSpecialCharacter.MatchString(inputStr)
	}
	return false
}

type DomainChecker struct{}

func (DomainChecker) IsFormat(input interface{}) bool {
	if inputStr, ok := input.(string); ok {
		return domain.MatchString(inputStr)
	}
	return false
}

type UrlChecker struct{}

func (UrlChecker) IsFormat(input interface{}) bool {
	if inputStr, ok := input.(string); ok {
		return regexUrl.MatchString(inputStr)
	}
	return false
}

type IDSnSChecker struct{}

func (IDSnSChecker) IsFormat(input interface{}) bool {
	if inputStr, ok := input.(string); ok {
		return regexIDSnS.MatchString(inputStr)
	}
	return false
}

type MaxLengthChecker struct{}

func (MaxLengthChecker) IsFormat(input interface{}) bool {
	if inputStr, ok := input.(string); ok {
		return len([]rune(inputStr)) <= 50
	}
	return false
}

func NewMaxLengthError(ctx *gojsonschema.JsonContext, val interface{}, details gojsonschema.ErrorDetails) *gojsonschema.ResultErrorFields {
	err := &gojsonschema.ResultErrorFields{}
	err.SetContext(ctx)
	err.SetType("max_length_byte")
	err.SetDescriptionFormat(utils.ErrorInputByteLimit)
	err.SetValue(val)
	err.SetDetails(details)
	return err
}
