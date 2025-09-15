package main

type Routes map[string]HttpHandlerFunc

const (
	MethodGet  = "GET"
	MethodPost = "POST"
)

type Router struct {
	routes   Routes
	NotFound HttpHandlerFunc
}

func NewRouter() *Router {
	return &Router{
		routes: make(Routes),
	}
}

func (r *Router) HandlerFunc(method string, path string, h HttpHandlerFunc) {
	r.routes[method+path] = h
}

func (r *Router) ServeHTTP(rw ResponseWriter, req *Request) {
	h, ok := r.routes[req.Method+req.Path]
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

	router.HandlerFunc(MethodGet, "/v1/healthcheck", app.Healthcheck)
	router.NotFound = app.NotFound

	return router
}
