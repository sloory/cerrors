package cerrors

import (
	"errors"
	"fmt"
	"regexp"
	"runtime"
	"strings"
	"testing"
)

var initpc = caller()

type X struct{}

// val returns a Frame pointing to itself.
func (x X) val() Frame {
	return caller()
}

// ptr returns a Frame pointing to itself.
func (x *X) ptr() Frame {
	return caller()
}

func TestFrameFormat(t *testing.T) {
	var tests = []struct {
		Frame
		format string
		want   string
	}{{
		initpc,
		"%s",
		"stacktrace_test.go",
	}, {
		initpc,
		"%+s",
		"github.com/sloory/cerrors.init\n" +
			"\t.+/cerrors/stacktrace_test.go",
	}, {
		0,
		"%s",
		"unknown",
	}, {
		0,
		"%+s",
		"unknown",
	}, {
		initpc,
		"%v",
		"stacktrace_test.go:12",
	}, {
		initpc,
		"%+v",
		"github.com/sloory/cerrors.init\n" +
			"\t.+/cerrors/stacktrace_test.go:12",
	}, {
		0,
		"%v",
		"unknown:0",
	}}

	for i, tt := range tests {
		testFormatRegexp(t, i, tt.Frame, tt.format, tt.want)
	}
}

func TestStackTrace(t *testing.T) {

	tests := []struct {
		err  error
		want []string
	}{{
		WithStack(errors.New("ooh")), []string{
			"github.com/sloory/cerrors.TestStackTrace\n" +
				"\t.+/cerrors/stacktrace_test.go:74",
		},
	}, {
		Wrap("ahh",
			WithStack(
				errors.New("ooh"),
			),
		), []string{
			"github.com/sloory/cerrors.TestStackTrace\n" +
				"\t.+/cerrors/stacktrace_test.go:80", // this is the stack of newWithStack, not New
		},
	}, {
		func() error {
			return WithStack(errors.New("ooh"))
		}(), []string{
			`github.com/sloory/cerrors.TestStackTrace.func1` +
				"\n\t.+/cerrors/stacktrace_test.go:89", // this is the stack of newWithStack
			"github.com/sloory/cerrors.TestStackTrace\n" +
				"\t.+/cerrors/stacktrace_test.go:90", // this is the stack of newWithStack's caller
		},
	}}
	for i, tt := range tests {

		var errWithStack interface {
			StackTrace() StackTrace
		}

		if !errors.As(tt.err, &errWithStack) {
			t.Errorf("expected %#v to implement StackTrace() StackTrace", tt.err)
			continue
		}
		st := errWithStack.StackTrace()
		for j, want := range tt.want {
			testFormatRegexp(t, i, st[j], "%+v", want)
		}
	}
}

func stackTraceTest() StackTrace {
	const depth = 8
	var pcs [depth]uintptr
	n := runtime.Callers(1, pcs[:])
	f := make([]Frame, n)
	for i := 0; i < n; i++ {
		f[i] = Frame(pcs[i])
	}
	return f
}

func TestStackTraceFormat(t *testing.T) {
	tests := []struct {
		StackTrace
		format string
		want   string
	}{{
		nil,
		"%s",
		`\[\]`,
	}, {
		nil,
		"%v",
		`\[\]`,
	}, {
		nil,
		"%+v",
		"",
	}, {
		make(StackTrace, 0),
		"%s",
		`\[\]`,
	}, {
		make(StackTrace, 0),
		"%v",
		`\[\]`,
	}, {
		make(StackTrace, 0),
		"%+v",
		"",
	}, {
		stackTraceTest()[:2],
		"%s",
		`\[stacktrace_test.go stacktrace_test.go\]`,
	}, {
		stackTraceTest()[:2],
		"%v",
		`\[stacktrace_test.go:117 stacktrace_test.go:159\]`,
	}, {
		stackTraceTest()[:2],
		"%+v",
		"\n" +
			"github.com/sloory/cerrors.stackTraceTest\n" +
			"\t.+/cerrors/stacktrace_test.go:117\n" +
			"github.com/sloory/cerrors.TestStackTraceFormat\n" +
			"\t.+/cerrors/stacktrace_test.go:163",
	}, {
		stackTraceTest()[:2],
		"%#v",
		`\[\]cerrors.Frame{stacktrace_test.go:117, stacktrace_test.go:171}`,
	}}

	for i, tt := range tests {
		testFormatRegexp(t, i, tt.StackTrace, tt.format, tt.want)
	}
}

// a version of runtime.Caller that returns a Frame, not a uintptr.
func caller() Frame {
	var pcs [3]uintptr
	n := runtime.Callers(2, pcs[:])
	frames := runtime.CallersFrames(pcs[:n])
	frame, _ := frames.Next()
	return Frame(frame.PC)
}

func testFormatRegexp(t *testing.T, n int, arg interface{}, format, want string) {
	t.Helper()
	got := fmt.Sprintf(format, arg)
	gotLines := strings.SplitN(got, "\n", -1)
	wantLines := strings.SplitN(want, "\n", -1)

	if len(wantLines) > len(gotLines) {
		t.Errorf("test %d: wantLines(%d) > gotLines(%d):\n got: %q\nwant: %q", n+1, len(wantLines), len(gotLines), got, want)
		return
	}

	for i, w := range wantLines {
		match, err := regexp.MatchString(w, gotLines[i])
		if err != nil {
			t.Fatal(err)
		}
		if !match {
			t.Errorf("test %d: line %d: fmt.Sprintf(%q, err):\n got: %q\nwant: %q", n+1, i+1, format, got, want)
		}
	}
}
