package cli

import (
	"fmt"
)

func ExampleKeyFunc() {
	path := []string{"Foo", "bar", "BAZ"}
	fmt.Println(DotKey(path))
	fmt.Println(UnderscoreKey(path))
	fmt.Println(DashKey(path))
	// Output:
	// foo.bar.baz
	// foo_bar_baz
	// foo-bar-baz
}
