package main

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

type Route struct {
	Pattern string
	Handler HttpHandlerFunc
	Method  string
	regex   *regexp.Regexp
	params  []string
}

type Params map[string]string

const (
	MethodGet  = "GET"
	MethodPost = "POST"
)

type Router struct {
	routes   []Route
	NotFound HttpHandlerFunc
}

func NewRouter() *Router {
	return &Router{
		routes: make([]Route, 0),
	}
}

func (r *Router) add(method string, path string, h HttpHandlerFunc) error {
	route := &Route{
		Pattern: path,
		Handler: h,
		Method:  method,
	}

	regex, params, err := r.compilePattern(path)
	if err != nil {
		return fmt.Errorf("failed to compile pattern %s: %w", path, err)
	}

	route.params = params
	route.regex = regex
	r.routes = append(r.routes, *route)

	return nil
}

func (r *Router) match(method, rawPath string) (HttpHandlerFunc, Params, bool) {
	parsedURL, err := url.Parse(rawPath)
	if err != nil {
		return nil, nil, false
	}

	cleanPath := parsedURL.Path

	for _, route := range r.routes {
		if route.Method != method {
			continue
		}

		matches := route.regex.FindStringSubmatch(cleanPath)
		if matches == nil {
			continue
		}

		params := make(Params)
		for i, param := range route.params {
			if i+1 < len(matches) {

				// URL decode the parameter value
				decodedValue, err := url.QueryUnescape(matches[i+1]) // Using package function, not method
				if err != nil {
					decodedValue = matches[i+1]
				}
				params[param] = decodedValue
			}
		}

		return route.Handler, params, true
	}

	return nil, nil, false
}

func (r *Router) compilePattern(pattern string) (*regexp.Regexp, []string, error) {
	params := make([]string, 0)
	regexPattern := pattern

	// Colon parameters: /users/:id/posts/:postId
	colonRegex := regexp.MustCompile(`:([a-zA-Z_][a-zA-Z0-9_]*)`)
	matches := colonRegex.FindAllStringSubmatch(pattern, -1)

	for _, match := range matches {
		params = append(params, match[1])
		regexPattern = strings.Replace(regexPattern, match[0], `([^/]+)`, 1)
	}

	// Anchor the pattern
	regexPattern = "^" + regexPattern + "$"

	regex, err := regexp.Compile(regexPattern)
	if err != nil {
		return nil, nil, err
	}

	return regex, params, nil
}

func (r *Router) GET(path string, h HttpHandlerFunc) {
	r.add(MethodGet, path, h)
}

func (r *Router) ServeHTTP(rw ResponseWriter, req *Request) {
	h, params, ok := r.match(req.Method, req.Path)
	req.Params = params
	if !ok {
		if r.NotFound == nil {
			rw.WriteHeader(StatusNotFound)
			return
		}
		r.NotFound(rw, req)
		return
	}

	h(rw, req)
}

func (app *application) Routes() HttpHandler {
	router := NewRouter()

	router.GET("/v1/healthcheck", app.Healthcheck)
	router.GET("/v1/echo/:str", app.Echo)
	router.NotFound = app.NotFound

	return router
}
