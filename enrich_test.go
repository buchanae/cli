package cli

import (
	"fmt"
	"os"
)

func ExampleEnrich() {
	cmd := &Cmd{
		RawName: "WordCount",
		RawDoc: `Count the number of words in a file.

word-count counts the number of words in a file and writes a single number to stdout.
This is a second line of docs.

Deprecated: please don't use word-count, it will be removed.
Hidden
Aliases: count-words count-word wc
Example: word-count <file.txt>
`,
		Opts: []*Opt{
			{
				RawDoc: `Opt detail synopsis.`,
			},
			{
				DefaultValue: os.Stderr,
			},
		},
	}

	Enrich(cmd)

	fmt.Println("NAME: ", cmd.Name)
	fmt.Println("PATH: ", cmd.Path)
	fmt.Println("SYN:  ", cmd.Synopsis)
	fmt.Println("EX:   ", cmd.Example)
	fmt.Println("DEP:  ", cmd.Deprecated)
	fmt.Println("HIDE: ", cmd.Hidden)
	fmt.Println("ALIAS:", cmd.Aliases)
	fmt.Println("SYN:  ", cmd.Opts[0].Synopsis)
	fmt.Println("DEF:  ", cmd.Opts[1].DefaultString)
	// Output:
	// NAME:  count
	// PATH:  [word count]
	// SYN:   Count the number of words in a file.
	// EX:    word-count <file.txt>
	// DEP:   please don't use word-count, it will be removed.
	// HIDE:  true
	// ALIAS: [count-words count-word wc]
	// SYN:   Opt detail synopsis.
	// DEF:   os.Stderr
}
