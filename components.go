package cerrors

import (
	"context"
	"errors"
	"fmt"
)

type key int

const (
	componentsKey key = 1
)

type withComponents interface {
	error
	fmt.Formatter
	
	Components() []string
}

// check interface implementation
var _ withComponents = (*withComponentsError)(nil)

type withComponentsError struct {
	cause      error
	components []string
}

func enrichWithComponents(ctx context.Context, err error) error {
	if err == nil {
		return nil
	}

	var cErr withComponents
	if errors.As(err, &cErr) {
		return err
	}

	comps := getCtxComponents(ctx)
	if len(comps) == 0 {
		return err
	}

	return &withComponentsError{cause: err, components: comps}
}

func (w *withComponentsError) Error() string        { return w.cause.Error() }
func (w *withComponentsError) Unwrap() error        { return w.cause }
func (w *withComponentsError) Components() []string { return w.components }
func (w *withComponentsError) Format(f fmt.State, verb rune) {
	fmt.Printf(fmt.FormatString(f, verb), w.cause)
}

func getCtxComponents(ctx context.Context) []string {
	c, ok := ctx.Value(componentsKey).([]string)
	if !ok {
		return nil
	}

	return c
}

func InComponent(ctx context.Context, component string) context.Context {
	components := getCtxComponents(ctx)
	components = append(components, component)

	return context.WithValue(ctx, componentsKey, components)
}
