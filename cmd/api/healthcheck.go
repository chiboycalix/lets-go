package main

import (
	"fmt"
	"net/http"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Status: Available")
	fmt.Fprintf(w, "environment: %s\n", app.configuration.env)
	fmt.Fprintf(w, "version: %s\n", version)
}
