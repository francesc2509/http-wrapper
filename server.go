package servergo

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"
)

// Param stores name and position of a url parameter
type Param struct {
	name  string
	start uint32
}

type paramsWrapper struct {
	Params []*Param
}

type contextKey int

const (
	paramsKey contextKey = iota
	routeKey
)

// Params returns a map with the values of current route's parameters
func Params(r *http.Request) map[string]string {
	params := make(map[string]string)
	ctxVars, ok := r.Context().Value(paramsKey).(*paramsWrapper)

	if ok && ctxVars != nil {
		splitURL := strings.Split(r.URL.Path[1:], "/")
		for _, param := range ctxVars.Params {
			params[param.name] = splitURL[param.start]
		}
	}
	return params
}

// arrContains returns true if the provided slice contains
// the seached value
func arrContains(arr interface{}, element interface{}) bool {
	switch reflect.TypeOf(arr).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(arr)
		for i := 0; i < s.Len(); i++ {
			if element == s.Index(i).Interface() {
				return true
			}
		}
	}
	return false
}

func methodNotAllowedHandler(w http.ResponseWriter, message interface{}) {
	w.WriteHeader(http.StatusMethodNotAllowed)
	fmt.Fprint(w)
}

func notFoundHandler(w http.ResponseWriter, message interface{}) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w)
}

func unauthorizedHandler(w http.ResponseWriter, message interface{}) {
	w.WriteHeader(http.StatusUnauthorized)
	fmt.Fprint(w)
}
