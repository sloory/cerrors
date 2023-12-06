package cerrors

import (
	//	"context"
	"context"
	"errors"
	"reflect"
	"testing"
)

var (
	testFields = map[string]interface{}{
		"key": "value",
	}
	// testCtx = log.WithFields(context.Background(), testFields)
)

// func TestEnrich(t *testing.T) {
// 	ctx := context.Background()

// 	t.Run("nil_err", func(t *testing.T) {
// 		err := Enrich(ctx, nil)

// 		require.Equal(t, nil, err)
// 	})

// 	t.Run("ctx_without_fields", func(t *testing.T) {
// 		err := Enrich(ctx, errors.New("err"))

// 		requireNoFields(t, err)
// 		requireStack(t, err)
// 	})

// 	expFields := map[string]interface{}{
// 		"itemID":    1,
// 		"requestID": 2,
// 	}

// 	ctxWithFields := log.WithFields(ctx, expFields)

// 	t.Run("ctx_test", func(t *testing.T) {
// 		err := Enrich(ctxWithFields, errors.New("err"))
// 		err = Enrich(ctxWithFields, err)

// 		fieldsErr := requireFields(t, err)
// 		require.Equal(t, expFields, fieldsErr.Fields())

// 		i := 0
// 		for err := err; err != nil; err = errors.Unwrap(err) {
// 			i++
// 		}

// 		require.Equalf(t, 3, i, "Ошибка должна быть обернута только в контекст и стэктрейс.")
// 	})

// 	t.Run("stacktrace_test", func(t *testing.T) {
// 		err := Enrich(ctxWithFields, errors.New("err"))

// 		var stErr stackTrace

// 		require.ErrorAs(t, err, &stErr)

// 		t.Logf("%+v", stErr.StackTrace())
// 		require.Len(t, stErr.StackTrace(), 3)
// 	})

// 	expFields = map[string]interface{}{
// 		"itemID":    1,
// 		"requestID": 2,
// 		"phoneID":   3,
// 	}

// 	ctxWithFields = log.WithFields(ctxWithFields, map[string]interface{}{
// 		"phoneID": 3,
// 	})

// 	t.Run("ext_ctx_test", func(t *testing.T) {
// 		err := Enrich(ctxWithFields, errors.New("err"))

// 		fieldsErr := requireFields(t, err)
// 		require.EqualValues(t, expFields, fieldsErr.Fields())
// 	})
// }

func TestWithStack(t *testing.T) {

	t.Run("with stack", func(t *testing.T) {
		err := WithStack(errors.New("err"))
		var errWithStack stackTrace
		if !errors.As(err, &errWithStack) || len(errWithStack.StackTrace()) == 0 {
			t.Error("after EnrichStack error must have stack")
		}
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

// // func TestEnrichContext(t *testing.T) {
// // 	var err error

// // 	ctx := context.Background()

// // 	t.Run("nil_err", func(t *testing.T) {
// // 		got := EnrichContext(ctx, err)

// // 		assert.Equal(t, got, context.Background())
// // 	})

// // 	err = errors.New("error")

// // 	t.Run("not_fields_err", func(t *testing.T) {
// // 		got := EnrichContext(ctx, err)

// // 		assert.Equal(t, got, context.Background())
// // 	})

// // 	expFields := map[string]interface{}{
// // 		"itemID":    1,
// // 		"requestID": 2,
// // 	}

// // 	ctxWithFields := log.WithFields(ctx, expFields)

// // 	err = Enrich(ctxWithFields, err)

// // 	t.Run("fields_err", func(t *testing.T) {
// // 		got := EnrichContext(ctx, err)

// // 		f := log.Fields(got)
// // 		assert.Equal(t, expFields, f)
// // 	})
// // }

func TestWithField(t *testing.T) {
	t.Run("ordinal error", func(t *testing.T) {
		err := WithField(errors.New("err"), "some key", "field value")

		var fieldsErr withFields
		if !errors.As(err, &fieldsErr) {
			t.Error("expect error with fields")
		}

		expectedFields := map[string]any{"some key": "field value"}
		if !reflect.DeepEqual(expectedFields, fieldsErr.Fields()) {
			t.Errorf("unexpected fields: expected %v, got %v", expectedFields, fieldsErr.Fields())
		}
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

		var fieldsErr withFields
		if !errors.As(err2, &fieldsErr) {
			t.Error("expect error with fields")
		}

		expectedFields := map[string]any{"key1": "value1", "key2": "value2"}
		if !reflect.DeepEqual(expectedFields, fieldsErr.Fields()) {
			t.Errorf("unexpected fields: expected %v, got %v", expectedFields, fieldsErr.Fields())
		}
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

		var fieldsErr withFields
		if !errors.As(err2, &fieldsErr) {
			t.Error("expect error with fields")
		}

		expectedFields := map[string]any{
			"key1": "value1",
			"key2": "NEW",
			"key3": "value4",
		}
		if !reflect.DeepEqual(expectedFields, fieldsErr.Fields()) {
			t.Errorf("unexpected fields: expected %v, got %v", expectedFields, fieldsErr.Fields())
		}
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

// func requireStack(t *testing.T, err error) {
// 	var stackErr stackTrace
// 	require.ErrorAs(t, err, &stackErr)
// }

// func requireFields(t *testing.T, err error) withFields {
// 	var fieldsErr withFields
// 	require.ErrorAs(t, err, &fieldsErr)

// 	return fieldsErr
// }

// func requireNoFields(t *testing.T, err error) {
// 	var fieldsErr withFields

// 	errors.As(err, &fieldsErr)
// 	require.Equal(t, nil, fieldsErr)
// }
