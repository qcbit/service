// Package web contains a small web framework for extension.
package web

import (
	"context"
	"net/http"
	"os"

	"github.com/dimfeld/httptreemux/v5"
)

// A Handler is a type that handles an http request within our own little mini framework.
type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

// App is the entrypoint into our application and what configures
// our context object for each of our http handlers.
// Feel free to add any configuration data/logic to this App struct.
type App struct {
	*httptreemux.ContextMux
	shutdown chan os.Signal
	mw       []Middleware
}

// NewApp creates an App value that handle a set of routes for the application.
func NewApp(shutdown chan os.Signal, mw ...Middleware) *App {
	return &App{
		ContextMux: httptreemux.NewContextMux(),
		shutdown:   shutdown,
		mw:         mw,
	}
}

// Handle sets a handler function for a given HTTP method and path pair to the application server mux.
func (a *App) Handle(method string, path string, handler Handler, mw ...Middleware) {
	handler = wrapMiddleware(mw, handler)
	handler = wrapMiddleware(a.mw, handler)

	h := func(w http.ResponseWriter, r *http.Request) {

		// ADD CODE HERE: LOG

		if err := handler(r.Context(), w, r); err != nil {
			// ADD CODE HERE: ERROR HANDLING
			return
		}

		// ADD CODE HERE: LOG
	}

	a.ContextMux.Handle(method, path, h)

}
