package main

import (
	"fmt"
	"net/http"

	"golang.org/x/time/rate"
)

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				res.Header().Set("Connection", "close")
				app.serverErrorResponse(res, req, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(res, req)
	})
}

func (app *application) rateLimit(next http.Handler) http.Handler {
	limit := rate.NewLimiter(2, 4)
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if !limit.Allow() {
			app.rateLimitExceededResponse(res, req)
			return
		}
		next.ServeHTTP(res, req)
	})
}
