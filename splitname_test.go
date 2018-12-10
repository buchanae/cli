package cli

import (
  "fmt"
)

func ExampleSplitIdent() {
  fmt.Println(SplitIdent("HelloWorld"))
  fmt.Println(SplitIdent("Hello"))
  fmt.Println(SplitIdent("HTTPServer"))
  fmt.Println(SplitIdent("ConfigureTLS"))
  // Output:
  // [Hello World]
  // [Hello]
  // [HTTP Server]
  // [Configure TLS]
}
