// Package handlers manages the different versions of the API.
package handlers

import (
	"net/http"
	"os"

	"github.com/qcbit/services/app/services/sales-api/handlers/v1/testgrp"
	"github.com/qcbit/services/business/web/v1/mid"
	"github.com/qcbit/services/foundation/web"
	"go.uber.org/zap"
)

// APIMuxConfig contains all the mandatory systems required by handlers.
type APIMuxConfig struct {
	Shutdown chan os.Signal
	Log      *zap.SugaredLogger
}

// APIMux constructs a http.Handler with all application routes defined.
func APIMux(cfg APIMuxConfig) *web.App {
	app := web.NewApp(cfg.Shutdown, mid.Logger(cfg.Log), mid.Errors(cfg.Log), mid.Metrics(), mid.Panics())

	app.Handle(http.MethodGet, "/status", testgrp.Status)

	return app
}
