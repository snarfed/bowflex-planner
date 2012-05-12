// The complete bowflex-planner App Engine app.
//
// TODO: fan out doesnt_matters to both possibilities. e.g. x => dm => y should
// cost something but x => dm => x shouldn't.
package app

import (
	"fmt"
	"net/http"
	"strconv"
)

// possible values for enum params
var WEIGHT_MIN = 0
var WEIGHT_MAX = 95
var WEIGHT_STEP = 5
var ARMS_MIN = 0
var ARMS_MAX = 9
var HANDLES_VALUES = map[string]bool{
	"arms": true, "inner ground": true, "lat bar": true, "outer ground": true}
var HANDLE_LENGTH_VALUES = map[string]bool{
	"doesnt_matter": true, "long": true, "short": true}
var BACK_VALUES = map[string]bool{
	"doesnt_matter": true, "curved": true, "flat": true}
var SEAT_VALUES = map[string]bool{
	"doesnt_matter": true, "no": true, "yes": true}

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
	msg := fmt.Sprintf("Bad value %v for parameter %v", value, param)
	http.Error(w, msg, http.StatusBadRequest)
}

func generate(w http.ResponseWriter, r *http.Request) {
	routine, ok := parse_params(w, r)
	// fmt.Printf("@ %v %v\n", ok, len(routine))
	// fmt.Printf("@ %v\n", w.Header())
	if !ok {
		return
	}
	routines, min, avg, max := min_avg_max(routine)

	fmt.Fprintf(w, "<p>Minimum cost: %v<br />\n", min)
	fmt.Fprintf(w, "Average cost: %v<br />\n", avg)
	fmt.Fprintf(w, "Maximum cost: %v<br /></p>\n", max)

	fmt.Fprintf(w, "<p>%v minimum cost path(s):</p>\n", len(routines))
	fmt.Fprintln(w, "<table><tr><th>Name</th><th>Weight</th><th>Arms</th>"+
		"<th>Handles</th><th>Handle length</th><th>Back</th><th>Seat</th></tr>")
	for i, r := range routines {
		for _, e := range r {
			fmt.Fprintf(w, "<tr><td>%v</td><td>%v</td><td>%v</td><td>%v</td>"+
				"<td>%v</td><td>%v</td><td>%v</td></tr>\n",
				e.name, e.weight, e.arms, e.handles, e.handle_length, e.back, e.seat)
		}
		if i < len(r)-2 {
			fmt.Fprintln(w, "<tr><td><b>OR</b></td></tr>")
		}
	}
	fmt.Fprintln(w, "</table>")
}

func parse_params(w http.ResponseWriter, r *http.Request) (Routine, bool) {
	routine := Routine{}
	ok := true

	for i := 1; ; i++ {
		si := strconv.Itoa(i)

		// parse
		e := Exercise{
			name:          r.FormValue("name" + si),
			weight:        int_param("weight"+si, &ok, w, r),
			arms:          int_param("arms"+si, &ok, w, r),
			handles:       r.FormValue("handles" + si),
			handle_length: r.FormValue("handle_length" + si),
			back:          r.FormValue("back" + si),
			seat:          r.FormValue("seat" + si),
		}

		// validate
		if e.name == "" && e.weight == -1 && e.arms == -1 && e.handles == "" &&
			e.handle_length == "" && e.back == "" && e.seat == "" {
			break
		}

		ok = false
		if !HANDLES_VALUES[e.handles] {
			badParamError(w, "handles"+si, e.handles)
		} else if !HANDLE_LENGTH_VALUES[e.handle_length] {
			badParamError(w, "handle_length"+si, e.handle_length)
		} else if !BACK_VALUES[e.back] {
			badParamError(w, "back"+si, e.back)
		} else if !SEAT_VALUES[e.seat] {
			badParamError(w, "seat"+si, e.seat)
		} else if e.arms < ARMS_MIN || e.arms > ARMS_MAX {
			badParamError(w, "arms"+si, e.arms)
		} else if e.weight < WEIGHT_MIN || e.weight > WEIGHT_MAX ||
			e.weight%WEIGHT_STEP != 0 {
			badParamError(w, "weight"+si, e.weight)
		} else if e.name == "" {
			badParamError(w, "name"+si, "''")
		} else {
			ok = true
		}

		routine = append(routine, &e)
	}

	return routine, ok
}

// Returns the int value of a query parameter. If the parameter is not provided,
// returns -1. If it can't be converted to an integer, sets ok to false.
func int_param(param string, ok *bool, w http.ResponseWriter, r *http.Request) int {
	val_str := r.FormValue(param)
	if val_str == "" {
		return -1
	}

	val_int, err := strconv.ParseInt(val_str, 0, 0)
	if err != nil {
		badParamError(w, param, val_str)
		*ok = false
	}
	return int(val_int)
}

// Returns the min cost routine(s) and the mean and max costs.
func min_avg_max(routine Routine) ([]Routine, int, int, int) {
	steps := all_steps(routine)
	paths := make([]Path, Factorial(len(routine)))

	// generate all possible paths and calculate their costs.
	for i, perm := range Permutations(len(routine)) {
		paths[i].order = perm
		paths[i].cost = 0
		for j := 0; j < len(routine)-1; j++ {
			paths[i].cost += steps[[2]int{perm[j], perm[j+1]}]
		}
	}

	// find the min, max, and mean average cost paths
	var min_paths []*Path = nil
	min, max, sum := 0, 0, 0
	for i, path := range paths {
		sum += path.cost
		if min == 0 || path.cost < min {
			min = path.cost
			// not &path because path is a separate copy and changes
			min_paths = []*Path{&paths[i]}
		} else if path.cost == min {
			min_paths = append(min_paths, &paths[i])
		}
		if path.cost > max {
			max = path.cost
		}
	}
	avg := sum / len(paths)

	// populate min routines
	min_routines := make([]Routine, len(min_paths))
	for i, path := range min_paths {
		min_routines[i] = make(Routine, len(routine))
		for j, exercise_j := range path.order {
			min_routines[i][j] = routine[exercise_j]
		}
	}
	return min_routines, min, avg, max
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

// requires same order
func RoutinesEqual(this []Routine, that []Routine) bool {
	if len(this) != len(that) {
		return false
	}
	for i, this_i := range this {
		if !this_i.Equal(that[i]) {
			return false
		}
	}
	return true
}

// ugh, i want generics
func (this Routine) Equal(that Routine) bool {
	if len(this) != len(that) {
		return false
	}
	for i, this_i := range this {
		if this_i != that[i] {
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
