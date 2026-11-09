package handler

import "net/http"

type AuthHandlerInterface interface {
	AuthRegister(w http.ResponseWriter, r *http.Request)
	AuthLogin(w http.ResponseWriter, r *http.Request)
	AuthLogout(w http.ResponseWriter, r *http.Request)
	UserGetByID(w http.ResponseWriter, r *http.Request)
	UserUpdate(w http.ResponseWriter, r *http.Request)
	UserDelete(w http.ResponseWriter, r *http.Request)
}
