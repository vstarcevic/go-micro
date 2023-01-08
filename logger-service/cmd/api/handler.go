package main

import (
	"log-service/data"
	"net/http"
)

type JSONpayload struct {
	Name string `json:"Name"`
	Data string `json:"Data"`
}

func (app *Config) WriteLog(w http.ResponseWriter, r *http.Request) {
	// read the json into a var
	var requestPayload JSONpayload
	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// insert data
	event := data.LogEntry{
		Name: requestPayload.Name,
		Data: requestPayload.Data,
	}

	err = app.Models.LogEntry.Insert(event)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	resp := jsonResponse{
		Error:   false,
		Message: "logged",
	}

	app.writeJSON(w, http.StatusAccepted, resp)

}
