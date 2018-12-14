package cli

import (
	"fmt"
	"os"
)

func ExampleEnv() {
	// Note: Env converts keys to uppercase + underscore.
	os.Setenv("CLI_FOO_BAR", "baz")

	// TODO gross
	val := ""
	key := []string{"foo", "bar"}
	opts := []*Opt{
		{Key: key, Value: &val},
	}

	l := NewLoader(opts, Env("cli"))
	l.Load()
	fmt.Println(val)
	// Output:
	// baz
}
