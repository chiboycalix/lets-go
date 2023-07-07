package main

import (
	"net/http"
	"time"

	"github.com/chiboycalix/go-further/internal/data"
)

func (app *application) createMovieHandler(res http.ResponseWriter, req *http.Request) {
	var input struct {
		Title   string       `json:"title"`
		Year    int32        `json: "year"`
		Runtime data.Runtime `json:runtime`
		Genres  []string     `json:genres`
	}

	err := app.readJson(res, req, &input)
	if err != nil {
		app.errorResponse(res, req, http.StatusBadRequest, err.Error())
		return
	}

	env := envelope{"movie": input}

	err = app.writeJson(res, http.StatusCreated, env, nil)
	if err != nil {
		app.badRequestResponse(res, req, err)
	}
}

func (app *application) showMovieHandler(res http.ResponseWriter, req *http.Request) {
	id, err := app.readIdFromParams(req)
	if err != nil {
		app.notFoundResponse(res, req)
		return
	}

	movie := data.Movie{
		ID:        id,
		CreatedAt: time.Now(),
		Title:     "Casablanca",
		Runtime:   102,
		Genres:    []string{"drama", "romance", "war"},
		Version:   1,
	}

	err = app.writeJson(res, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
		app.serverErrorResponse(res, req, err)
	}
}
