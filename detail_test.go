package cli

import (
  "fmt"
)

func ExampleParseOptDetail() {
  det := ParseOptDetail(OptSpec{
    Doc: `Opt detail synopsis.
Deprecated: this option is deprecated.
`,
  })
  fmt.Println(det.Synopsis)
  fmt.Println(det.Deprecated)
  // Output:
  // Opt detail synopsis.
  // this option is deprecated.
}

func ExampleParseCmdDetail() {
  det := ParseCmdDetail(&cmdSpec{
    name: "WordCount",
    doc: `Count the number of words in a file.

word-count counts the number of words in a file and writes a single number to stdout.
This is a second line of docs.

Name: word-count
Deprecated: please don't use word-count, it will be removed.
Hidden
Aliases: count-words count-word wc
Example: word-count <file.txt>
`,
  })
  fmt.Println(det.Synopsis)
  fmt.Println(det.Doc)
  fmt.Println(det.Name)
  fmt.Println(det.Example)
  fmt.Println(det.Hidden)
  fmt.Println(det.Aliases)
  // Output:
  // Count the number of words in a file.
  // word-count counts the number of words in a file and writes a single number to stdout.
  // This is a second line of docs.
  // word-count
  // word-count <file.txt>
  // true
  // [count-words count-word wc]
}

type cmdSpec struct {
  CmdSpec
  name string
  doc string
}

func (c *cmdSpec) Name() string {
  return c.name
}

func (c *cmdSpec) Doc() string {
  return c.doc
}
