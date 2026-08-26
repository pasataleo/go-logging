package logging

import (
	"encoding/json"
	"testing"

	"github.com/pasataleo/go-testingx/pkg/testingx"
)

func TestLevel_SetStringIsCaseInsensitive(t *testing.T) {
	for _, name := range []string{"info", "INFO", "Info", "iNfO"} {
		var level Level
		testingx.Call(t, level.SetString, name).NoError()
		testingx.Capture(t, level).Equal(Info)
	}
}

func TestLevel_SetStringUnknown(t *testing.T) {
	var level Level
	testingx.Call(t, level.SetString, "nonsense").ErrorContains("unknown log level")
}

func TestLevel_TextRoundTrip(t *testing.T) {
	for _, level := range []Level{None, Trace, Debug, Info, Warn, Error, Fatal, Panic} {
		text := testingx.CallAs[[]byte](t, level.MarshalText).Fatal().NoError().Capture()
		testingx.Capture(t, string(text)).Equal(level.String())

		var decoded Level
		testingx.Call(t, decoded.UnmarshalText, text).NoError()
		testingx.Capture(t, decoded).Equal(level)
	}
}

// A Level field has to survive encoding/json, which is what logging/json.Decode
// depends on.
func TestLevel_JSONRoundTrip(t *testing.T) {
	type payload struct {
		Level Level `json:"level"`
	}

	data := testingx.CallAs[[]byte](t, json.Marshal, payload{Level: Warn}).Fatal().NoError().Capture()
	testingx.Capture(t, string(data)).Equal(`{"level":"warn"}`)

	var decoded payload
	testingx.Call(t, json.Unmarshal, data, &decoded).NoError()
	testingx.Capture(t, decoded.Level).Equal(Warn)
}

func TestLevel_JSONRejectsUnknown(t *testing.T) {
	var decoded struct {
		Level Level `json:"level"`
	}
	testingx.Call(t, json.Unmarshal, []byte(`{"level":"nonsense"}`), &decoded).Error()
}

func TestLowerASCII(t *testing.T) {
	for _, tc := range []struct{ in, want string }{
		{"", ""},
		{"info", "info"},
		{"INFO", "info"},
		{"InFo", "info"},
		{"log-level", "log-level"},
	} {
		testingx.Call(t, lowerASCII, tc.in).Equal(tc.want)
	}
}
