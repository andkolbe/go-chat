package main

import (
	"fmt"
	"net/http"
	"testing"
)

func TestNoSurf(t *testing.T) {
	// we have to pass in a handler to no surf so it will hand back a handler

	var myH myHandler // type myHandler is defined in setup_test.go

	h := NoSurf(&myH) // we have to pass a http.Handler into NoSurf for it to run

	// test that No Surf takes in a handler and returns a handler
	switch v := h.(type) { // store in v whatever type h is
	case http.Handler:
		// do nothing. this is what we expect
	default:
		t.Error(fmt.Sprintf("type is not http.Handler, but is %T", v))
	}
}

func TestSessionLoad(t *testing.T) {
	var myH myHandler

	s := SessionLoad(&myH)

	switch v := s.(type) {
	case http.Handler:

	default:
		t.Error(fmt.Sprintf("type is not http.Handler, but is %T", v))
	}
}