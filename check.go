package cli

import (
  "os"
  "github.com/spf13/cobra"
)

func Check(err error) {
  if err != nil {
    panic(err)
  }
}

func Open(path string) *os.File {
  fh, err := os.Open(path)
  Check(err)
  return fh
}

func CheckArgs(args []string, specs []ArgSpec) {
}

func CoerceString(arg string) string {
  return arg
}

func CoerceInts(args []string) []int {
  return nil
}

func CoerceInt(args string) int {
  return 0
}

func Cobra(specs ...CmdSpec) *cobra.Command {
  return nil
}
