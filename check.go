package cli

import (
	"fmt"
	"strings"
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

// Check panics with an instance of ErrFatal if err != nil.
func Check(err error) {
	if err != nil {
		panic(ErrFatal{err})
	}
}

// Run runs the Cmd with the given args.
// Panics of type ErrFatal and ErrUsage are recovered and returned as an error,
// all other panics are passed through.
func Run(spec Spec, l *Loader, raw []string) (err error) {
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

	cmd := spec.Cmd()

	err = validateArgs(cmd.Args, raw)
	if err != nil {
		return err
	}

	err = loadArgs(cmd.Args, l, raw)
	if err != nil {
		return
	}

	// load option values.
	l.Load()
	errs := l.Errors()
	if errs != nil {
		return combineErrors(errs)
	}

	spec.Run()
	return
}

func combineErrors(errs []error) error {
	var lines []string
	for _, err := range errs {
		lines = append(lines, err.Error())
	}
	return fmt.Errorf(strings.Join(lines, "\n"))
}

func loadArgs(args []*Arg, l *Loader, raw []string) error {

	for i := 0; i < len(args); i++ {
		arg := args[i]

		var val interface{}
		if arg.Variadic {
			val = raw[i:]
		} else {
			val = raw[i]
		}

		err := l.Coerce(arg.Value, val)
		if err != nil {
			return ErrUsage{err}
		}
	}
	return nil
}

// validateArgs checks that the number of args given on the CLI
// matches the number of args needed by the CLI function.
func validateArgs(specs []*Arg, args []string) error {
	if len(specs) == 0 && len(args) > 0 {
		return ErrUsage{fmt.Errorf("unexpected args %v", args)}
	}
	if len(specs) == 0 {
		return nil
	}

	// variadic functions e.g. func HelloWorld(names ...string)
	variadic := specs[len(specs)-1].Variadic
	if variadic {
		min := len(specs) - 1
		if len(args) < min {
			return ErrUsage{fmt.Errorf("expected at least %d arg", min)}
		}
		return nil
	}

	if len(args) != len(specs) {
		return ErrUsage{fmt.Errorf("expected exactly %d args", len(specs))}
	}
	return nil
}
