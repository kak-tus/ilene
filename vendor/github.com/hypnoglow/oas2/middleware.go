package oas

import (
	"fmt"
	"net/http"
	"strings"
)

// Middleware describes a middleware that can be applied to a http.handler.
type Middleware func(next http.Handler) http.Handler

// RequestErrorHandler is a function that handles an error occurred in
// middleware while working with request. It is the library user responsibility
// to implement this to handle various errors that can occur during middleware
// work. This errors can include request validation errors, json encoding errors
// and other. Also, user must return proper boolean value that indicates if
// the request should continue or it should be stopped (basically, call "next"
// or not).
type RequestErrorHandler func(w http.ResponseWriter, req *http.Request, err error) (resume bool)

// ResponseErrorHandler is a function that handles an error occurred in
// middleware while working with response. It is the library user responsibility
// to implement this to handle various errors that can occur on middleware work.
// This errors can include response validation errors, json serialization errors
// and others.
type ResponseErrorHandler func(w http.ResponseWriter, req *http.Request, err error)

// JSONError occurs on json encoding or decoding.
// It can happen both in request and response validation.
type JSONError struct {
	error
}

// ValidationError occurs on request or response validation.
type ValidationError struct {
	error
	errs []error
}

// Error implements error.
func (ve ValidationError) Error() string {
	var s []string
	for _, err := range ve.errs {
		s = append(s, err.Error())
	}
	return fmt.Sprintf("%s: %s", ve.error, strings.Join(s, ", "))
}

// Errors returns validation errors.
func (ve ValidationError) Errors() []error {
	return ve.errs
}
