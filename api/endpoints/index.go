package endpoints

import (
	"classroom/functions"
	"classroom/models"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// GET /
func (e *Endpoints) IndexGet(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	// Struct for response
	resp := models.IndexResponse{}
	resp.WelcomeMessage = "Hello, Kyung Hee!"

	// Response with JSON
	functions.ResponseOK(w, "success", resp)
}
