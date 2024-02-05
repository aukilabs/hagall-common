package errors

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Run("new error", func(t *testing.T) {
		err := New("hello").(richError)
		require.Equal(t, "hello", err.Message())
		require.Equal(t, "errors_test.go:14", err.Line())
		t.Log(err)
	})

	t.Run("new error with format", func(t *testing.T) {
		err := Newf("hello %v", 42).(richError)
		require.Equal(t, "hello 42", err.Message())
		require.Equal(t, "errors_test.go:21", err.Line())
		t.Log(err)
	})
}

func TestUnwrap(t *testing.T) {
	t.Run("enriched error is unwraped", func(t *testing.T) {
		werr := fmt.Errorf("werr")
		err := New("err").Wrap(werr)
		require.Equal(t, werr, Unwrap(err))
	})

	t.Run("enriched error is not unwraped", func(t *testing.T) {
		err := New("err")
		require.Nil(t, Unwrap(err))
	})
}

func TestIs(t *testing.T) {
	t.Run("is enriched error is true", func(t *testing.T) {
		err := New("test")
		require.True(t, Is(err, err))
	})

	t.Run("is enriched error is false", func(t *testing.T) {
		err := New("test")
		require.False(t, Is(err, New("test b")))
	})

	t.Run("is nested enriched error is true", func(t *testing.T) {
		werr := New("werr")
		err := fmt.Errorf("err: %w", werr)
		require.True(t, Is(err, werr))
	})

	t.Run("is nested enriched error is false", func(t *testing.T) {
		werr := New("werr")
		err := fmt.Errorf("err: %w", New("werr"))
		require.False(t, Is(err, werr))
	})

	t.Run("is not enriched nested error is true", func(t *testing.T) {
		werr := fmt.Errorf("werr")
		err := New("err").Wrap(werr)
		require.True(t, Is(err, werr))
	})

	t.Run("is not enriched nested error is false", func(t *testing.T) {
		werr := fmt.Errorf("werr")
		err := New("err").Wrap(fmt.Errorf("werr"))
		require.False(t, Is(err, werr))
	})
}

func TestAs(t *testing.T) {
	t.Run("has enriched error is true", func(t *testing.T) {
		var ierr Error
		err := New("err")
		require.True(t, As(err, &ierr))
	})

	t.Run("has not enriched error is false", func(t *testing.T) {
		var ierr Error
		err := fmt.Errorf("err")
		require.False(t, As(err, &ierr))
	})

	t.Run("has nested enriched error is true", func(t *testing.T) {
		var ierr Error
		err := fmt.Errorf("err: %w", New("werr"))
		require.True(t, As(err, &ierr))
	})
}

func TestIsType(t *testing.T) {
	t.Run("nil error is empty", func(t *testing.T) {
		require.True(t, IsType(nil, ""))
	})

	t.Run("enriched error is of the default type", func(t *testing.T) {
		err := New("err")
		require.True(t, IsType(err, "errors.richError"))
	})

	t.Run("enriched error is of the defined type", func(t *testing.T) {
		err := New("err").WithType("foo")
		require.True(t, IsType(err, "foo"))
	})

	t.Run("enriched error is not of the requested type", func(t *testing.T) {
		err := New("err").WithType("foo")
		require.False(t, IsType(err, "bar"))
	})

	t.Run("non enriched error is of the default type", func(t *testing.T) {
		err := fmt.Errorf("err")
		require.True(t, IsType(err, "*errors.errorString"))
	})

	t.Run("non enriched error is not of the requested type", func(t *testing.T) {
		err := fmt.Errorf("err")
		require.False(t, IsType(err, "foo"))
	})

	t.Run("enriched error is of the nested enriched type", func(t *testing.T) {
		err := New("err").Wrap(New("werr").WithType("foo"))
		require.True(t, IsType(err, "foo"))
	})

	t.Run("enriched error is of the nested non enriched type", func(t *testing.T) {
		err := New("err").Wrap(fmt.Errorf("werr"))
		require.True(t, IsType(err, "*errors.errorString"))
	})

	t.Run("non enriched error is of the nested enriched type", func(t *testing.T) {
		err := fmt.Errorf("err: %w", New("werr").WithType("foo"))
		require.True(t, IsType(err, "foo"))
	})
}

func TestTag(t *testing.T) {
	t.Run("enriched error returns the tag value", func(t *testing.T) {
		err := New("test").WithTag("foo", "bar")
		require.Equal(t, "bar", Tag(err, "foo"))
	})

	t.Run("enriched error does not returns the tag value", func(t *testing.T) {
		err := New("test")
		require.Empty(t, Tag(err, "foo"))
	})

	t.Run("nested enriched error in enriched error returns the tag value", func(t *testing.T) {
		err := New("err").Wrap(New("werr").WithTag("foo", "bar"))
		require.Equal(t, "bar", Tag(err, "foo"))
	})

	t.Run("nested enriched error in non enriched error returns the tag value", func(t *testing.T) {
		err := fmt.Errorf("err: %w", New("werr").WithTag("foo", "bar"))
		require.Equal(t, "bar", Tag(err, "foo"))
	})

	t.Run("non enriched error does not returns the tag value", func(t *testing.T) {
		err := fmt.Errorf("err")
		require.Empty(t, Tag(err, "foo"))
	})
}

func TestTags(t *testing.T) {
	t.Run("error without tags returns nil", func(t *testing.T) {
		err := New("err")
		require.Nil(t, err.Tags())
	})

	t.Run("error tags are returned", func(t *testing.T) {
		err := New("err").WithTag("foo", "bar")
		require.Equal(t, map[string]string{"foo": "bar"}, err.Tags())
	})
}

func TestError(t *testing.T) {
	SetIndentEncoder()
	defer SetInlineEncoder()

	t.Run("stringify an enriched error", func(t *testing.T) {
		err := New("err").
			WithTag("foo", "bar").
			Error()
		require.Contains(t, err, "err")
		require.Contains(t, err, "errors.richError")
		t.Log(err)
	})

	t.Run("stringify an enriched error wrapped in an enriched error", func(t *testing.T) {
		err := New("err").
			WithTag("foo", "bar").
			Wrap(New("werr").WithType("boo")).
			Error()
		require.Contains(t, err, "err")
		require.NotContains(t, err, "errors.richError")
		require.Contains(t, err, "werr")
		require.Contains(t, err, "boo")
		t.Log(err)
	})

	t.Run("stringify a non enriched error wrapped in an enriched error", func(t *testing.T) {
		err := New("err").
			WithTag("foo", "bar").
			Wrap(fmt.Errorf("werr")).
			Error()
		require.Contains(t, err, "err")
		require.NotContains(t, err, "errors.richError")
		require.Contains(t, err, "werr")
		require.Contains(t, err, "*errors.errorString")
		t.Log(err)

	})

	t.Run("stringify a non enriched error wrapped in an enriched error", func(t *testing.T) {
		err := fmt.Errorf("err: %w", New("werr")).Error()
		require.Contains(t, err, "werr")
		require.NotContains(t, err, "*errors.errorString")
		require.Contains(t, err, "err")
		require.Contains(t, err, "errors.richError")
		t.Log(err)
	})
}

func TestMessage(t *testing.T) {
	t.Run("enriched error message is returned", func(t *testing.T) {
		err := New("hello").WithTag("name", "buu")
		require.Equal(t, "hello", Message(err))
	})

	t.Run("standard error message is returned", func(t *testing.T) {
		err := fmt.Errorf("hello world")
		require.Equal(t, "hello world", Message(err))
	})
}

func TestToString(t *testing.T) {
	utests := []struct {
		in  interface{}
		out string
	}{
		{
			in:  "hello",
			out: "hello",
		},
		{
			in:  []byte("bye"),
			out: "bye",
		},
		{
			in:  -42,
			out: "-42",
		},
		{
			in:  int64(-42),
			out: "-42",
		},
		{
			in:  int32(-42),
			out: "-42",
		},
		{
			in:  int16(-42),
			out: "-42",
		},
		{
			in:  int8(-42),
			out: "-42",
		},
		{
			in:  uint(84),
			out: "84",
		},
		{
			in:  uint64(84),
			out: "84",
		},
		{
			in:  uint32(84),
			out: "84",
		},
		{
			in:  uint16(84),
			out: "84",
		},
		{
			in:  uint8(84),
			out: "84",
		},
		{
			in:  42.42,
			out: "42.42",
		},
		{
			in:  float32(42.42),
			out: "42.42",
		},
		{
			in:  true,
			out: "true",
		},
		{
			in:  false,
			out: "false",
		},
		{
			in:  time.Minute,
			out: "1m0s",
		},
		{
			in:  map[string]string{"foo": "bar"},
			out: `{"foo":"bar"}`,
		},
	}

	for _, u := range utests {
		t.Run(reflect.TypeOf(u.in).String(), func(t *testing.T) {
			require.Equal(t, u.out, toString(u.in))
		})
	}
}
