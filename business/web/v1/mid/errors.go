package mid

import (
	"context"
	"net/http"

	"github.com/qcbit/services/foundation/web"
	v1web "github.com/qcbit/services/business/web/v1"
	"go.uber.org/zap"
)

// Errors handle errors coming out of the call chain. It detects normal
// application errors which are used to respond to the client in a uniform way.
// Unexpected errors (status >= 500) are logged.
func Errors(log *zap.SugaredLogger) web.Middleware {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			if err := handler(ctx, w, r); err != nil {
				log.Errorw("ERROR", "trace_id", web.GetTraceID(ctx), "message", err)

				var er v1web.ErrorResponse
				var status int

				switch {
				case v1web.IsRequestError(err):
					reqErr := v1web.GetRequestError(err)
					er = v1web.ErrorResponse{
						Error: reqErr.Error(),
					}
					status = reqErr.Status
				default:
					er = v1web.ErrorResponse{
						Error: http.StatusText(http.StatusInternalServerError),
					}
					status = http.StatusInternalServerError
				}

				if err := web.Respond(ctx, w, er, status); err != nil {
					return err
				}

				// If we receive the shutdown err we need to return it
				// back to the base handler to shutdown the service
				if web.IsShutdown(err) {
					return err
				}
			}
			return nil
		}

		return h
	}
	return m
}