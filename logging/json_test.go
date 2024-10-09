package logging

import (
	"bytes"
	"encoding/json"
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

	tests.ExecuteE(logger.Infof("this is my %s", "message")).NoError(t)
	tests.Execute(logs.String()).Diff(t, mustJson(t, map[string]interface{}{
		"datetime": time.Now().Format(logger.DateTimeFormat),
		"level":    "info",
		"message":  "this is my message",
	}))
}

func TestJson_WithError(t *testing.T) {
	var logs bytes.Buffer

	json := JsonLogger()
	json.Out = &logs
	json.DateTimeFormat = "2006-01-02"

	logger := json.WithError(errors.New(nil, errors.ErrorCodeUnknown, "bad error"))

	tests.ExecuteE(logger.Errorf("this is my %s", "message")).NoError(t)
	tests.Execute(logs.String()).Diff(t, mustJson(t, map[string]interface{}{
		"datetime": time.Now().Format(json.DateTimeFormat),
		"error":    "bad error",
		"level":    "error",
		"message":  "this is my message",
	}))
}

func TestJson_WithField(t *testing.T) {
	var logs bytes.Buffer

	json := JsonLogger()
	json.Out = &logs
	json.DateTimeFormat = "2006-01-02"

	logger := json.WithField("key", "value")

	tests.ExecuteE(logger.Infof("this is my %s", "message")).NoError(t)
	tests.Execute(logs.String()).Diff(t, mustJson(t, map[string]interface{}{
		"datetime": time.Now().Format(json.DateTimeFormat),
		"key":      "value",
		"level":    "info",
		"message":  "this is my message",
	}))
}

func mustJson(t *testing.T, data interface{}) string {
	t.Helper()
	return string(tests.Execute2E(json.Marshal(data)).NoError(t).Capture()) + "\n"
}
