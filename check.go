package cli

import (
	"fmt"
	"os"
)

type ErrFatal struct {
  error
}

type ErrUsage struct {
  error
}

// Fatal panics with an instance of ErrFatal with a formatted message.
func Fatal(msg string, args ...interface{}) {
	panic(ErrFatal{fmt.Errorf(msg, args...)})
}

// Usage panics with an instance of ErrFatal with a formatted message.
func Usage(msg string, args ...interface{}) {
	panic(ErrUsage{fmt.Errorf(msg, args...)})
}

// Check panics with an instance of ErrFatal if err != nil.
func Check(err error) {
	if err != nil {
		panic(ErrFatal{err})
	}
}

// Open opens a file with os.Open(path) but will panic if os.Open returns an error.
func Open(path string) *os.File {
	fh, err := os.Open(path)
	Check(err)
	return fh
}

// Run runs the CmdSpec with the given args.
// All panics are recovered and returned as an error.
// TODO only recover from instances of ErrFatal?
func Run(spec CmdSpec, args []string) (err error) {
	defer func() {
		if r := recover(); r != nil {
      switch z := r.(type) {
      case ErrUsage:
        err = z
        return
      case ErrFatal:
        err = z
        return
      default:
        panic(r)
      }
		}
	}()
	spec.Run(args)
	return
}

// CheckArgs panics if args cannot correctly fulfill the given ArgSpecs.
func CheckArgs(args []string, specs []ArgSpec) {
	if len(specs) == 0 && len(args) > 0 {
    Usage("unexpected args %v", args)
	}
  if len(specs) == 0 {
    return
  }

	variadic := specs[len(specs)-1].Variadic
	if variadic {
		min := len(specs) - 1
		if len(args) < min {
      Usage("expected at least %d arg", min)
		}

    for i := 0; i < min; i++ {
      spec := specs[i]
      err := coerceSet(spec.Value, args[i])
      if err != nil {
        Usage("coercing %q to %T", args[i], spec.Value)
      }
    }

    if len(args) > min {
      spec := specs[len(specs) - 1]
      err := coerceSet(spec.Value, args[min:])
      if err != nil {
        // TODO these error messages suck. "coercing ["x"] to *[]int"
        Usage("coercing %q to %T", args[min:], spec.Value)
      }
    }

	} else {
		if len(args) != len(specs) {
      Usage("expected exactly %d args", len(specs))
		}
    for i := 0; i < len(args); i++ {
      spec := specs[i]
      err := coerceSet(spec.Value, args[i])
      if err != nil {
        Usage("coercing %q to %T", args[i], spec.Value)
      }
    }
	}
}
