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
var BASE_PARAMS = url.Values{
	"name":          {"foo"},
	"weight":        {"40"},
	"arms":          {"5"},
	"handles":       {"arms"},
	"handle_length": {"long"},
	"back":          {"flat"},
	"seat":          {"yes"},
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
			generateAndExpect(map[string]string{key: "bad"}, 400, t)
		}
	}
}

func TestValidateGoodEnumParams(t *testing.T) {
	generateAndExpect(nil, 0, t)
}

func TestCost(t *testing.T) {
	from := Exercise{"", 10, 0, "arms", "short", "flat", "yes"}

	// weight (3) + arms (1) + handles arms => outer ground (1) == 5
	to := Exercise{"", 20, 1, "outer ground", "short", "flat", "yes"}
	if cost(&from, &to) != 5 {
		t.Error(from, to)
	}

	// handle length short => long (1) + back (1) + seat (1) == 3
	to = Exercise{"", 10, 0, "arms", "long", "curved", "no"}
	if cost(&from, &to) != 3 {
		t.Error(from, to)
	}

	// handle length, back, seat all doesn't matter == 0
	to = Exercise{"", 10, 0, "arms", "doesnt_matter", "doesnt_matter", "doesnt_matter"}
	if cost(&from, &to) != 0 {
		t.Error(from, to)
	}
}

func TestCostSum(t *testing.T) {
	exercises := []*Exercise{}
	sum := cost_sum(exercises)
	if sum != 0 {
		t.Error(sum)
	}

	exercises = append(exercises,
		&Exercise{"", 10, 0, "arms", "short", "flat", "yes"})
	sum = cost_sum(exercises)
	if sum != 0 {
		t.Error(sum)
	}

	exercises = append(exercises,
		&Exercise{"", 10, 0, "arms", "doesnt_matter", "doesnt_matter", "doesnt_matter"})
	sum = cost_sum(exercises)
	if sum != 0 {
		t.Error(sum)
	}

	// weight (3) + arms (1) + handles arms => outer ground (1) == 5
	// handle length short => long (1) + back (1) + seat (1) == 3
	exercises = append(exercises,
		&Exercise{"", 20, 1, "outer ground", "short", "flat", "yes"},
		&Exercise{"", 20, 1, "outer ground", "long", "curved", "no"})
	sum = cost_sum(exercises)
	if sum != 8 {
		t.Error(sum)
	}
}

func (this Steps) Equal(that Steps) bool {
	if len(this) != len(that) {
		return false
	}

	for key, val := range this {
		if val != that[key] {
			return false
		}
	}

	return true
}

func TestAllSteps(t *testing.T) {
	exercises := []*Exercise{}
	steps := all_steps(exercises)
	if len(steps) != 0 {
		t.Error(steps)
	}

	exercises = append(exercises,
		&Exercise{"", 10, 0, "arms", "short", "flat", "yes"})
	steps = all_steps(exercises)
	if len(steps) != 0 {
		t.Error(steps)
	}

	exercises = append(exercises,
		&Exercise{"", 10, 0, "arms", "doesnt_matter", "doesnt_matter", "doesnt_matter"})
	steps = all_steps(exercises)
	if !steps.Equal(Steps{[2]int{0, 1}: 0}) {
		t.Error(steps)
	}

	// weight (3) + arms (1) + handles arms => outer ground (1) == 5
	// handle length short => long (1) + back (1) + seat (1) == 3
	exercises = append(exercises,
		&Exercise{"", 20, 1, "outer ground", "short", "flat", "yes"},
		&Exercise{"", 20, 1, "outer ground", "long", "curved", "no"})
	steps = all_steps(exercises)
	if !steps.Equal(Steps{
		[2]int{0, 1}: 0,
		[2]int{0, 2}: 5,
		[2]int{0, 3}: 8,
		[2]int{1, 2}: 5,
		[2]int{1, 3}: 5,
		[2]int{2, 3}: 3,
	}) {
		t.Error(steps)
	}
}

// func TestMinPath(t *testing.T) {
// 	exercises := []*Exercise{}
// 	path, cost := min_path(exercises)
// 	if path != exercises {
// 		t.Error(sum)
// 	} else if cost != 0 {
// 		t.Error(cost)
// 	}

// 	exercises = append(exercises,
// 		&Exercise{"", 10, 0, "arms", "short", "flat", "yes"})
// 	path = min_path(exercises)
// 	if path != exercises {
// 		t.Error(path)
// 	} else if cost != 0 {
// 		t.Error(cost)
// 	}

// 	exercises = append(exercises,
// 		&Exercise{"", 10, 0, "arms", "doesnt_matter", "doesnt_matter", "doesnt_matter"})
// 	path = min_path(exercises)
// 	if path != exercises {
// 		t.Error(path)
// 	} else if cost != 0 {
// 		t.Error(cost)
// 	}

// 	// weight (3) + arms (1) + handles arms => outer ground (1) == 5
// 	// handle length short => long (1) + back (1) + seat (1) == 3
// 	exercises = append(exercises,
// 		&Exercise{"", 20, 1, "outer ground", "short", "flat", "yes"},
// 		&Exercise{"", 20, 1, "outer ground", "long", "curved", "no"})
// 	path = min_path(exercises)
// 	if path != exercises {
// 		t.Error(path)
// 	} else if cost != 8 {
// 		t.Error(cost)
// 	}
// }
