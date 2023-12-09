package cerrors

import (
	"context"
	"errors"
	"reflect"
	"testing"
)

func TestEnrich(t *testing.T) {
	t.Run("nil error", func(t *testing.T) {
		err := Enrich(context.Background(), nil)
		if err != nil {
			t.Error("not nil error")
		}
	})

	t.Run("stacktrace for nil context", func(t *testing.T) {
		err := Enrich(nil, errors.New("err"))
		requireStack(t, err)
	})

	t.Run("stacktrace", func(t *testing.T) {
		err := Enrich(context.Background(), errors.New("err"))
		requireStack(t, err)
	})

	t.Run("ctx without fields", func(t *testing.T) {
		err := Enrich(context.Background(), errors.New("err"))

		var fieldsErr withFields
		if errors.As(err, &fieldsErr) {
			t.Error("expect error without fields")
		}
	})

	t.Run("ctx with fields", func(t *testing.T) {
		ctx := WithCtxField(context.Background(), "itemId", 11)
		ctx = WithCtxField(ctx, "requestId", "some text value")

		err := Enrich(ctx, errors.New("err"))

		var fieldsErr withFields
		if !errors.As(err, &fieldsErr) {
			t.Error("expect error with fields")
		}

		expected := map[string]interface{}{
			"itemId":    11,
			"requestId": "some text value",
		}

		if !reflect.DeepEqual(expected, fieldsErr.Fields()) {
			t.Errorf("unexpected fields: expected %v, got %v", expected, fieldsErr.Fields())
		}
	})

	t.Run("several calls - errors wrap not grow", func(t *testing.T) {
		ctx := WithCtxField(context.Background(), "itemId", 11)
		ctx = InComponent(ctx, "handler")

		err := Enrich(ctx, errors.New("err"))
		err = Enrich(ctx, err)
		err = Enrich(ctx, err)

		wrapsCount := 0
		for err := err; err != nil; err = errors.Unwrap(err) {
			wrapsCount++
		}

		// original, withStack, withFields, withComponents
		expected := 4
		if wrapsCount != expected {
			t.Errorf("unexpected wraps: expected %d, got %d", expected, wrapsCount)
		}
	})

	innerFunc1 := func() error {
		return Enrich(context.Background(), errors.New("err"))
	}

	innerFunc := func() error {
		return innerFunc1()
	}

	t.Run("stacktrace length", func(t *testing.T) {
		err := innerFunc()

		var stErr stackTrace
		if !errors.As(err, &stErr) {
			t.Error("expect error with stacktrace")
		}

		expected := 5
		if len(stErr.StackTrace()) != expected {
			t.Errorf("unexpected stack length: expected %d, got %d", expected, len(stErr.StackTrace()))
		}
	})
}

func TestWithStack(t *testing.T) {

	t.Run("with stack", func(t *testing.T) {
		err := WithStack(errors.New("err"))
		requireStack(t, err)
	})

	t.Run("nil", func(t *testing.T) {
		err := WithStack(nil)
		if err != nil {
			t.Error("not nil error")
		}
	})
}

func TestWrap(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		err := Wrap("some text", nil)
		if err != nil {
			t.Error("not nil error")
		}
	})

	t.Run("with message", func(t *testing.T) {
		initialError := errors.New("record not found")
		err := Wrap("user service", initialError)

		if !errors.Is(err, initialError) {
			t.Error("do not match initial error")
		}

		expected := "user service: record not found"
		if err.Error() != expected {
			t.Errorf("unexpected error message: expected %v, got %v", expected, err.Error())
		}
	})
}

func TestOpaque(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		err := Opaque("some text", nil)
		if err != nil {
			t.Error("not nil error")
		}
	})

	t.Run("error message hide", func(t *testing.T) {
		initialError := errors.New("record not found")
		err := Opaque("user not found", initialError)

		if !errors.Is(err, initialError) {
			t.Error("do not match initial error")
		}

		expected := "user not found"
		if err.Error() != expected {
			t.Errorf("unexpected error message: expected %v, got %v", expected, err.Error())
		}
	})
}

func TestInComponent(t *testing.T) {
	t.Run("one", func(t *testing.T) {
		ctx := InComponent(context.Background(), "api")

		err := Enrich(ctx, errors.New("error"))

		expected := []string{"api"}
		if !reflect.DeepEqual(expected, Components(err)) {
			t.Errorf("unexpected fields: expected %v, got %v", expected, err.Error())
		}
	})

	t.Run("several", func(t *testing.T) {
		ctx := InComponent(context.Background(), "api")
		ctx = InComponent(ctx, "service")
		ctx = InComponent(ctx, "storage")

		err := Enrich(ctx, errors.New("error"))

		expected := []string{"api", "service", "storage"}
		if !reflect.DeepEqual(expected, Components(err)) {
			t.Errorf("unexpected fields: expected %v, got %v", expected, err.Error())
		}
	})
}

func TestWithField(t *testing.T) {
	t.Run("ordinal error", func(t *testing.T) {
		err := WithField(errors.New("err"), "some key", "field value")

		requireFields(t, err, map[string]any{"some key": "field value"})
	})

	t.Run("nil", func(t *testing.T) {
		err := WithField(nil, "some key", "field value")
		if err != nil {
			t.Error("not nil error")
		}
	})

	t.Run("several fields", func(t *testing.T) {
		err1 := WithFields(
			errors.New("err"),
			map[string]any{"key1": "value1"},
		)

		err2 := WithField(err1, "key2", "value2")

		requireFields(t, err2, map[string]any{"key1": "value1", "key2": "value2"})
	})

	t.Run("previous fields overwritten", func(t *testing.T) {
		err1 := WithFields(
			errors.New("err"),
			map[string]any{
				"key1": "value1",
				"key2": "value2",
				"key3": "value4",
			},
		)

		err2 := WithField(err1, "key2", "NEW")

		requireFields(t, err2, map[string]any{
			"key1": "value1",
			"key2": "NEW",
			"key3": "value4",
		})
	})
}

func TestFields(t *testing.T) {
	t.Run("ordinal error", func(t *testing.T) {
		err := Fields(errors.New("err"))
		if err != nil {
			t.Error("not nil fields")
		}
	})

	t.Run("nil", func(t *testing.T) {
		err := Fields(nil)
		if err != nil {
			t.Error("not nil fields")
		}
	})

	t.Run("error with fields", func(t *testing.T) {
		err := WithFields(
			errors.New("err"),
			map[string]any{"key1": "value1", "key2": "value2"},
		)

		expectedFields := map[string]any{"key1": "value1", "key2": "value2"}
		if !reflect.DeepEqual(expectedFields, Fields(err)) {
			t.Errorf("unexpected fields: expected %v, got %v", expectedFields, Fields(err))
		}
	})
}

func requireStack(t *testing.T, err error) {
	var errWithStack stackTrace
	if !errors.As(err, &errWithStack) || len(errWithStack.StackTrace()) == 0 {
		t.Error("after WithStack error must have stack")
	}
}

func requireFields(t *testing.T, err error, expected map[string]any) {
	var fieldsErr withFields
	if !errors.As(err, &fieldsErr) {
		t.Error("expect error with fields")
	}

	if !reflect.DeepEqual(expected, fieldsErr.Fields()) {
		t.Errorf("unexpected fields: expected %v, got %v", expected, fieldsErr.Fields())
	}
}
