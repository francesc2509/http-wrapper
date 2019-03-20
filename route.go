package servergo

import (
	"bytes"
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

// Route manages an API's entry point
type Route struct {
	filePath   bool
	group      *Group
	handler    http.HandlerFunc
	methods    []string
	middleware Middleware
	params     []*Param
	path       string
	url        string
}

// Handler returns the handler assigned to the route
func (route *Route) Handler() http.HandlerFunc {
	handler := route.handler

	if route.middleware != nil {
		handler = route.middleware(handler)
	}

	if group := route.group; group != nil && group.Middleware() != nil {
		handler = group.Middleware()(handler)
	}

	return handler
}

// SetHandler assigns the specified handler to the route
func (route *Route) SetHandler(handler http.HandlerFunc) {
	route.handler = handler
}

// Group returns the group who owns the route
func (route *Route) Group() *Group {
	return route.group
}

// SetGroup returns the group who owns the route
func (route *Route) SetGroup(group *Group) {
	route.group = group
}

// Methods sets Route's allowed http methods
func (route *Route) Methods(methods ...string) *Route {
	route.methods = methods
	return route
}

// Middleware returns the middleware assigned to the route
func (route *Route) Middleware() Middleware {
	return route.middleware
}

// SetMiddleware assigns the specified middleware to the route
func (route *Route) SetMiddleware(middleware Middleware) {
	route.middleware = middleware
}

// Params returns the params of the route
func (route *Route) Params() []*Param {
	return route.params
}

// Path returns route's path
func (route *Route) Path() string {
	return route.path
}

// SetPath assigns the specified path to the route
func (route *Route) SetPath(path string) {
	route.path = path

	err := route.createURLFromPath()
	if err != nil {
		panic(err)
	}
}

// URL returns route's url
func (route *Route) URL() string {
	return route.url
}

// Use adds middlewares to the route
func (route *Route) Use(middlewares ...Middleware) {
	use(route, middlewares...)
}

func (route *Route) match(url string) bool {
	matched, _ := regexp.MatchString(route.url, url)
	return matched
}

func (route *Route) methodAllowed(method string) bool {
	return len(route.methods) == 0 || arrContains(route.methods, method)
}

func (route *Route) createURLFromPath() error {
	var buffer bytes.Buffer

	buffer.WriteString("^")
	if strings.Contains(route.path, ":") {
		splitURL := strings.Split(route.path[1:], "/")

		for i, str := range splitURL {
			buffer.WriteString("/")
			paramSepIndex := strings.Index(str, ":")
			if paramSepIndex == 0 {
				regexIniPos := strings.Index(str, "(")
				if regexIniPos > -1 {
					if paramSepIndex == regexIniPos-1 {
						return fmt.Errorf("Incorrect params: %s", route.path)
					}

					route.params = append(route.params, &Param{
						name:  str[paramSepIndex+1 : regexIniPos],
						start: uint32(i),
					})
					buffer.WriteString(str[regexIniPos:])
				} else {
					buffer.WriteString(".+")
					route.params = append(route.params, &Param{
						name:  str[paramSepIndex+1:],
						start: uint32(i),
					})
				}
			} else {
				buffer.WriteString(str)
			}
		}
	} else {
		buffer.WriteString(route.path)
	}

	url := buffer.String()
	buffer.Reset()

	url = strings.TrimSuffix(url, "/")
	buffer.WriteString(url)
	if route.filePath {
		if url != "" {
			buffer.WriteString("/")
		}
		buffer.WriteString(".*")
	}
	buffer.WriteString("(/)?$")

	route.url = buffer.String()
	fmt.Println(route.url)
	return nil
}
