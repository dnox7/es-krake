package utils

import "time"

const (
	ErrInputFail    = "ERR001"
	ErrInputRequire = "ERR002"
	ErrEmailFail    = "ERR003"
	ErrPasswordFail = "ERROO4"

	DefaultPageNo   = 1
	DefaultPageSize = 30

	FormatYearISO     = "2006"
	FormatDateISO     = time.DateOnly
	FormatTimeHHMM    = "15:04"
	FormatTimeHHMMSS  = time.TimeOnly
	FormatDateTimeISO = time.RFC3339
	FormatDateTimeSQL = time.DateTime
	FormatDateCompact = "20060102150405"

	ProdEnv    = "prod"
	DevEnv     = "dev"
	TestingEnv = "testing"

	ErrorInternalServer    = "Internal server error"
	ErrorLogRequestBody    = "failed to log request body"
	ErrorCloseResponseBody = "failed to close response body"
	ErrorCloseRows         = "failed to close Rows"
	ErrorCloseReader       = "failed to close reader"
	ErrorCloseWriter       = "failed to close writer"
	ErrorCloseFile         = "failed to close file"
	ErrorCloseSftp         = "failed to close sftp"
	ErrorParseUrl          = "failed to parse url"
	ErrorReadBody          = "failed to read body"
	ErrorDecodeBody        = "failed to decode body"
	ErrorMapToStruct       = "failed to map to struct"
	ErrorGetTx             = "failed to get Tx"
	ErrorGetSpec           = "failed to get specification"
	ErrorCreateReq         = "failed to create new HTTP request"
	ErrorMarshalFailed     = "failed to marshal object to JSON"
)
