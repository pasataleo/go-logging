package logging

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/pasataleo/go-errors/errors"
	"github.com/pasataleo/go-testing/tests"
)

func TestText(t *testing.T) {
	var logs bytes.Buffer

	logger := TextLogger()
	logger.Out = &logs
	logger.DateTimeFormat = "2006/01/02"

	tests.ExecFn(t, logger.Infof, "this is my %s", "message").
		Fatal().
		NoError()

	tests.ExecFn(t, logs.String).
		Fatal().
		Equals(fmt.Sprintf("[info] %s: this is my message\n", time.Now().Format(logger.DateTimeFormat)))
}

func TestText_WithError(t *testing.T) {
	var logs bytes.Buffer

	txt := TextLogger()
	txt.Out = &logs
	txt.DateTimeFormat = "2006-01-02"

	logger := txt.WithError(errors.New(nil, errors.ErrorCodeUnknown, "bad error"))

	tests.ExecFn(t, logger.Errorf, "this is my %s", "message").
		Fatal().
		NoError()

	tests.ExecFn(t, logs.String).
		Fatal().
		Equals(fmt.Sprintf("[error] %s: (bad error) this is my message\n", time.Now().Format(txt.DateTimeFormat)))
}

func TestText_WithField(t *testing.T) {
	var logs bytes.Buffer

	txt := TextLogger()
	txt.Out = &logs
	txt.DateTimeFormat = "2006-01-02"

	logger := txt.WithField("key", "value")

	tests.ExecFn(t, logger.Errorf, "this is my %s", "message").
		Fatal().
		NoError()

	tests.ExecFn(t, logs.String).
		Fatal().
		Equals(fmt.Sprintf("[error] %s: this is my message (value)\n", time.Now().Format(txt.DateTimeFormat)))
}
