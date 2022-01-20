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

var (
	ServiceName    = filepath.Base(os.Args[0])
	ServiceType    = ""
	ServiceEnv     = ""
	ServiceVersion = ""
)

func init() {
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
		Str("service.name", ServiceName).
		Str("process.executable", os.Args[0]).
		Int("process.pid", os.Getpid()).
		Int("process.ppid", os.Getppid()).
		Str("process.start", zerolog.TimestampFunc().Format(time.RFC3339Nano)).
		Str("os.platform", runtime.GOOS)

	if hostname, err := os.Hostname(); err == nil {
		zlc = zlc.Str("host.hostname", hostname)
	}

	if ServiceType != "" {
		zlc = zlc.Str("service.type", ServiceType)
	}

	if ServiceEnv != "" {
		zlc = zlc.Str("service.environment", ServiceEnv)
	}

	if ServiceVersion != "" {
		zlc = zlc.Str("service.version", ServiceVersion)
	}

	if len(os.Args) > 1 {
		zlc = zlc.Strs("process.args", os.Args[1:])
	}

	if wd, err := os.Getwd(); err == nil {
		zlc = zlc.Str("process.working_directory", wd)
	}

	log.Logger = zlc.Logger()
}
