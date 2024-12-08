package main

import (
	"net/http"
)

type JsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	var requestData map[string]interface{}

	err := app.readJson(w, r, &requestData)
	if err != nil {
		app.errorJson(w, err)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "Broker verileri başarıyla aldı",
		Data:    requestData,
	}

	err = app.writeJson(w, http.StatusAccepted, payload)
	if err != nil {
		app.errorJson(w, err)
	}
}
