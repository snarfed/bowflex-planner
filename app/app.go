// The complete bowflex-planner App Engine app.
//
// TODO: fan out doesnt_matters to both possibilities. e.g. x => dm => y should
// cost something but x => dm => x shouldn't.
// STATE: todos, more min cost tests
package app

import (
	"fmt"
	"net/http"
	"strconv"
)

// possible values for inputs in forms.
// TODO: make this a const. go seems to not allow that though?
var ENUM_PARAM_VALUES = map[string]map[string]bool{
	"handles":       {"arms": true, "inner ground": true, "lat bar": true, "outer ground": true},
	"handle_length": {"doesnt_matter": true, "long": true, "short": true},
	"back":          {"doesnt_matter": true, "curved": true, "flat": true},
	"seat":          {"doesnt_matter": true, "no": true, "yes": true},
}

type Exercise struct {
	name          string
	weight        int
	arms          int
	handles       string
	handle_length string
	back          string
	seat          string
}

func (e *Exercise) String() string {
	return fmt.Sprint(*e)
}

type Routine []*Exercise

type Steps map[[2]int]int

type Path struct {
	order []int
	cost  int
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

	for _, param := range []string{"weight", "arms"} {
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

func min_path(routine Routine) (Routine, int) {
	steps := all_steps(routine)
	n := len(routine)
	paths := make([]Path, Factorial(n))

	// generate all possible paths and calculate their costs.
	for i, perm := range Permutations(n) {
		paths[i].order = perm
		paths[i].cost = 0
		for j := 0; j < len(routine)-1; j++ {
			paths[i].cost += steps[[2]int{perm[j], perm[j+1]}]
		}
	}

	// linear search for path with lowest cost. (could use heap.)
	var min *Path = nil
	for _, path := range paths {
		if min == nil || path.cost < min.cost {
			min = &path
		}
	}

	// populate return value
	min_routine := make(Routine, 0, n)
	for _, i := range min.order {
		min_routine = append(min_routine, routine[i])
	}
	return min_routine, min.cost
}

func all_steps(routine Routine) Steps {
	n := len(routine)
	steps := make(Steps, n*n)

	if n < 2 {
		return steps
	}

	for i, from := range routine[:len(routine)-1] {
		for j, to := range routine[i+1:] {
			ij_cost := cost(from, to)
			steps[[2]int{i, i + j + 1}] = ij_cost
			steps[[2]int{i + j + 1, i}] = ij_cost
		}
	}
	return steps
}

func cost_sum(routine Routine) int {
	sum := 0
	for i := 0; i < len(routine)-1; i++ {
		sum += cost(routine[i], routine[i+1])
	}
	return sum
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

	if from.arms != to.arms {
		cost += 1
	}
	if from.weight != to.weight {
		cost += 3
	}
	if from.handles == "lat_bar" || to.handles == "lat_bar" {
		cost += 2
	} else if from.handles != to.handles {
		cost += 1
	}
	if (from.handle_length == "short" && to.handle_length == "long") ||
		(from.handle_length == "long" && to.handle_length == "short") {
		cost += 1
	}
	if (from.back == "curved" && to.back == "flat") ||
		(from.back == "flat" && to.back == "curved") {
		cost += 1
	}
	if (from.seat == "yes" && to.seat == "no") ||
		(from.seat == "no" && to.seat == "yes") {
		cost += 1
	}

	return cost
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

// i originally started to implement Routine as a ring to make this easier, but
// then thought it would be more complexity overall.
func (this Routine) Equal(that Routine) bool {
	if len(this) != len(that) {
		return false
	} else if len(this) == 0 {
		return true
	}

	return this.equal_ordered(that) || this.equal_ordered(that.Reversed())
}

func (this Routine) equal_ordered(that Routine) bool {
	// find the first exercise
	offset := -1
	for j, that_j := range that {
		if this[0] == that_j {
			offset = j
			break
		}
	}
	if offset == -1 {
		return false
	}

	// // determine direction
	// offset_prev := offset - 1
	// if offset_prev < 0 {
	// 	offset_prev += n
	// }

	// step := 0
	// if this[1] == that[(offset+1)%n] {
	// 	step = 1
	// } else if this[1] == that[offset_prev] {
	// 	step = -1
	// } else {
	// 	return false
	// }

	// start there and check the rest of the exercises
	for i, this_i := range this {
		// j := (i + offset) % n
		// if j < 0 {
		// 	j += n
		// }
		if this_i != that[(i+offset)%len(this)] {
			return false
		}
	}

	return true
}

func (this Routine) Reversed() Routine {
	reversed := make(Routine, len(this))
	for i, e := range this {
		reversed[len(this)-i-1] = e
	}
	return reversed
}
