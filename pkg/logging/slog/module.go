package slog

import (
	"log/slog"

	"github.com/pasataleo/go-config/pkg/inject"
	"github.com/pasataleo/go-logging/pkg/logging"
)

// Module is a self-configuring component that creates an *slog.Logger backed
// by an injected logging.Logger. It should be installed after a logger module
// (json, text, or test) so that logging.Logger is already bound.
type Module struct {
	Logger logging.Logger `inject:""`
}

// NewModule creates a Module that will produce an *slog.Logger when installed.
func NewModule() *Module {
	return &Module{}
}

func (m *Module) Install(s *inject.Source) error {
	s.Bind(slog.New(NewHandler(m.Logger)))
	return nil
}
