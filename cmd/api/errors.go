package main

import (
	"fmt"
	"net/http"
)

func (app *application) logError(req *http.Request, err error) {
	app.logger.Print(err)
}

func (app *application) errorResponse(res http.ResponseWriter, req *http.Request, status int, message any) {
	env := envelope{"error": message}

	err := app.writeJson(res, status, env, nil)
	if err != nil {
		app.logError(req, err)
		res.WriteHeader(500)
	}
}

func (app *application) serverErrorResponse(res http.ResponseWriter, req *http.Request, err error) {
	app.logError(req, err)
	message := "The server encountered an error and could not process your request"
	app.errorResponse(res, req, http.StatusInternalServerError, message)
}

func (app *application) notFoundResponse(res http.ResponseWriter, req *http.Request) {
	message := "The requested resource could not be found"
	app.errorResponse(res, req, http.StatusNotFound, message)
}

func (app *application) methodNotAllowedResponse(res http.ResponseWriter, req *http.Request) {
	message := fmt.Sprintf("the %s method is not supported for this resource", req.Method)
	app.errorResponse(res, req, http.StatusMethodNotAllowed, message)
}

func (app *application) badRequestResponse(res http.ResponseWriter, req *http.Request, err error) {
	app.errorResponse(res, req, http.StatusBadRequest, err.Error())
}

func (app *application) failedValidationResponse(res http.ResponseWriter, req *http.Request, errors map[string]string) {
	app.errorResponse(res, req, http.StatusUnprocessableEntity, errors)
}

func (app *application) editConflictResponse(res http.ResponseWriter, req *http.Request) {
	message := "unable to update the record due to an edit conflict, please try again"
	app.errorResponse(res, req, http.StatusConflict, message)
}
