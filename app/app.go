package app

import (
	"fmt"
	"net/http"
	"sort"
	"strings"
)

// possible values for inputs in forms
var PARAM_VALUES = map[string] []string {
	// note that each list must be sorted so we can use sort.Search :/
	"handles": { "arms", "inner ground", "lat bar", "outer ground", },
    "handle_length" : { "doesnt_matter", "long", "short", },
    "back": { "curved", "doesnt_matter", "flat", },
    "seat": { "doesnt_matter", "no", "yes", },
}

type Exercise struct {
	Name string
	Weight int
	Arms int
	Handles string
	Handle_length string
	Back string
	Seat bool
}

func init() {
	http.HandleFunc("/generate", generate)
}

func generate(w http.ResponseWriter, r *http.Request) {
	// validate input query parameters
	bad_param := false
	for param, expected := range PARAM_VALUES {
		actual := r.FormValue(param)
		if expected[sort.SearchStrings(expected, actual)] != actual {
			msg := fmt.Sprintf(
				"Bad value '%s' for parameter %s; expected one of %s",
				actual, param, strings.Join(expected, ", "))
			http.Error(w, msg, http.StatusBadRequest)
			bad_param = true
		}
	}

	if bad_param {
		return
	}

	
}

// cost heuristics:
// arms: 0 if same, 1 if different
// weight: 0 if same, 3 if different (TODO)
// handles: lat bar <=> * 2, otherwise 1
// handle_length: long <=> short 1
// back: curved <=> flat 1
// seat: no <=> yes 1
func cost(from *Exercise, to *Exercise) int {
	
}

