package handler

import (
	httpresponse "dc-transaction-service/app/controller/handler/httpResponse"
	"dc-transaction-service/app/model/request"
	"dc-transaction-service/app/model/response"
	"dc-transaction-service/app/service"
	"dc-transaction-service/util"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type authHandler struct {
	authService service.AuthServiceInterface
}

func NewAuthHandler(authService service.AuthServiceInterface) AuthHandlerInterface {
	return &authHandler{
		authService: authService,
	}
}

// AuthRegister handles user registration
// @Summary register user
// @Description register user
// @Tags auth
// @Produce json
// @Param	RegisterRequest body  request.RegisterRequest true "Register Request Body"
// @Param	X-Version	header	string	false "app version"
// @Param	X-DEVICE-UIID	header	string	false "device uiid"
// @Success 201 {object}  httpresponse.SetResponse
// @Failure 400 {object}  httpresponse.SetResponse
// @Router /account/register [post]
func (h *authHandler) AuthRegister(w http.ResponseWriter, r *http.Request) {

	var resp response.AuthRegisterResponse
	var req request.RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpresponse.ErrorHandler(r.Context(), err, http.StatusBadRequest, &resp.GeneralResponse)
		httpresponse.Response(w, r, resp, http.StatusBadRequest)
		return
	}

	if status, err := util.ValidateRequest(r.Context(), req); err != nil {
		httpresponse.ErrorHandler(r.Context(), err, status, &resp.GeneralResponse)
		httpresponse.Response(w, r, resp, status)
		return
	}

	uid, status, err := h.authService.AuthRegister(r.Context(), req)
	if err != nil {
		httpresponse.ErrorHandler(r.Context(), err, status, &resp.GeneralResponse)
		httpresponse.Response(w, r, resp, status)
		return
	}

	resp.IsSuccess = true
	resp.UID = uid
	httpresponse.Response(w, r, resp, status)
}

// AuthLogin handles user login
// @Summary login user
// @Description login user, returns JWT token
// @Tags auth
// @Produce json
// @Param	LoginRequest body  request.LoginRequest true "Login Request Body"
// @Param	X-Version	header	string	false "app version"
// @Param	X-DEVICE-UIID	header	string	false "device uiid"
// @Success 200 {object}  httpresponse.SetResponse
// @Failure 400 {object}  httpresponse.SetResponse
// @Failure 401 {object}  httpresponse.SetResponse
// @Router /account/login [post]
func (h *authHandler) AuthLogin(w http.ResponseWriter, r *http.Request) {

	var resp response.AuthLoginResponse
	var req request.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpresponse.ErrorHandler(r.Context(), err, http.StatusBadRequest, &resp.GeneralResponse)
		httpresponse.Response(w, r, resp, http.StatusBadRequest)
		return
	}

	if status, err := util.ValidateRequest(r.Context(), req); err != nil {
		httpresponse.ErrorHandler(r.Context(), err, status, &resp.GeneralResponse)
		httpresponse.Response(w, r, resp, status)
		return
	}

	user, token, status, err := h.authService.AuthLogin(r.Context(), req)
	if err != nil {
		httpresponse.ErrorHandler(r.Context(), err, status, &resp.GeneralResponse)
		httpresponse.Response(w, r, resp, status)
		return
	}

	resp.IsSuccess = true
	resp.Token = token.Token
	resp.RefreshToken = token.RefreshToken
	resp.Data = user
	httpresponse.Response(w, r, resp, status)
}

// AuthLogout handles user logout
// @Summary logout user
// @Description logout user (stateless — client discards the token)
// @Tags auth
// @Produce json
// @Param	Authorization	header	string	false "Bearer token"
// @Success 200 {object}  httpresponse.SetResponse
// @Router /account/logout [post]
func (h *authHandler) AuthLogout(w http.ResponseWriter, r *http.Request) {

	var resp response.AuthLogoutResponse

	status, err := h.authService.AuthLogout(r.Context())
	if err != nil {
		httpresponse.ErrorHandler(r.Context(), err, status, &resp.GeneralResponse)
		httpresponse.Response(w, r, resp, status)
		return
	}

	resp.IsSuccess = true
	resp.Message = "Logout success"
	httpresponse.Response(w, r, resp, status)
}

// UserGetByID handles get user by uid
// @Summary get user by uid
// @Description get user by uid
// @Tags auth
// @Produce json
// @Param	uid	path	string	true "User UID"
// @Param	Authorization	header	string	false "Bearer token"
// @Success 200 {object}  httpresponse.SetResponse
// @Failure 404 {object}  httpresponse.SetResponse
// @Router /account/{uid} [get]
func (h *authHandler) UserGetByID(w http.ResponseWriter, r *http.Request) {

	var resp response.UserDataResponse

	uid := chi.URLParam(r, "uid")

	data, status, err := h.authService.UserGetByID(r.Context(), uid)
	if err != nil {
		httpresponse.ErrorHandler(r.Context(), err, status, &resp.GeneralResponse)
		httpresponse.Response(w, r, resp, status)
		return
	}

	resp.IsSuccess = true
	resp.Data = data
	httpresponse.Response(w, r, resp, status)
}

// UserUpdate handles update user by uid
// @Summary update user by uid
// @Description update user by uid
// @Tags auth
// @Produce json
// @Param	uid	path	string	true "User UID"
// @Param	UserUpdateRequest body  request.UserUpdateRequest true "Update User Request Body"
// @Param	Authorization	header	string	false "Bearer token"
// @Success 200 {object}  httpresponse.SetResponse
// @Failure 400 {object}  httpresponse.SetResponse
// @Failure 404 {object}  httpresponse.SetResponse
// @Router /account/{uid} [put]
func (h *authHandler) UserUpdate(w http.ResponseWriter, r *http.Request) {

	var resp response.UserResponse
	var req request.UserUpdateRequest

	uid := chi.URLParam(r, "uid")

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpresponse.ErrorHandler(r.Context(), err, http.StatusBadRequest, &resp.GeneralResponse)
		httpresponse.Response(w, r, resp, http.StatusBadRequest)
		return
	}

	if status, err := util.ValidateRequest(r.Context(), req); err != nil {
		httpresponse.ErrorHandler(r.Context(), err, status, &resp.GeneralResponse)
		httpresponse.Response(w, r, resp, status)
		return
	}

	status, err := h.authService.UserUpdate(r.Context(), uid, req)
	if err != nil {
		httpresponse.ErrorHandler(r.Context(), err, status, &resp.GeneralResponse)
		httpresponse.Response(w, r, resp, status)
		return
	}

	resp.IsSuccess = true
	httpresponse.Response(w, r, resp, status)
}

// UserDelete handles delete user by uid
// @Summary delete user by uid
// @Description delete user by uid
// @Tags auth
// @Produce json
// @Param	uid	path	string	true "User UID"
// @Param	Authorization	header	string	false "Bearer token"
// @Success 200 {object}  httpresponse.SetResponse
// @Failure 404 {object}  httpresponse.SetResponse
// @Router /account/{uid} [delete]
func (h *authHandler) UserDelete(w http.ResponseWriter, r *http.Request) {

	var resp response.UserResponse

	uid := chi.URLParam(r, "uid")

	status, err := h.authService.UserDelete(r.Context(), uid)
	if err != nil {
		httpresponse.ErrorHandler(r.Context(), err, status, &resp.GeneralResponse)
		httpresponse.Response(w, r, resp, status)
		return
	}

	resp.IsSuccess = true
	httpresponse.Response(w, r, resp, status)
}
