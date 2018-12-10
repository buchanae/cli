package cli

import (
	"fmt"
	"os"
)

type ErrFatal struct {
	err error
}

func (e ErrFatal) Error() string {
	return e.err.Error()
}

// Fatal panics with an instance of ErrFatal with a formatted message.
func Fatal(msg string, args ...interface{}) {
	panic(ErrFatal{fmt.Errorf(msg, args...)})
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

// CheckArgs panics if args cannot correctly fulfill the given ArgSpecs.
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
