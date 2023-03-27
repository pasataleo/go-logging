package logging

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/pasataleo/go-errors/errors"
	"github.com/pasataleo/go-testing/tests"
)

func TestJson(t *testing.T) {
	var logs bytes.Buffer

	logger := JsonLogger()
	logger.Out = &logs
	logger.DateTimeFormat = "2006/01/02"

	tests.ExecFn(t, logger.Infof, "this is my %s", "message").
		Fatal().
		NoError()

	tests.ExecFn(t, logs.String).
		Fatal().
		Equals(fmt.Sprintf("{\"datetime\":\"%s\",\"level\":\"info\",\"message\":\"this is my message\"}\n", time.Now().Format(logger.DateTimeFormat)))
}

func TestJson_WithError(t *testing.T) {
	var logs bytes.Buffer

	json := JsonLogger()
	json.Out = &logs
	json.DateTimeFormat = "2006-01-02"

	logger := json.WithError(errors.New(nil, errors.ErrorCodeUnknown, "bad error"))

	tests.ExecFn(t, logger.Errorf, "this is my %s", "message").
		Fatal().
		NoError()

	tests.ExecFn(t, logs.String).
		Fatal().
		Equals(fmt.Sprintf("{\"datetime\":\"%s\",\"error\":\"bad error\",\"level\":\"error\",\"message\":\"this is my message\"}\n", time.Now().Format(json.DateTimeFormat)))

}

func TestJson_WithField(t *testing.T) {
	var logs bytes.Buffer

	json := JsonLogger()
	json.Out = &logs
	json.DateTimeFormat = "2006-01-02"

	logger := json.WithField("key", "value")

	tests.ExecFn(t, logger.Errorf, "this is my %s", "message").
		Fatal().
		NoError()

	tests.ExecFn(t, logs.String).
		Fatal().
		Equals(fmt.Sprintf("{\"datetime\":\"%s\",\"key\":\"value\",\"level\":\"error\",\"message\":\"this is my message\"}\n", time.Now().Format(json.DateTimeFormat)))
}
