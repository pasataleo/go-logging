package logging

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/pasataleo/go-errors/errors"
)

var _ Logger = (*Json)(nil)

type Json struct {
	err error // The last error requested to be tracked by this logger.

	// Level is the log level this Text logger will include in its output.
	//
	// This defaults to Info.
	Level Level

	// Out is the stream that this Text logger will write to.
	//
	// This defaults to os.Stderr.
	Out io.Writer

	// DateTimeFormat is the format any log messages should write the date and
	// time of the format.
	//
	// This should be provided in the same format that is supported by the
	// Time.Format function: https://pkg.go.dev/time#Time.Format.
	//
	// This defaults to time.RFC3339Nano.
	DateTimeFormat string

	// AdditionalFields contains any extra fields that should be attached to the
	// written log entries.
	AdditionalFields map[string]interface{}
}

func JsonLogger() *Json {
	return &Json{
		Level:            Info,
		Out:              os.Stderr,
		DateTimeFormat:   time.RFC3339Nano,
		AdditionalFields: map[string]interface{}{},
	}
}

func (j *Json) Copy() *Json {
	fields := map[string]interface{}{}
	for key, value := range j.AdditionalFields {
		fields[key] = value
	}

	return &Json{
		err:              j.err,
		Level:            j.Level,
		Out:              j.Out,
		DateTimeFormat:   j.DateTimeFormat,
		AdditionalFields: fields,
	}
}

func (j *Json) Print(level Level, message string) error {
	if !j.Level.includes(level) {
		// Only print if our level includes the target level.
		return nil
	}

	entry := map[string]interface{}{}
	now := time.Now().UTC()
	entry["datetime"] = now.Format(j.DateTimeFormat)
	entry["level"] = level.String()
	if j.err != nil {
		entry["error"] = j.err.Error()

		// Automatically include any embedded data in the error.
		data := errors.GetAllEmbeddedData(j.err)
		for key, value := range data {
			entry[fmt.Sprintf("error.%s", key)] = value
		}

		// Automatically include the error code if it's not ErrorCodeUnknown.
		errorCode := errors.GetErrorCode(j.err)
		if errorCode != errors.ErrorCodeUnknown {
			entry["error.code"] = errorCode
		}
	}
	entry["message"] = message

	for key, field := range j.AdditionalFields {
		entry[key] = field
	}

	jEntry, err := json.Marshal(entry)
	if err != nil {
		return errors.Wrap(err, "could not marshal json entry")
	}

	_, err = j.Out.Write([]byte(fmt.Sprintf("%s\n", jEntry)))
	if err != nil {
		return errors.Wrap(err, "failed to write log message")
	}
	return nil
}

func (j *Json) Printf(level Level, format string, args ...any) error {
	return j.Print(level, fmt.Sprintf(format, args...))
}

func (j *Json) Error(message string) error {
	return j.Print(Error, message)
}

func (j *Json) Errorf(format string, args ...any) error {
	return j.Printf(Error, format, args...)
}

func (j *Json) Warn(message string) error {
	return j.Print(Warn, message)
}

func (j *Json) Warnf(format string, args ...any) error {
	return j.Printf(Warn, format, args...)
}

func (j *Json) Info(message string) error {
	return j.Print(Info, message)
}

func (j *Json) Infof(format string, args ...any) error {
	return j.Printf(Info, format, args...)
}

func (j *Json) Trace(message string) error {
	return j.Print(Trace, message)
}

func (j *Json) Tracef(format string, args ...any) error {
	return j.Printf(Trace, format, args...)
}

func (j *Json) WithError(err error) Logger {
	newj := j.Copy()
	newj.err = err
	return newj
}

func (j *Json) WithField(key string, value interface{}) Logger {
	newj := j.Copy()
	newj.AdditionalFields[key] = value
	return newj
}

func (j *Json) WithFields(fields map[string]interface{}) Logger {
	newj := j.Copy()
	for key, value := range fields {
		newj.AdditionalFields[key] = value
	}
	return newj
}
