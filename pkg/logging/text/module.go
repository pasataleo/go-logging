package text

import (
	"github.com/pasataleo/go-config/pkg/inject"
	"github.com/pasataleo/go-logging/pkg/logging"
)

// Module is a self-configuring component that creates a text logger. The
// framework populates Level from flags/env/defaults, then Install creates the
// logger and binds it for injection.
type Module struct {
	Level   logging.Level `flag:"log-level" env:"LOG_LEVEL" default:"info" hidden:"true"`
	options []logging.Opt
}

// NewModule creates a Module that will produce a text logger with the given
// options when installed.
func NewModule(options ...logging.Opt) *Module {
	return &Module{options: options}
}

func (m *Module) Install(s *inject.Source) error {
	s.Bind(New(m.Level, m.options...))
	return nil
}
