package cerrors

import (
	"errors"
)

type withFields interface {
	error
	Fields() map[string]interface{}
	AddField(name string, value any)
	AddFields(fields map[string]interface{})
}

type withFieldsError struct {
	cause  error
	fields map[string]interface{}
}

func newWithFields(err error) withFields {
	if err == nil {
		return nil
	}

	var fErr *withFieldsError
	if errors.As(err, &fErr) {
		return fErr
	}

	return &withFieldsError{cause: err, fields: make(map[string]interface{})}
}

func (w *withFieldsError) AddField(name string, value any) {
	w.fields[name] = value
}

func (w *withFieldsError) AddFields(fields map[string]interface{}) {
	for name, value := range fields {
		w.fields[name] = value
	}
}

func (w *withFieldsError) Error() string                  { return w.cause.Error() }
func (w *withFieldsError) Unwrap() error                  { return w.cause }
func (w *withFieldsError) Fields() map[string]interface{} { return w.fields }
