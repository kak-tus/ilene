package api

import (
	"net/http"

	oas "github.com/hypnoglow/oas2"
)

func (a *Type) middlewareRequestErrorHandler() oas.RequestErrorHandler {
	return func(w http.ResponseWriter, req *http.Request, err error) (resume bool) {
		switch err.(type) {
		case oas.ValidationError:
			e := err.(oas.ValidationError)
			a.respondClientErrors(w, e.Errors())
			return false
		default:
			a.log.Errorf("oas middleware: %s", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return false
		}
	}
}

func (a *Type) middlewareResponseErrorHandler() oas.ResponseErrorHandler {
	return func(w http.ResponseWriter, req *http.Request, err error) {
		switch err.(type) {
		case oas.ValidationError:
			e := err.(oas.ValidationError)
			a.respondClientErrors(w, e.Errors())
		default:
			a.log.Errorf("oas middleware: %s", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func (a *Type) respondClientErrors(w http.ResponseWriter, errs []error) {
	type (
		errorItem struct {
			Message string      `json:"message"`
			Field   string      `json:"field"`
			Value   interface{} `json:"value"`
		}
		payload struct {
			Errors []errorItem `json:"errors"`
		}
	)

	type fielder interface {
		Field() string
	}

	type valuer interface {
		Value() interface{}
	}

	p := payload{Errors: make([]errorItem, 0)}
	for _, e := range errs {
		item := errorItem{Message: e.Error()}
		if fe, ok := e.(fielder); ok {
			item.Field = fe.Field()
		}
		if ve, ok := e.(valuer); ok {
			item.Value = ve.Value()
		}
		p.Errors = append(p.Errors, item)
	}

	err := a.enc.NewEncoder(w).Encode(p)

	if err != nil {
		a.log.Error(err)
	}
}
