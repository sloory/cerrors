package cerrors

import (
	"context"
	"errors"
	"fmt"
)

func Enrich(ctx context.Context, err error) error {
	if ctx == nil {
		return newWithStack(err)
	}

	return enrichWithComponents(ctx, newWithStack(err))
	// newWithFields( // replace on enrichWithFields
	// ),
}

func WithStack(err error) error {
	return newWithStack(err)
}

func Wrap(msg string, err error) error {
	if err == nil {
		return nil
	}

	return fmt.Errorf("%v: %w", msg, err)
}

func Opaque(msg string, err error) error {
	if err == nil {
		return nil
	}

	return newOpaque(msg, err)
}

func Nested(parent, child error) error {
	return fmt.Errorf("%w: %w", parent, child)
}

func WithField(err error, key string, value any) error {
	if err == nil {
		return nil
	}

	fErr := newWithFields(err)
	fErr.AddField(key, value)
	return fErr
}

func WithFields(err error, fields map[string]any) error {
	if err == nil {
		return nil
	}

	fErr := newWithFields(err)
	fErr.AddFields(fields)
	return fErr
}

func Fields(err error) map[string]any {
	if err == nil {
		return nil
	}

	var fErr *withFieldsError
	if !errors.As(err, &fErr) {
		return nil
	}

	return fErr.Fields()
}

func Components(err error) []string {
	if err == nil {
		return nil
	}

	var cErr withComponents
	if errors.As(err, &cErr) {
		return cErr.Components()
	}

	return nil
}

// func EnrichContext(ctx context.Context, err error) context.Context {
// 	var errFields withFields
// 	if errors.As(err, &errFields) {
// 		return log.WithFields(ctx, errFields.Fields())
// 	}

// 	return ctx
// }
