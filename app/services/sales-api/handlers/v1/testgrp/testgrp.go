package testgrp

import (
	"context"
	"net/http"
	"errors"
	"math/rand"

	// v1web "github.com/qcbit/services/business/web/v1"
	"github.com/qcbit/services/foundation/web"
)

// Status represents a test handler.
func Status(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	if n := rand.Intn(100); n%2 == 0 {
		// return v1web.NewRequestError(errors.New("trusted error"), http.StatusBadRequest)
		return errors.New("non-trusted error")
	}
	status := struct {
		Status string
	}{
		Status: "OK",
	}
	return web.Respond(ctx, w, status, http.StatusOK)
}