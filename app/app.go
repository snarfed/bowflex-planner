// The complete bowflex-planner App Engine app.
package app

import (
	"fmt"
	"net/http"
	"strconv"
)

// possible values for inputs in forms.
// TODO: make this a const. go seems to not allow that though?
var ENUM_PARAM_VALUES = map[string] map[string] bool {
	"handles": map[string] bool { "arms": true, "inner ground": true, "lat bar": true,
		"outer ground": true, },
    "handle_length" : { "doesnt_matter": true, "long": true, "short": true, },
    "back": { "doesnt_matter": true, "curved": true, "flat": true, },
    "seat": { "doesnt_matter": true, "no": true, "yes": true, },
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

func badParamError(w http.ResponseWriter, param string, value interface{}) {
	msg := fmt.Sprintf("Bad value %s for parameter %s", value, param)
	http.Error(w, msg, http.StatusBadRequest)
}

func generate(w http.ResponseWriter, r *http.Request) {
	// parse and validate input query parameters
	error := false

	for param, expected := range ENUM_PARAM_VALUES {
		actual := r.FormValue(param)
		if !expected[actual] {
			badParamError(w, param, actual)
			error = true
		}
	}

	for _, param := range []string {"weight", "arms"} {
		val, err := strconv.ParseUint(r.FormValue(param), 0, 0)
		if err != nil {
			badParamError(w, param, val)
			error = true
		}
	}

	if r.FormValue("name") == "" {
		badParamError(w, "name", "''")
	}

	if error {
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
	cost := 0

	// if from.arms != to.arms {
	// 	cost += 1
	// }
	// if from.weight != to.weight {
	// 	cost += 3
	// }
	// if from.handles == "lat_bar" || to.handles == "lat_bar" {
	// 	cost += 2
	// } else {
	// 	cost += 1
	// }
	// if (from.handle_length == "short" && to.handle_length == "long") ||
	// 	 (from.handle_length == "long" && to.handle_length == "short") {
	// 	cost += 1
	// }
	// if (from.back == "curved" && to.back == "flat") ||
	// 	 (from.back == "flat" && to.back == "curved") {
	// 	cost += 1
	// }
	// if (from.seat == "yes" && to.seat == "no") ||
	// 	 (from.seat == "no" && to.seat == "yes") {
	// 	cost += 1
	// }

	return cost
}

