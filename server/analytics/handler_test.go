package analytics_test

import (
	"net/http"
	"server/analytics"
	"testing"

	"github.com/gorilla/mux"
)

func TestHandler_RegisterRoutes(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		service analytics.Service
		// Named input parameters for target function.
		router *mux.Router
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := analytics.NewHandler(tt.service)
			h.RegisterRoutes(tt.router)
		})
	}
}

func TestHandler_HandleTimePatterns(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		service analytics.Service
		// Named input parameters for target function.
		w http.ResponseWriter
		r *http.Request
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := analytics.NewHandler(tt.service)
			h.HandleTimePatterns(tt.w, tt.r)
		})
	}
}
