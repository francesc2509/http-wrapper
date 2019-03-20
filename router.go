package servergo

import (
	"context"
	"net/http"
)

// Router manages API's entry points
type Router struct {
	routes                  []*Route
	MethodNotAllowedHandler func(w http.ResponseWriter, message interface{})
	NotFoundHandler         func(w http.ResponseWriter, message interface{})
	UnauthorizedHandler     func(w http.ResponseWriter, message interface{})
	middleware              Middleware
}

// New Creates a new router
func New() *Router {
	router := &Router{
		routes:                  []*Route{},
		MethodNotAllowedHandler: methodNotAllowedHandler,
		NotFoundHandler:         notFoundHandler,
		UnauthorizedHandler:     unauthorizedHandler,
	}

	return router
}

// Middleware returns the middleware assigned to the router
func (router *Router) Middleware() Middleware {
	return router.middleware
}

// SetMiddleware assigns the specified middleware to the router
func (router *Router) SetMiddleware(middleware Middleware) {
	router.middleware = middleware
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	route := router.getRoute(r.URL.Path)

	if route == nil {
		router.NotFoundHandler(w, nil)
		return
	}

	if !route.methodAllowed(r.Method) {
		router.MethodNotAllowedHandler(w, nil)
		return
	}

	ctx := context.WithValue(context.Background(), paramsKey, &paramsWrapper{Params: route.Params()})
	r = r.WithContext(ctx)
	handler := route.Handler()

	if router.middleware != nil {
		handler = router.middleware(handler)
	}

	disableCors(handler)(w, r)
}

// HandleFunc push a new entry point to the router
func (router *Router) HandleFunc(path string, handler http.HandlerFunc, middlewares ...Middleware) *Route {
	route := &Route{}
	route.SetPath(path)
	route.SetHandler(handler)

	return router.AddRoute(route, middlewares...)
}

// HandleFileFunc push a path to filesystem
func (router *Router) HandleFileFunc(path string, handler http.HandlerFunc, middlewares ...Middleware) *Route {
	route := &Route{}
	route.filePath = true
	route.SetPath(path)
	route.SetHandler(handler)

	return router.AddRoute(route, middlewares...)
}

// AddGroup adds a new Group to the router
func (router *Router) AddGroup(path string, middlewares ...Middleware) *Group {
	group := &Group{}
	group.SetPath(path)
	group.SetRouter(router)
	group.SetMiddleware(chainMiddleware(middlewares...))

	return group
}

// AddRoute adds a route to the router
func (router *Router) AddRoute(route *Route, middlewares ...Middleware) *Route {
	route.SetMiddleware(chainMiddleware(middlewares...))
	router.routes = append(router.routes, route)

	return route
}

// Use adds middlewares to the route
func (router *Router) Use(middlewares ...Middleware) {
	use(router, middlewares...)
}

func (router *Router) getRoute(path string) *Route {
	for _, route := range router.routes {
		if route.match(path) {
			return route
		}
	}
	return nil
}

func disableCors(handler http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, Content-Length, Accept-Encoding")
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Max-Age", "86400")
			w.WriteHeader(http.StatusOK)
			return
		}
		handler(w, r)
	})
}
