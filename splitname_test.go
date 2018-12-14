package cli

import (
	"reflect"
	"testing"
)

func TestSplitIdent(t *testing.T) {
	check := func(res []string, expect ...string) {
		if !reflect.DeepEqual(res, expect) {
			t.Errorf("expected %v got %v", expect, res)
		}
	}

	check(splitIdent("HelloWorld"), "Hello", "World")
	check(splitIdent("Hello"), "Hello")
	check(splitIdent("HTTPServer"), "HTTP", "Server")
	check(splitIdent("ConfigureTLS"), "Configure", "TLS")
}
