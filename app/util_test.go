package app

import (
	"testing"
	"net/url"
)

func TestCopy(t *testing.T) {
	orig := url.Values {"foo": {"bar"}}
	copied := Copy(orig)

	if len(copied) != 1 || len(copied["foo"]) != 1 || copied["foo"][0] != "bar" {
		t.Error(orig, copied)
	}

	orig["baz"] = []string {"baj"}

	if copied["baz"] != nil {
		t.Error(orig, copied)
	}
}
