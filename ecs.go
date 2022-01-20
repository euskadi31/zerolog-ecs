// Copyright 2022 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package zerologecs

import (
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const Version = "1.12"

type Option func(*config)

type config struct {
	serviceName    string
	serviceType    string
	serviceEnv     string
	serviceVersion string
}

func WithServiceName(name string) Option {
	return func(c *config) {
		c.serviceName = name
	}
}

func WithServiceEnv(env string) Option {
	return func(c *config) {
		c.serviceEnv = env
	}
}

func WithServiceType(t string) Option {
	return func(c *config) {
		c.serviceType = t
	}
}

func WithServiceVersion(version string) Option {
	return func(c *config) {
		c.serviceVersion = version
	}
}

func Configure(opts ...Option) {
	cfg := &config{
		serviceName: filepath.Base(os.Args[0]),
	}

	for _, opt := range opts {
		opt(cfg)
	}

	zerolog.TimeFieldFormat = time.RFC3339Nano
	zerolog.TimestampFieldName = "@timestamp"
	zerolog.LevelFieldName = "log.level"
	zerolog.MessageFieldName = "message"
	zerolog.ErrorFieldName = "error.message"
	zerolog.CallerFieldName = "log.origin.function"
	zerolog.ErrorStackFieldName = "error.stack_trace"

	zerolog.TimestampFunc = func() time.Time {
		return time.Now().UTC()
	}

	zlc := log.With().
		Str("ecs.version", Version).
		Str("service.name", cfg.serviceName).
		Str("process.executable", os.Args[0]).
		Int("process.pid", os.Getpid()).
		Int("process.ppid", os.Getppid()).
		Str("process.start", zerolog.TimestampFunc().Format(time.RFC3339Nano)).
		Str("os.platform", runtime.GOOS)

	if hostname, err := os.Hostname(); err == nil {
		zlc = zlc.Str("host.hostname", hostname)
	}

	if cfg.serviceType != "" {
		zlc = zlc.Str("service.type", cfg.serviceType)
	}

	if cfg.serviceEnv != "" {
		zlc = zlc.Str("service.environment", cfg.serviceEnv)
	}

	if cfg.serviceVersion != "" {
		zlc = zlc.Str("service.version", cfg.serviceVersion)
	}

	if len(os.Args) > 1 {
		zlc = zlc.Strs("process.args", os.Args[1:])
	}

	if wd, err := os.Getwd(); err == nil {
		zlc = zlc.Str("process.working_directory", wd)
	}

	log.Logger = zlc.Logger()
}
