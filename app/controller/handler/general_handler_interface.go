package handler

import "net/http"

type GeneralHandlerInterface interface {
	Ping(w http.ResponseWriter, r *http.Request)
	HealthCheck(w http.ResponseWriter, r *http.Request)
	HealthCheckDB(w http.ResponseWriter, r *http.Request)
}
