package web

// Middleware is a function designed to run some code before and/or after
// another Handler. It is desinged to remove boilerplat or other concerns not
// directed to any given Handler.
type Middleware func(Handler) Handler

// wrapMiddleware creates a new hanlder by wrapping middlewre around a final
// handler. The middlewares' Handlers will be executed by requests in the order
// they are provided.
func wrapMiddleware(mw []Middleware, handler Handler) Handler {
	// Loop backwards through the Middleware invoking each one. Replace the
	// handler with the new wrapped handler. Looping backwards ensures that the
	// first middleware of the slice is the first to be executed by requests.
	for i := len(mw)-1; i >= 0; i-- {
		h := mw[i]
		if h != nil {
			handler = h(handler)
		}
	}
	return handler
}