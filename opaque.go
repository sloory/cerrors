package cerrors

import "fmt"

// check interface implementation
var _ error = (*opaqueError)(nil)
var _ fmt.Formatter = (*opaqueError)(nil)

type opaqueError struct {
	cause   error
	message string
}

func newOpaque(msg string, err error) error {
	if err == nil {
		return nil
	}

	return &opaqueError{cause: err, message: msg}
}

func (w *opaqueError) Error() string { return w.message }
func (w *opaqueError) Unwrap() error { return w.cause }
func (w *opaqueError) Format(f fmt.State, verb rune) {
	fmt.Printf(fmt.FormatString(f, verb), w.cause)
}
