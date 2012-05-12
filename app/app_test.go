package app

import (
	// "fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

// TODO: make this a const. go seems to not allow that though?
var BASE_PARAMS = url.Values{
	"name1":          {"foo"},
	"weight1":        {"40"},
	"arms1":          {"5"},
	"handles1":       {"arms"},
	"handle_length1": {"long"},
	"back1":          {"flat"},
	"seat1":          {"yes"},
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
		val := "bad"
		if strings.HasPrefix(key, "name") {
			val = ""
		}
		generateAndExpect(map[string]string{key: val}, 400, t)
	}

	generateAndExpect(map[string]string{"arms1": "-1"}, 400, t)
	generateAndExpect(map[string]string{"arms1": "10"}, 400, t)
	generateAndExpect(map[string]string{"arms1": "3.14"}, 400, t)

	generateAndExpect(map[string]string{"weight1": "-5"}, 400, t)
	generateAndExpect(map[string]string{"weight1": "100"}, 400, t)
	generateAndExpect(map[string]string{"weight1": "7"}, 400, t)
	generateAndExpect(map[string]string{"weight1": "3.14"}, 400, t)
}

func TestValidateGoodEnumParams(t *testing.T) {
	generateAndExpect(nil, 200, t)
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
	routine := Routine{}
	sum := cost_sum(routine)
	if sum != 0 {
		t.Error(sum)
	}

	routine = append(routine,
		&Exercise{"", 10, 0, "arms", "short", "flat", "yes"})
	sum = cost_sum(routine)
	if sum != 0 {
		t.Error(sum)
	}

	routine = append(routine,
		&Exercise{"", 10, 0, "arms", "doesnt_matter", "doesnt_matter", "doesnt_matter"})
	sum = cost_sum(routine)
	if sum != 0 {
		t.Error(sum)
	}

	// weight (3) + arms (1) + handles arms => outer ground (1) == 5
	// handle length short => long (1) + back (1) + seat (1) == 3
	routine = append(routine,
		&Exercise{"", 20, 1, "outer ground", "short", "flat", "yes"},
		&Exercise{"", 20, 1, "outer ground", "long", "curved", "no"})
	sum = cost_sum(routine)
	if sum != 8 {
		t.Error(sum)
	}
}

func TestAllSteps(t *testing.T) {
	routine := Routine{}
	steps := all_steps(routine)
	if len(steps) != 0 {
		t.Error(steps)
	}

	routine = append(routine,
		&Exercise{"", 10, 0, "arms", "short", "flat", "yes"})
	steps = all_steps(routine)
	if len(steps) != 0 {
		t.Error(steps)
	}

	routine = append(routine,
		&Exercise{"", 10, 0, "arms", "doesnt_matter", "doesnt_matter", "doesnt_matter"})
	steps = all_steps(routine)
	if !steps.Equal(Steps{[2]int{0, 1}: 0,
		[2]int{1, 0}: 0}) {
		t.Error(steps)
	}

	// weight (3) + arms (1) + handles arms => outer ground (1) == 5
	// handle length short => long (1) + back (1) + seat (1) == 3
	routine = append(routine,
		&Exercise{"", 20, 1, "outer ground", "short", "flat", "yes"},
		&Exercise{"", 20, 1, "outer ground", "long", "curved", "no"})
	steps = all_steps(routine)
	if !steps.Equal(Steps{
		[2]int{0, 1}: 0, [2]int{1, 0}: 0,
		[2]int{0, 2}: 5, [2]int{2, 0}: 5,
		[2]int{0, 3}: 8, [2]int{3, 0}: 8,
		[2]int{1, 2}: 5, [2]int{2, 1}: 5,
		[2]int{1, 3}: 5, [2]int{3, 1}: 5,
		[2]int{2, 3}: 3, [2]int{3, 2}: 3,
	}) {
		t.Error(steps)
	}
}

func TestMinAvgMax(t *testing.T) {
	r := Routine{}
	ExpectMinAvgMax(t, r, []Routine{{}}, 0, 0, 0)

	r = append(r,
		&Exercise{"", 10, 0, "arms", "short", "flat", "yes"})
	ExpectMinAvgMax(t, r, []Routine{r}, 0, 0, 0)

	r = append(r,
		&Exercise{"", 10, 0, "arms", "doesnt_matter", "doesnt_matter", "doesnt_matter"})
	expected := []Routine{r.Reversed()}
	ExpectMinAvgMax(t, r, expected, 0, 0, 0)

	// weight (3) + arms (1) + handles arms => outer ground (1) == 5
	// handle length short => long (1) + back (1) + seat (1) == 3
	r = append(r,
		&Exercise{"", 20, 1, "outer ground", "short", "flat", "yes"},
		&Exercise{"", 20, 1, "outer ground", "long", "curved", "no"})
	expected = []Routine{
		{r[0], r[1], r[2], r[3]},
		{r[1], r[0], r[2], r[3]},
		{r[0], r[1], r[3], r[2]},
		{r[3], r[2], r[0], r[1]},
		{r[2], r[3], r[1], r[0]},
		{r[3], r[2], r[1], r[0]},
	}
	ExpectMinAvgMax(t, r, expected, 8, 13, 18)

	// weight (3) + handles => lat_bar (2) == 6
	// weight (3) + arms (1) + handles <= lat_bar (2) + seat (1) == 7
	r = append(r,
		&Exercise{"", 60, 1, "lat_bar", "long", "curved", "no"},
		&Exercise{"", 10, 2, "outer ground", "long", "doesnt_matter", "yes"})
	expected = []Routine{
		{r[4], r[3], r[2], r[0], r[1], r[5]},
		{r[0], r[1], r[5], r[2], r[3], r[4]},
		{r[5], r[1], r[0], r[2], r[3], r[4]},
		{r[4], r[3], r[2], r[5], r[1], r[0]},
	}
	ExpectMinAvgMax(t, r, expected, 15, 25, 35)
}

func ExpectMinAvgMax(t *testing.T, input Routine, expected []Routine,
	expected_min int, expected_avg int, expected_max int) {
	paths, min, avg, max := min_avg_max(input)
	if !RoutinesEqual(paths, expected) {
		t.Error(expected, paths)
	} else if min != expected_min || avg != expected_avg || max != expected_max {
		t.Error(min, avg, max)
	}
}
