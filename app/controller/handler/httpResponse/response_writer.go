package httpresponse

import (
	"encoding/json"
	"net/http"
	"time"
)

type SetResponse struct {
	Warehouse WarehouseResponse `json:"warehouse"`
}

type WarehouseResponse struct {
	Status       string      `json:"status"`
	Results      interface{} `json:"results,omitempty"`
	AccessedTime string      `json:"accessedTime"`
}

func Response(w http.ResponseWriter, r *http.Request, results interface{}, code int) {
	warehouseResponse := WarehouseResponse{
		Status:       http.StatusText(code),
		AccessedTime: time.Now().Format("02-01-2006 15:04:05"),
		Results:      results}
	setResponse := SetResponse{Warehouse: warehouseResponse}
	response, _ := json.Marshal(setResponse)
	w.Header().Set("Content-Type", "Application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func ResponseWithoutTemplate(w http.ResponseWriter, r *http.Request, data interface{}, code int) {
	response, _ := json.Marshal(data)
	w.Header().Set("Content-Type", "Application/json")
	w.WriteHeader(code)
	w.Write(response)
}
