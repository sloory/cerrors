package cerrors

import (
	"errors"

	"github.com/cockroachdb/errors/errbase"
	"github.com/cockroachdb/errors/withstack"
)

type stackTrace interface {
	error
	StackTrace() errbase.StackTrace
}

func newWithStack(err error) error {
	if err == nil {
		return nil
	}

	var stErr stackTrace
	if errors.As(err, &stErr) {
		return err
	}

	return withstack.WithStackDepth(err, 2) //nolint:gomnd // self-explained
}
