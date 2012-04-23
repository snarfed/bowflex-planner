package app

import (
	// "fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)


// Calls generate and checks the result HTTP code and body.
func generateAndExpect(url string, expected_code int, t *testing.T) {
	recorder := httptest.NewRecorder()
	// fmt.Println(recorder.Code)
	req, _ := http.NewRequest("GET", url, nil)
	generate(recorder, req)
	if recorder.Code != expected_code {
		t.Error(url, recorder.Code, recorder.Body)
	}
}


func TestValidateParams(t *testing.T) {
	bad_urls := []string {
		"/generate?handles=&handle_length=long&back=flat&seat=no",
		"/generate?handles=b&handle_length=long&back=flat&seat=no",
		"/generate?handles=x&handle_length=long&back=flat&seat=no",
		"/generate?handles=arms&handle_length=x&back=flat&seat=no",
		"/generate?handles=arms&handle_length=long&back=x&seat=no",
		"/generate?handles=arms&handle_length=long&back=flat&seat=x",
	}

	for _, url := range bad_urls {
		generateAndExpect(url, 400, t)
	}

	good_url := "/generate?handles=arms&handle_length=long&back=flat&seat=no"
	generateAndExpect(good_url, 0, t)
}


func TestCost(t *testing.T) {
}
