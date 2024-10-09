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

	tests.ExecuteE(logger.Infof("this is my %s", "message")).NoError(t)
	tests.Execute(logs.String()).Diff(t, fmt.Sprintf("[info] %s: this is my message\n", time.Now().Format(logger.DateTimeFormat)))
}

func TestText_WithError(t *testing.T) {
	var logs bytes.Buffer

	txt := TextLogger()
	txt.Out = &logs
	txt.DateTimeFormat = "2006-01-02"

	logger := txt.WithError(errors.New(nil, errors.ErrorCodeUnknown, "bad error"))

	tests.ExecuteE(logger.Errorf("this is my %s", "message")).NoError(t)
	tests.Execute(logs.String()).Diff(t, fmt.Sprintf("[error] %s: (bad error) this is my message\n", time.Now().Format(txt.DateTimeFormat)))
}

func TestText_WithField(t *testing.T) {
	var logs bytes.Buffer

	txt := TextLogger()
	txt.Out = &logs
	txt.DateTimeFormat = "2006-01-02"

	logger := txt.WithField("key", "value")

	tests.ExecuteE(logger.Infof("this is my %s", "message")).NoError(t)
	tests.Execute(logs.String()).Diff(t, fmt.Sprintf("[info] %s: this is my message (value)\n", time.Now().Format(txt.DateTimeFormat)))
}
