package handler

import (
	"net/http"
	httpresponse "wh-auth-service/app/controller/handler/httpResponse"
	"wh-auth-service/app/model/response"
	"wh-auth-service/app/service"
)

type generalHandler struct {
	generalService service.GeneralServiceInterface
}

func NewGeneralHandler(generalService service.GeneralServiceInterface) GeneralHandlerInterface {
	return &generalHandler{
		generalService: generalService,
	}
}

func (h *generalHandler) Ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Pong!"))
}

func (h *generalHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Server is Running!"))
}

func (h *generalHandler) HealthCheckDB(w http.ResponseWriter, r *http.Request) {
	var resp response.GeneralResponse

	status, err := h.generalService.PingDB(r.Context())
	if err != nil {
		httpresponse.ErrorHandler(r.Context(), err, status, &resp)
		httpresponse.Response(w, r, resp, status)
		return
	}

	resp.IsSuccess = true
	resp.Message = "DB Ping Success"
	httpresponse.Response(w, r, resp, status)
}
