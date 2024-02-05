package logs

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/aukilabs/hagall-common/errors"
	"github.com/stretchr/testify/require"
)

func TestParseLevel(t *testing.T) {
	utests := []struct {
		in  string
		out Level
	}{
		{
			in:  "debug",
			out: DebugLevel,
		},
		{
			in:  "info",
			out: InfoLevel,
		},
		{
			in:  "warning",
			out: WarningLevel,
		},
		{
			in:  "error",
			out: ErrorLevel,
		},
		{
			in:  "42",
			out: 42,
		},
	}

	for _, u := range utests {
		t.Run(u.in, func(t *testing.T) {
			require.Equal(t, u.out, ParseLevel(u.in))
		})
	}
}

func TestLevelString(t *testing.T) {
	utests := []struct {
		in  Level
		out string
	}{
		{
			in:  DebugLevel,
			out: "debug",
		},
		{
			in:  InfoLevel,
			out: "info",
		},
		{
			in:  WarningLevel,
			out: "warning",
		},
		{
			in:  ErrorLevel,
			out: "error",
		},
		{
			in:  42,
			out: "42",
		},
	}

	for _, u := range utests {
		t.Run(u.in.String(), func(t *testing.T) {
			require.Equal(t, u.out, u.in.String())
		})
	}
}

func TestSetLevel(t *testing.T) {
	SetLevel(WarningLevel)

	for i := DebugLevel; i < WarningLevel; i++ {
		require.Equal(t, fmt.Sprintf("%p", empytLogger), fmt.Sprintf("%p", loggers[i]))
	}

	for i := WarningLevel; i <= ErrorLevel; i++ {
		require.Equal(t, fmt.Sprintf("%p", logger), fmt.Sprintf("%p", loggers[i]))
	}
}

func TestEntry(t *testing.T) {
	SetIndentEncoder()
	SetLogger(func(e Entry) {
		t.Log(e)
	})

	t.Run("debug", func(t *testing.T) {
		entry := New()
		entry.Debug("logging a entry")
	})

	t.Run("debugf", func(t *testing.T) {
		New().Debugf("logging a entry %v", 42)
	})

	t.Run("info", func(t *testing.T) {
		New().Info("logging a entry")
	})

	t.Run("infof", func(t *testing.T) {
		New().Infof("logging a entry %v", 42)
	})

	t.Run("warn", func(t *testing.T) {
		New().Warn("logging a entry")
	})

	t.Run("warnf", func(t *testing.T) {
		New().Warnf("logging a entry %v", 42)
	})

	t.Run("error", func(t *testing.T) {
		New().Error(fmt.Errorf("simulated error"))
	})

	t.Run("rich error", func(t *testing.T) {
		New().Error(errors.New("simulated error").WithTag("style", "red"))

		New().Error(errors.New("simulated error").
			WithTag("style", "red").
			Wrap(fmt.Errorf("sub error")))

		WithTag("foo", "bar").Error(errors.New("simulated error").
			WithTag("style", "red").
			Wrap(fmt.Errorf("sub error")))
	})
}

func TestEntryTime(t *testing.T) {
	SetLogger(func(e Entry) {
		require.NotZero(t, e.Time())
		t.Log(e)
	})

	New().Debug("hi")
}

func TestEntryGetError(t *testing.T) {
	err := errors.New("test")

	SetLogger(func(e Entry) {
		require.Equal(t, err, e.GetError())
		t.Log(e)
	})

	Error(err)
}

func TestEntryTags(t *testing.T) {
	e := WithTag("hello", "max")
	require.Equal(t, map[string]any{"hello": "max"}, e.Tags())
}

func TestNormalizeTag(t *testing.T) {
	SetInlineEncoder()

	testValues := []struct {
		in  interface{}
		out any
	}{
		{
			in:  "hello",
			out: "hello",
		},
		{
			in:  fmt.Errorf("hi"),
			out: "hi",
		},
		{
			in:  []byte("bye"),
			out: "bye",
		},
		{
			in:  -42,
			out: -42,
		},
		{
			in:  int64(-42),
			out: int64(-42),
		},
		{
			in:  int32(-42),
			out: int32(-42),
		},
		{
			in:  int16(-42),
			out: int16(-42),
		},
		{
			in:  int8(-42),
			out: int8(-42),
		},
		{
			in:  uint(84),
			out: uint(84),
		},
		{
			in:  uint64(84),
			out: uint64(84),
		},
		{
			in:  uint32(84),
			out: uint32(84),
		},
		{
			in:  uint16(84),
			out: uint16(84),
		},
		{
			in:  uint8(84),
			out: uint8(84),
		},
		{
			in:  42.42,
			out: 42.42,
		},
		{
			in:  float32(42.42),
			out: float32(42.42),
		},
		{
			in:  true,
			out: true,
		},
		{
			in:  false,
			out: false,
		},
		{
			in:  time.Minute,
			out: "1m0s",
		},
		{
			in:  map[string]string{"foo": "bar"},
			out: map[string]string{"foo": "bar"},
		},
	}

	for _, u := range testValues {
		t.Run(reflect.TypeOf(u.in).String(), func(t *testing.T) {
			require.Equal(t, u.out, normalizeTag(u.in))
		})
	}
}

func TestEntryMarshalJSON(t *testing.T) {
	utest := []struct {
		in           any
		expectedJSON string
	}{
		{
			in:           42,
			expectedJSON: "42",
		},
		{
			in:           42.42,
			expectedJSON: "42.42",
		},
		{
			in:           "bar",
			expectedJSON: `"bar"`,
		},
		{
			in:           []int{42, 21},
			expectedJSON: `[42,21]`,
		},
		{
			in:           map[string]int{"n": 42},
			expectedJSON: `{"n":42}`,
		},
		{
			in: struct {
				Hello string
			}{
				Hello: "world",
			},
			expectedJSON: `{"Hello":"world"}`,
		},
	}

	for _, u := range utest {
		t.Run(fmt.Sprintf("marshal %T", u.in), func(t *testing.T) {
			entry := WithTag("foo", u.in)
			b, err := json.Marshal(entry)
			t.Log(string(b))

			require.NoError(t, err)
			require.Contains(t, string(b), fmt.Sprintf(`"foo":%s`, u.expectedJSON))
		})
	}

	t.Run("marshal enriched error without tags", func(t *testing.T) {
		entry := entry{
			err: errors.New("error test"),
		}

		_, err := json.Marshal(entry)
		require.NoError(t, err)
	})

	t.Run("marshal enriched error with tags", func(t *testing.T) {
		entry := entry{
			err: errors.New("error test").WithTag("foo", "bar"),
		}

		b, err := json.Marshal(entry)
		require.NoError(t, err)
		require.Contains(t, string(b), `"foo":"bar"`)
	})
}
