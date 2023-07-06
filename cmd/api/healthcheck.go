package main

import (
	"net/http"
)

func (app *application) healthcheckHandler(res http.ResponseWriter, req *http.Request) {
	env := envelope{
		"status": "available",
		"system_info": map[string]string{
			"environment": app.configuration.env,
			"version":     version,
		},
	}

	err := app.writeJson(res, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(res, req, err)
	}
}
