package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// table test

// holds whatever we are posting to a page
type postData struct {
	key string
	value string
}

var theTests = []struct { // slice of structs because we will have more than one test we want to run in our table test
	name string
	url string
	method string
	params []postData
	expectedStatusCode int
}{
	{"login", "/", "GET", []postData{}, http.StatusOK},
	{"chat", "/chat", "GET", []postData{}, http.StatusOK},
	{"register", "/register", "GET", []postData{}, http.StatusOK},

}

func testHandlers(t *testing.T) {
	routes := getRoutes() // getRoutes comes from setup_test

	// we need to create a server and a client that can call that server 
	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	for _, e := range theTests {
		// determine whether we are going to run a GET or POST request
		if e.method == "GET" {
			// we want to test as if we are a client accessing a web page
			resp, err := ts.Client().Get(ts.URL + e.url)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}
			// look at the status code we get back from the server and compare that against the expected status code
			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf("for %s, expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
			}
		} else { // POST

		}
	}
}