package cerrors

import (
	"context"
	"reflect"
	"testing"
)

func TestContextField(t *testing.T) {
	t.Run("without ctx", func(t *testing.T) {
		fields := CtxFields(nil)
		if fields != nil {
			t.Error("not nil fields")
		}
	})

	t.Run("context without fields", func(t *testing.T) {
		fields := CtxFields(context.Background())
		if fields != nil {
			t.Error("not nil fields")
		}
	})

	t.Run("context with fields", func(t *testing.T) {
		ctx := WithCtxField(context.Background(), "userId", 1)
		ctx = WithCtxField(ctx, "handler", "addUser")

		fields := CtxFields(ctx)

		expectedFields := map[string]any{"userId": 1, "handler": "addUser"}
		if !reflect.DeepEqual(expectedFields, fields) {
			t.Errorf("unexpected fields: expected %v, got %v", expectedFields, fields)
		}
	})

	t.Run("nil context", func(t *testing.T) {
		ctx := WithCtxField(nil, "userId", 1)
		if ctx == nil {
			t.Error("cxt is nil")
		}

		fields := CtxFields(ctx)

		expectedFields := map[string]any{"userId": 1}
		if !reflect.DeepEqual(expectedFields, fields) {
			t.Errorf("unexpected fields: expected %v, got %v", expectedFields, fields)
		}
	})

	t.Run("overwrite values", func(t *testing.T) {
		ctx := WithCtxField(context.Background(), "userId", 1)
		ctx = WithCtxField(ctx, "userId", "some id")

		fields := CtxFields(ctx)

		expectedFields := map[string]any{"userId": "some id"}
		if !reflect.DeepEqual(expectedFields, fields) {
			t.Errorf("unexpected fields: expected %v, got %v", expectedFields, fields)
		}
	})

}
