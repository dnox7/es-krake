package dto

type BaseErrorResponse struct {
	Error *ErrorResponse `json:"error"`
}

type ErrorResponse struct {
	Message   interface{} `json:"message"`
	Details   interface{} `json:"details"`
	DebugInfo interface{} `json:"debug_information"`
}
