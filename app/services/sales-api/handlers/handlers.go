// Package handlers manges the different versions of the API
package handlers

import (
	"net/http"
	"os"

	"github.com/dimfeld/httptreemux/v5"
	"go.uber.org/zap"

	"github.com/qcbit/service/app/services/sales-api/handlers/v1/testgrp"
)

// APIMuxConfig contains all the mandatory systems required by handlers.
type APIMuxConfig struct {
	Shutdown chan os.Signal
	Log      *zap.SugaredLogger
}

// APIMux constructs a http.Handler with all application routes defined.
func APIMux(cfg APIMuxConfig) *httptreemux.ContextMux {
	mux := httptreemux.NewContextMux()

	mux.Handle(http.MethodGet, "/test", testgrp.Test)

	return mux
}
