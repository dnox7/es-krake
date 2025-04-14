package utils

const (
	ErrInputFail    = "ERR001"
	ErrInputRequire = "ERR002"
	ErrEmailFail    = "ERR003"
	ErrPasswordFail = "ERROO4"

	ErrInternalServerMsg     = "Internal server error"
	ErrQueryBuilderFailedMsg = "The query builder failed"

	DefaultPageNo   = 1
	DefaultPageSize = 30

	FormatYearISO     = "2006"
	FormatDateISO     = "2006-01-02"
	FormatTimeHHMM    = "15:04"
	FormatTimeHHMMSS  = "15:04:05"
	FormatDateTimeISO = "2006-01-02T15:04:05Z07:00"
	FormatDateTimeSQL = "2006-01-02 15:04:05"
	FormatDateCompact = "20060102150405"

	ProdEnv    = "prod"
	DevEnv     = "dev"
	TestingEnv = "testing"
)
