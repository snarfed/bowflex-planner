// Misc utility functions.
package app

import (
	"net/url"
)

func Copy(input url.Values) url.Values {
	copied := make(url.Values, len(input))
	for key, val := range input {
		copied[key] = val
	}
	return copied
}
