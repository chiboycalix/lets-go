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
	id, err := app.readIdFromParams(res, req)
	if err != nil {
		http.NotFound(res, req)
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

	er := app.writeJson(res, http.StatusOK, movie, nil)
	if er != nil {
		app.logger.Print(er)
		http.Error(res, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
	}
}
