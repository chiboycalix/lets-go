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
		Year    int32        `json:"year"`
		Runtime data.Runtime `json:"runtime"`
		Genres  []string     `json:"genres"`
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
		switch {
		case errors.Is(err, data.ErrorEditConflict):
			app.editConflictResponse(res, req)
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

func (app *application) showAllMovieHander(res http.ResponseWriter, req *http.Request) {
	var input struct {
		Title  string
		Genres []string
		data.Filters
	}

	v := validator.New()
	qs := req.URL.Query()

	input.Title = app.readString(qs, "title", "")
	input.Genres = app.readCSV(qs, "genres", []string{})
	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "pageSize", 10, v)
	input.Filters.Sort = app.readString(qs, "sort", "id")
	input.Filters.SortSafelist = []string{"id", "title", "year", "runtime", "-id", "-title", "-year", "-runtime"}

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(res, req, v.Errors)
		return
	}

	movies, metadata, err := app.models.Movies.GetAll(input.Title, input.Genres, input.Filters)
	if err != nil {
		app.serverErrorResponse(res, req, err)
		return
	}
	err = app.writeJson(res, http.StatusOK, envelope{"movies": movies, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(res, req, err)
	}
}

func (app *application) patchMovieHandler(res http.ResponseWriter, req *http.Request) {
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
		Title   *string       `json:"title"`
		Year    *int32        `json: "year"`
		Runtime *data.Runtime `json:runtime`
		Genres  []string      `json:genres`
	}

	err = app.readJson(res, req, &input)
	if err != nil {
		app.badRequestResponse(res, req, err)
		return
	}

	if input.Title != nil {
		movie.Title = *input.Title
	}
	if input.Year != nil {
		movie.Year = *input.Year
	}
	if input.Runtime != nil {
		movie.Runtime = *input.Runtime
	}
	if input.Genres != nil {
		movie.Genres = input.Genres
	}

	v := validator.New()
	if data.ValidateMovie(v, movie); !v.Valid() {
		app.failedValidationResponse(res, req, v.Errors)
		return
	}

	err = app.models.Movies.Update(movie)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrorEditConflict):
			app.editConflictResponse(res, req)
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
