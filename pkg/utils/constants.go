package utils

const (
	ErrInputFail    = "ERR001"
	ErrInputRequire = "ERR002"
	ErrEmailFail    = "ERR003"
	ErrPasswordFail = "ERROO4"

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

	HeaderAuthorization = "Authorization"
	HeaderContentType   = "Content-Type"

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
)
