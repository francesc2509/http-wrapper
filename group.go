package servergo

import (
	"fmt"
	"net/http"
)

// Group groups actions
type Group struct {
	path       string
	middleware Middleware
	router     *Router
}

// Path returns the path of the group
func (group *Group) Path() string {
	return group.path
}

// SetPath sets the path of the group
func (group *Group) SetPath(path string) {
	group.path = path
}

// HandleFunc adds a new route to the group
func (group *Group) HandleFunc(path string, handler http.HandlerFunc, middlewares ...Middleware) *Route {
	route := &Route{}
	route.SetGroup(group)
	route.SetHandler(handler)
	route.SetPath(fmt.Sprint(group.path, path))

	group.router.AddRoute(route, middlewares...)
	return route
}

// Middleware returns group's middleware
func (group *Group) Middleware() Middleware {
	return group.middleware
}

// SetMiddleware sets group's middleware
func (group *Group) SetMiddleware(middleware Middleware) {
	group.middleware = middleware
}

// Router returns the router that owns the group
func (group *Group) Router() *Router {
	return group.router
}

// SetRouter sets the router owner of the group
func (group *Group) SetRouter(router *Router) {
	group.router = router
}

// Use adds middlewares to the group
func (group *Group) Use(middlewares ...Middleware) {
	use(group, middlewares...)
}
