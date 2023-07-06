package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/chiboycalix/go-further/internal/data"
)

func (app *application) createMovieHandler(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(res, "create a new movie")
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
