package cerrors

import (
	"context"
)

const (
	fieldsKey key = 2
)

func WithCtxField(ctx context.Context, k string, v any) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	fields, ok := ctx.Value(fieldsKey).(map[string]any)
	if !ok || fields == nil {
		fields = make(map[string]any)
	}

	fields[k] = v
	return context.WithValue(ctx, fieldsKey, fields)
}

func CtxFields(ctx context.Context) map[string]any {
	if ctx == nil {
		return nil
	}

	fields, ok := ctx.Value(fieldsKey).(map[string]any)
	if !ok {
		return nil
	}

	return fields
}

func enrichWithFields(ctx context.Context, err error) error {
	if err == nil {
		return nil
	}

	fields := CtxFields(ctx)
	if fields == nil {
		return err
	}

	fErr := newWithFields(err)
	if len(fields) == 0 {
		return err
	}

	for key, value := range fields {
		fErr.AddField(key, value)
	}

	return fErr
}
