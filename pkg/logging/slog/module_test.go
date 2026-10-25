package slog

import (
	"bytes"
	"log/slog"
	"testing"

	"github.com/pasataleo/go-config/pkg/config"
	"github.com/pasataleo/go-config/pkg/inject"
	"github.com/pasataleo/go-testingx/pkg/testingx"

	"github.com/pasataleo/go-logging/pkg/logging"
	"github.com/pasataleo/go-logging/pkg/logging/json"
)

func TestModule(t *testing.T) {
	var out bytes.Buffer
	injector := inject.New()
	injector.Bind(json.New(logging.Info, logging.WithWriter(&out)))
	injector.Register(NewModule())

	var cfg struct {
		Logger *slog.Logger
	}
	testingx.Call(t, config.New(injector).Assign, &cfg).Fatal().NoError()

	cfg.Logger.Info("hello from slog")
	testingx.Call(t, out.String).Contains("hello from slog")
}

func TestModule_noLogger(t *testing.T) {
	injector := inject.New()
	injector.Register(NewModule())

	var cfg struct {
		Logger *slog.Logger
	}
	testingx.Panics(t, nil, func() { _ = config.New(injector).Assign(&cfg) })
}
