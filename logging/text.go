package logging

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/pasataleo/go-errors/errors"
)

var _ Logger = (*Text)(nil)

// Text is a simple Logger that writes messages in a given format to a specified
// output. It is designed to write out human-readable logs.
type Text struct {
	err error // The last error requested to be tracked by this logger.

	// Level is the log level this Text logger will include in its output.
	//
	// This defaults to Info.
	Level Level

	// Out is the stream that this Text logger will write to.
	//
	// This defaults to os.Stderr.
	Out io.Writer

	// Format is the format this Text logger will write log messages.
	//
	// This defaults to "[%L] %T: %e%m%f\n".
	//
	// It can be overridden to customise the format of the log messages. The
	// following elements are supported:
	//   - %m: The actual log message.
	//   - %T: The date time of the log message, this can be further customised
	//         by editing the DateTimeFormat field.
	//   - %L: The current log level.
	//   - %e: The value for any attached error, this can be further customised.
	//   - %f: The value for any additional fields, this can be further
	//         customised by editing the AdditionalFieldsFormat field.
	Format string

	// DateTimeFormat is the format any log messages should write the date and
	// time of the format.
	//
	// This should be provided in the same format that is supported by the
	// Time.Format function: https://pkg.go.dev/time#Time.Format.
	//
	// This defaults to time.RFC3339Nano.
	DateTimeFormat string

	// ErrorFormat describes how any attached error should be rendered.
	//
	// This defaults to "(%s) ", the %s is where the error will be rendered.
	ErrorFormat string

	// AdditionalFieldsFormat describes how any additional fields should be
	// rendered.
	//
	// This defaults to " [%s]", where %s is the fields joined into a writable
	// list format using the ',' character.
	AdditionalFieldsFormat string

	// AdditionalFields contains any extra fields that should be attached to the
	// written log entries.
	//
	// The key for a field can be specified in the Format preceded by a '%'
	// character in order to customise exactly where the additional fields will
	// be rendered. You can also override existing keys using the
	// AdditionalFields parameter.
	//
	// Or, you can specify %f to render all the fields using the
	// AdditionalFieldsFormat parameter.
	AdditionalFields map[string]interface{}
}

func TextLogger() *Text {
	return &Text{
		Level:                  Info,
		Out:                    os.Stderr,
		Format:                 "[%L] %T: %e%m%f\n",
		DateTimeFormat:         time.RFC3339Nano,
		ErrorFormat:            "(%s) ",
		AdditionalFieldsFormat: " (%s)",
		AdditionalFields:       map[string]interface{}{},
	}
}

// Copy creates a copy of the current Text logger.
func (text *Text) Copy() *Text {

	fields := map[string]interface{}{}
	for key, value := range text.AdditionalFields {
		fields[key] = value
	}

	return &Text{
		err:                    text.err,
		Level:                  text.Level,
		Out:                    text.Out,
		Format:                 text.Format,
		DateTimeFormat:         text.DateTimeFormat,
		ErrorFormat:            text.ErrorFormat,
		AdditionalFieldsFormat: text.AdditionalFieldsFormat,
		AdditionalFields:       fields,
	}
}

func (text *Text) Print(level Level, message string) error {
	if !text.Level.includes(level) {
		// Only print if our level includes the target level.
		return nil
	}

	now := time.Now().UTC()
	processed := text.Format
	for key, value := range text.AdditionalFields {
		processed = strings.ReplaceAll(processed, fmt.Sprintf("%%%s", key), fmt.Sprintf("%v", value))
	}

	processed = strings.ReplaceAll(processed, "%T", now.Format(text.DateTimeFormat))
	processed = strings.ReplaceAll(processed, "%L", level.String())
	if text.err != nil {
		processed = strings.ReplaceAll(processed, "%e", fmt.Sprintf(text.ErrorFormat, text.err.Error()))
	} else {
		processed = strings.ReplaceAll(processed, "%e", "")
	}
	if len(text.AdditionalFields) > 0 {
		processed = strings.ReplaceAll(processed, "%f", fmt.Sprintf(text.AdditionalFieldsFormat, strings.Join(func() []string {
			var fields []string
			for _, field := range text.AdditionalFields {
				fields = append(fields, fmt.Sprintf("%v", field))
			}
			return fields
		}(), ",")))
	} else {
		processed = strings.ReplaceAll(processed, "%f", "")
	}
	processed = strings.ReplaceAll(processed, "%m", message)
	_, err := text.Out.Write([]byte(processed))
	if err != nil {
		return errors.Wrap(err, "failed to write log message")
	}
	return nil
}

func (text *Text) Printf(level Level, format string, args ...any) error {
	return text.Print(level, fmt.Sprintf(format, args...))
}

func (text *Text) Error(message string) error {
	return text.Print(Error, message)
}

func (text *Text) Errorf(format string, args ...any) error {
	return text.Printf(Error, format, args...)
}

func (text *Text) Warn(message string) error {
	return text.Print(Warn, message)
}

func (text *Text) Warnf(format string, args ...any) error {
	return text.Printf(Warn, format, args...)
}

func (text *Text) Info(message string) error {
	return text.Print(Info, message)
}

func (text *Text) Infof(format string, args ...any) error {
	return text.Printf(Info, format, args...)
}

func (text *Text) Trace(message string) error {
	return text.Print(Trace, message)
}

func (text *Text) Tracef(format string, args ...any) error {
	return text.Printf(Trace, format, args...)
}

func (text *Text) WithError(err error) Logger {
	newt := text.Copy()
	newt.err = err
	return newt
}

func (text *Text) WithField(key string, value interface{}) Logger {
	newt := text.Copy()
	newt.AdditionalFields[key] = value
	return newt
}

func (text *Text) WithFields(fields map[string]interface{}) Logger {
	newt := text.Copy()
	for key, value := range fields {
		newt.AdditionalFields[key] = value
	}
	return newt
}
