package cerrors

import (
	"errors"
	"fmt"
	"io"
	"path"
	"runtime"
	"strconv"
)

type stackTrace interface {
	error
	StackTrace() StackTrace
}

type withStack struct {
	cause error
	stack StackTrace
}

func newWithStack(err error) error {
	if err == nil {
		return nil
	}

	var stErr stackTrace
	if errors.As(err, &stErr) {
		return err
	}

	return &withStack{cause: err, stack: callers(2)} //nolint:gomnd // self-explained
}

// *** Code from https://github.com/pkg/errors/blob/master/stack.go ** //

var _ stackTrace = (*withStack)(nil)

func (w *withStack) Error() string          { return w.cause.Error() }
func (w *withStack) Cause() error           { return w.cause }
func (w *withStack) Unwrap() error          { return w.cause }
func (w *withStack) StackTrace() StackTrace { return w.stack }

// Frame represents a program counter inside a stack frame.
// For historical reasons if Frame is interpreted as a uintptr
// its value represents the program counter + 1.
type Frame uintptr

// pc returns the program counter for this frame;
// multiple frames may have the same PC value.
func (f Frame) pc() uintptr { return uintptr(f) - 1 }

// file returns the full path to the file that contains the
// function for this Frame's pc.
func (f Frame) file() string {
	fn := runtime.FuncForPC(f.pc())
	if fn == nil {
		return "unknown"
	}
	file, _ := fn.FileLine(f.pc())
	return file
}

// line returns the line number of source code of the
// function for this Frame's pc.
func (f Frame) line() int {
	fn := runtime.FuncForPC(f.pc())
	if fn == nil {
		return 0
	}
	_, line := fn.FileLine(f.pc())
	return line
}

// name returns the name of this function, if known.
func (f Frame) name() string {
	fn := runtime.FuncForPC(f.pc())
	if fn == nil {
		return "unknown"
	}
	return fn.Name()
}

// Format formats the frame according to the fmt.Formatter interface.
func (f Frame) Format(s fmt.State, verb rune) {
	switch verb {
	case 's':
		switch {
		case s.Flag('+'):
			io.WriteString(s, f.name())
			io.WriteString(s, "\n\t")
			io.WriteString(s, f.file())
		default:
			io.WriteString(s, path.Base(f.file()))
		}
	case 'v':
		f.Format(s, 's')
		io.WriteString(s, ":")
		io.WriteString(s, strconv.Itoa(f.line()))
	}
}

// StackTrace is stack of Frames from innermost (newest) to outermost (oldest).
type StackTrace []Frame

// Format formats the stack of Frames according to the fmt.Formatter interface.
//
//	%s	lists source files for each Frame in the stack
//	%v	lists the source file and line number for each Frame in the stack
//
// Format accepts flags that alter the printing of some verbs, as follows:
//
//	%+v   Prints filename, function, and line number for each Frame in the stack.
func (st StackTrace) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		switch {
		case s.Flag('+'):
			for _, f := range st {
				io.WriteString(s, "\n")
				f.Format(s, verb)
			}
		case s.Flag('#'):
			fmt.Fprintf(s, "%#v", []Frame(st))
		default:
			st.formatSlice(s, verb)
		}
	case 's':
		st.formatSlice(s, verb)
	}
}

// formatSlice will format this StackTrace into the given buffer as a slice of
// Frame, only valid when called with '%s' or '%v'.
func (st StackTrace) formatSlice(s fmt.State, verb rune) {
	io.WriteString(s, "[")
	for i, f := range st {
		if i > 0 {
			io.WriteString(s, " ")
		}
		f.Format(s, verb)
	}
	io.WriteString(s, "]")
}

// callers mirrors the code in github.com/pkg/errors,
// but makes the depth customizable.
func callers(depth int) StackTrace {
	const numFrames = 32
	var pcs [numFrames]uintptr
	n := runtime.Callers(2+depth, pcs[:])
	f := make([]Frame, n)
	for i := 0; i < n; i++ {
		f[i] = Frame(pcs[i])
	}
	return f
}
