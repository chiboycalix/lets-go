package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/chiboycalix/go-further/internal/data"
	"github.com/chiboycalix/go-further/internal/validator"
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

	movie := data.Movie{
		Title:   input.Title,
		Year:    input.Year,
		Runtime: input.Runtime,
		Genres:  input.Genres,
	}
	v := validator.New()

	if data.ValidateMovie(v, &movie); !v.Valid() {
		app.failedValidationResponse(res, req, v.Errors)
		return
	}

	err = app.models.Movies.Insert(&movie)
	if err != nil {
		fmt.Println(err, "erro")
		app.serverErrorResponse(res, req, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/movies/%d", movie.ID))

	env := envelope{"movie": movie}

	err = app.writeJson(res, http.StatusCreated, env, headers)
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
	movie, err := app.models.Movies.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrorRecordNotFound):
			app.notFoundResponse(res, req)
		default:
			app.serverErrorResponse(res, req, err)
		}
		return
	}

	err = app.writeJson(res, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
		app.serverErrorResponse(res, req, err)
	}
}

func (app *application) updateMovieHandler(res http.ResponseWriter, req *http.Request) {
	id, err := app.readIdFromParams(req)
	if err != nil {
		app.notFoundResponse(res, req)
		return
	}
	movie, err := app.models.Movies.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrorRecordNotFound):
			app.notFoundResponse(res, req)
		default:
			app.serverErrorResponse(res, req, err)
		}
		return
	}

	var input struct {
		Title   string       `json:"title"`
		Year    int32        `json: "year"`
		Runtime data.Runtime `json:runtime`
		Genres  []string     `json:genres`
	}

	err = app.readJson(res, req, &input)
	if err != nil {
		app.badRequestResponse(res, req, err)
		return
	}

	movie.Title = input.Title
	movie.Year = input.Year
	movie.Runtime = input.Runtime
	movie.Genres = input.Genres

	v := validator.New()
	if data.ValidateMovie(v, movie); !v.Valid() {
		app.failedValidationResponse(res, req, v.Errors)
		return
	}

	err = app.models.Movies.Update(movie)
	if err != nil {
		app.serverErrorResponse(res, req, err)
		return
	}

	err = app.writeJson(res, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
		app.serverErrorResponse(res, req, err)
	}
}

func (app *application) deleteMovieHandler(res http.ResponseWriter, req *http.Request) {
	id, err := app.readIdFromParams(req)
	if err != nil {
		app.notFoundResponse(res, req)
		return
	}

	err = app.models.Movies.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrorRecordNotFound):
			app.notFoundResponse(res, req)
		default:
			app.serverErrorResponse(res, req, err)
		}
		return
	}

	err = app.writeJson(res, http.StatusOK, envelope{"message": "movie successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(res, req, err)
	}
}
