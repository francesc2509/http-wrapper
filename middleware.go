package servergo

import (
	"net/http"
)

// Middleware is function that must be executed before
// the corresponding HTTP handler
type Middleware func(next http.HandlerFunc) http.HandlerFunc

// hasMiddleware interface defines the methods
// required to have the possibility of adding
// a middleware function
type hasMiddleware interface {
	Use(middlewares ...Middleware)
	Middleware() Middleware
	SetMiddleware(middleware Middleware)
}

// chainMiddleware provides syntactic sugar to create a new middleware
// which will be the result of chaining the ones received as parameters.
func chainMiddleware(mw ...Middleware) Middleware {
	if mw != nil && len(mw) > 0 {
		return func(final http.HandlerFunc) http.HandlerFunc {
			return func(w http.ResponseWriter, r *http.Request) {
				last := final
				for i := len(mw) - 1; i >= 0; i-- {
					last = mw[i](last)
				}
				last(w, r)
			}
		}
	}
	return nil
}

// use chains the middlewares provided by param
func use(h hasMiddleware, middlewares ...Middleware) {
	var m []Middleware
	if h.Middleware() != nil {
		m = append(m, h.Middleware())
	}

	m = append(m, middlewares...)
	h.SetMiddleware(chainMiddleware(m...))
}
