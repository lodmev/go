package log

import (
	"flag"
	stdlog "log"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/rs/zerolog"
)

type LoggerConfig struct {
	Level  string
	Human  bool
	Caller bool
}

func ExportConfig(fs *flag.FlagSet) *LoggerConfig {
	conf := LoggerConfig{}
	fs.StringVar(&conf.Level, "lvl", "info", "level of logging")
	fs.BoolVar(&conf.Human, "hum", false, "set pretty output of logging")
	fs.BoolVar(&conf.Caller, "cal", false, "add info about code, that call log message")
	return &conf

}

func (conf LoggerConfig) Setup() {
	if conf.Human {
		Logger = Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.Stamp})
	}
	level, err := parseLevel(conf.Level)
	if err != nil {
		Fatal().
			Err(err).
			Msgf("Can't parse log level %s\n", conf.Level)
	}
	zerolog.SetGlobalLevel(level)
	if conf.Caller {
		Logger = With().Caller().Logger()
	}
	stdlog.SetFlags(0)
	stdlog.SetOutput(Logger)

}

func parseLevel(l string) (zerolog.Level, error) {
	switch strings.ToLower(l) {
	case "trace":
		return zerolog.TraceLevel, nil
	case "debug":
		return zerolog.DebugLevel, nil
	case "info":
		return zerolog.InfoLevel, nil
	case "warn":
		return zerolog.WarnLevel, nil
	case "error":
		return zerolog.ErrorLevel, nil
	}

	return zerolog.InfoLevel, errors.Errorf("invalid level: %v", l)
}
