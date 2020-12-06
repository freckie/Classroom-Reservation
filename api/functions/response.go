package functions

import (
	"encoding/json"
	"log"
	"net/http"

	"classroom/models"
)

// ResponseOK make 200 response.
func ResponseOK(w http.ResponseWriter, msg string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(&models.Response{
		Status:  200,
		Message: msg,
		Data:    data,
	})
}

// ResponseError make error response.
func ResponseError(w http.ResponseWriter, errorCode int, errorMsg string) {
	log.Println("[Error] :", errorCode, errorMsg)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	w.WriteHeader(errorCode)
	json.NewEncoder(w).Encode(&models.Response{
		Status:  errorCode,
		Message: errorMsg,
		Data:    nil,
	})
}
