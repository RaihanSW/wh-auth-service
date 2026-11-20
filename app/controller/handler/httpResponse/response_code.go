package httpresponse

import (
	"context"
	"net/http"
	"strings"
	"wh-auth-service/app/model/response"
	"wh-auth-service/log"
	"wh-auth-service/util"
)

// service code
const (
	ServiceCodeMandiri  = "25"
	ServiceErrorDefault = "The system is currently experiencing issues. Please try again later"
)

// general codes
const (
	SuccessResp                   = "00"
	ErrCodeJsonDecode             = "01"
	ErrCodeInvalidRequest         = "02"
	ErrCodeInvalidService         = "03"
	ErrCodeInvalidExternalRequest = "04"
	ErrCodeInvalidAuth            = "05"
)

func GetErrorMessage(err string) string {
	errorMessage := strings.Split(err, ":")
	return strings.TrimSpace(errorMessage[len(errorMessage)-1])
}

func ErrorHandler(ctx context.Context, err error, status int, resp *response.GeneralResponse) {
	resp.SetError(GetErrorMessage(err.Error()))
	if status >= http.StatusInternalServerError {
		resp.SetError(http.StatusText(http.StatusInternalServerError))
		log.Logging.Error(util.RequestIDLog(ctx) + err.Error())
		return
	}
	log.Logging.Warn(util.RequestIDLog(ctx) + err.Error())
}
