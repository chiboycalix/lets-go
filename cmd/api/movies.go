package main

import (
	"fmt"
	"net/http"
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

	fmt.Fprintf(res, "show the details of movie %d\n", id)
}
