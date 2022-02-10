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

type Config struct {
	Level  string
	Human  bool
	Caller bool
}

func ExportConf(fs *flag.FlagSet) *Config {
	conf := Config{}
	fs.StringVar(&conf.Level, "lvl", "inf", "level of logging, may be: trc, dbg, inf, warn, err")
	fs.BoolVar(&conf.Human, "hum", false, "set pretty output of logging")
	fs.BoolVar(&conf.Caller, "cal", false, "add info about code, that call log message")
	return &conf

}

func (conf Config) Setup() (err error) {
	if conf.Human {
		Logger = Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.Stamp})
	}
	level, err := parseLevel(conf.Level)
	zerolog.SetGlobalLevel(level)
	if conf.Caller {
		Logger = With().Caller().Logger()
	}
	stdlog.SetFlags(0)
	stdlog.SetOutput(Logger)
	return
}

func parseLevel(l string) (zerolog.Level, error) {
	switch strings.ToLower(l) {
	case "trc":
		return zerolog.TraceLevel, nil
	case "dbg":
		return zerolog.DebugLevel, nil
	case "inf":
		return zerolog.InfoLevel, nil
	case "warn":
		return zerolog.WarnLevel, nil
	case "err":
		return zerolog.ErrorLevel, nil
	}

	return zerolog.InfoLevel, errors.Errorf("invalid level: %v", l)
}
