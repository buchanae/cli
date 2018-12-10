package cli

import (
  "fmt"
  "os"
)

func ExampleEnv() {
  p := Env("cli")
  // Note: by default, EnvProvider converts keys to uppercase + underscore.
  os.Setenv("CLI_FOO_BAR", "baz")
  val, ok := p.Lookup([]string{"foo", "bar"})
  fmt.Println(val, ok)
  // Output:
  // baz true
}
