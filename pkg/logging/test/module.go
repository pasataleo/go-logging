package test

import (
	"testing"

	"github.com/pasataleo/go-config/pkg/inject"
	"github.com/pasataleo/go-logging/pkg/logging"
)

// Module is a self-configuring component that creates a test logger. The
// framework populates Level from flags/env/defaults, then Install creates the
// logger and binds it for injection.
type Module struct {
	Level logging.Level `flag:"log-level" env:"LOG_LEVEL" default:"info" hidden:"true"`
	tb    testing.TB
}

// NewModule creates a Module that will produce a test logger forwarding to tb
// when installed.
func NewModule(tb testing.TB) *Module {
	return &Module{tb: tb}
}

func (m *Module) Install(s *inject.Source) error {
	s.Bind(New(m.tb, m.Level))
	return nil
}
