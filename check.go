package cli

import (
  "bytes"
  "bufio"
  "fmt"
  "go/doc"
  "os"
  "strings"
  "github.com/spf13/cobra"
)

type ErrFatal struct {
  err error
}
func (e ErrFatal) Error() string {
  return e.err.Error()
}

func Fatal(msg string, args ...interface{}) {
  panic(ErrFatal{fmt.Errorf(msg, args...)})
}

func Check(err error) {
  if err != nil {
    panic(ErrFatal{err})
  }
}

func Run(spec CmdSpec, args []string) (err error) {
  defer func() {
    if r := recover(); r != nil {
      if x, ok := r.(error); ok {
        err = x
      } else {
        err = fmt.Errorf("%v", r)
      }
    }
  }()
  spec.Run(args)
  return
}

func Open(path string) *os.File {
  fh, err := os.Open(path)
  Check(err)
  return fh
}

func CheckArgs(args []string, specs []ArgSpec) {
  if len(specs) == 0 && len(args) > 0 {
    Fatal("unexpected args %v", args)
  }

  variadic := specs[len(specs)-1].Variadic
  if variadic {
    min := len(specs) - 1
    if len(args) < min {
      Fatal("expected at least %d arg", min)
    }
  } else {
    if len(args) != len(specs) {
      Fatal("expected exactly %d args", len(specs))
    }
  }
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

type CmdDetail struct {
  Name string
  Path []string
  Doc string
  Synopsis string
  Example string
  Deprecated string
  Hidden bool
  Aliases []string
}

func ParseDetail(spec CmdSpec) CmdDetail {

  det := CmdDetail{
    Doc: spec.Doc(),
    Synopsis: doc.Synopsis(spec.Doc()),
  }

  setNamePath := func(parts []string) {
    if len(parts) == 0 {
      return
    }

    for _, p := range parts {
      det.Path = append(det.Path, strings.ToLower(p))
    }
    det.Name = det.Path[len(det.Path)-1]
  }

  parts := SplitIdent(spec.Name())
  setNamePath(parts)

  scan := bufio.NewScanner(bytes.NewBufferString(spec.Doc()))
  for scan.Scan() {
    line := strings.TrimSpace(scan.Text())
    if strings.HasPrefix(line, "Name: ") {
      name := strings.TrimPrefix(line, "Name: ")
      parts := strings.Split(name, " ")
      setNamePath(parts)
    }
    if strings.HasPrefix(line, "Deprecated: ") {
      det.Deprecated = strings.TrimPrefix(line, "Deprecated: ")
    }
    if strings.HasPrefix(line, "Example: ") {
      det.Example = strings.TrimPrefix(line, "Example: ")
    }
    if line == "Hidden" {
      det.Hidden = true
    }
    if strings.HasPrefix(line, "Aliases: ") {
      det.Aliases = strings.Split(strings.TrimPrefix(line, "Aliases: "), " ")
    }
  }

  return det
}

func AddCobra(cmd *cobra.Command, specs ...CmdSpec) {
  addCobra(cmd, 0, specs...)
}

func addCobra(cmd *cobra.Command, depth int, specs ...CmdSpec) {
  sub := map[string]*cobra.Command{}

  for _, spec := range specs {
    det := ParseDetail(spec)

    if depth == len(det.Path) - 1 {
      addCobraDet(cmd, det, spec)
      continue
    }

    name := det.Path[depth]
    parent, ok := sub[name]
    if !ok {
      parent = &cobra.Command{
        Use: name,
      }
      cmd.AddCommand(parent)
      sub[name] = parent
    }
    addCobra(parent, depth + 1, spec)
  }
}

func addCobraDet(cmd *cobra.Command, det CmdDetail, spec CmdSpec) {
  cmd.AddCommand(&cobra.Command{
    Use: det.Name,
    Short: det.Synopsis,
    Long: det.Doc,
    Example: det.Example,
    Deprecated: det.Deprecated,
    Hidden: det.Hidden,
    Aliases: det.Aliases,
    RunE: func(cmd *cobra.Command, args []string) error {
      return Run(spec, args)
    },
  })
}
