package main

import (
	"fmt"
	"testing"

	"github.com/go-chi/chi/v5"
)


func TestRoutes(t *testing.T) {

	mux := routes()

	switch mux.(type) {
	case *chi.Mux:
		// do nothing test passed
	default: 
		t.Error(fmt.Sprintln("type is not *chi.Mux"))
	}
}