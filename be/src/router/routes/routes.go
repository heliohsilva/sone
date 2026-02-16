package routes

import (
	"api/src/middlewares"
	"net/http"

	"github.com/gorilla/mux"
)

type Route struct {
	URI      string
	Method   string
	Function func(http.ResponseWriter, *http.Request)
	NeedAuth bool
}

func Config(r *mux.Router) *mux.Router {
	var routes []Route
	routes = append(routes, userRoutes...)
	routes = append(routes, loginRoutes...)
	routes = append(routes, postRoutes...)

	for _, route := range routes {
		if route.NeedAuth {
			r.HandleFunc(route.URI, middlewares.Logger(middlewares.Auth(route.Function))).Methods(route.Method)
		} else {
			r.HandleFunc(route.URI, middlewares.Logger(route.Function)).Methods(route.Method)
		}
	}

	return r
}
