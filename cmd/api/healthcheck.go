package main

import (
	"net/http"
)

func (app *application) healthcheckHandler(res http.ResponseWriter, req *http.Request) {
	data := map[string]string{
		"status":      "available",
		"environment": app.configuration.env,
		"version":     version,
	}

	err := app.writeJson(res, http.StatusOK, data, nil)
	if err != nil {
		app.logger.Print(err)
		http.Error(res, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
	}
}
