package httpclient

type (
	ReqOpt struct{}

	reqOptSetter func(*reqOpt)

	reqOpt struct {
		canLog                      bool
		canLogRequestBody           bool
		canLogResponseBody          bool
		canLogRequestBodyOnlyError  bool
		canLogResponseBodyOnlyError bool
		loggedRequestBody           []string
		loggedResponseBody          []string
		markedQueryParamKeys        []string
		retryTimes                  uint
	}
)

func (ReqOpt) LoggedRequestBody(body []string) reqOptSetter {
	return func(ro *reqOpt) {
		ro.canLog = true
		ro.canLogRequestBody = true
		ro.loggedRequestBody = body
	}
}

func (ReqOpt) LoggedResponseBody(body []string) reqOptSetter {
	return func(ro *reqOpt) {
		ro.canLog = true
		ro.canLogResponseBody = true
		ro.loggedResponseBody = body
	}
}

func (ReqOpt) LoggedRequestBodyOnlyError(body []string) reqOptSetter {
	return func(ro *reqOpt) {
		ro.canLog = true
		ro.canLogRequestBodyOnlyError = true
		ro.loggedRequestBody = body
	}
}

func (ReqOpt) LoggedResponseBodyOnlyError(body []string) reqOptSetter {
	return func(ro *reqOpt) {
		ro.canLog = true
		ro.canLogResponseBodyOnlyError = true
		ro.loggedResponseBody = body
	}
}

func (ReqOpt) MarkedQueryParamKeys(keys []string) reqOptSetter {
	return func(ro *reqOpt) {
		ro.canLog = true
		ro.markedQueryParamKeys = keys
	}
}

func (ReqOpt) RetryTimes(retries uint) reqOptSetter {
	return func(ro *reqOpt) {
		ro.retryTimes = retries
	}
}
