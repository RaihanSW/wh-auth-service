package handler

import "net/http"

type RoleHandlerInterface interface {
	RoleCreate(w http.ResponseWriter, r *http.Request)
	RoleGetAll(w http.ResponseWriter, r *http.Request)
	RoleGetByID(w http.ResponseWriter, r *http.Request)
	RoleUpdate(w http.ResponseWriter, r *http.Request)
	RoleDelete(w http.ResponseWriter, r *http.Request)
}
