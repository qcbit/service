package testgrp

import (
	"context"
	"encoding/json"
	"net/http"
)

// Status represents a test handler.
func Status(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	status := struct {
		Status string
	}{
		Status: "OK",
	}
	json.NewEncoder(w).Encode(status)
}