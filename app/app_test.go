package app

import (
	// "fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	// "strings"
	"testing"
)


// TODO: make this a const. go seems to not allow that though?
var BASE_PARAMS = url.Values {
	"name": {"foo"},
	"weight": {"40"},
	"arms": {"5"},
	"handles": {"arms"},
	"handle_length": {"long"},
	"back": {"flat"},
	"seat": {"yes"},
}


// Calls generate and checks the result HTTP code and body.
func generateAndExpect(params map[string]string, expected_code int, t *testing.T) {
	values := Copy(BASE_PARAMS)
	for key, val := range params {
		values.Set(key, val)
	}

	recorder := httptest.NewRecorder()
	url := "?" + values.Encode()
	req, _ := http.NewRequest("GET", url, nil)

	generate(recorder, req)
	if recorder.Code != expected_code {
		t.Error(url, recorder.Code, recorder.Body)
	}
}

func TestValidateBadEnumParams(t *testing.T) {
	for key := range BASE_PARAMS {
		if key != "name" {
			generateAndExpect(map[string]string {key: "bad"}, 400, t)
		}
	}
}

func TestValidateGoodEnumParams(t *testing.T) {
	generateAndExpect(nil, 0, t)
}


func TestCost(t *testing.T) {
	// cost := cost(
	// if 
}
