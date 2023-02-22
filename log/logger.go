package log

import (
	"flag"
	"io"
	stdlog "log"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/rs/zerolog"
)

type Config struct {
	Level   string
	Human   bool
	Caller  bool
	Outfile string
}

func ExportConf(fs *flag.FlagSet) *Config {
	conf := Config{}
	fs.StringVar(&conf.Level, "lvl", "inf", "level of logging, may be: trc, dbg, inf, warn, err")
	fs.BoolVar(&conf.Human, "hum", false, "set pretty output of logging")
	fs.BoolVar(&conf.Caller, "cal", false, "add info about code, that call log message")
	fs.StringVar(&conf.Outfile, "outfile", "", "path to log file")
	return &conf

}

func (conf Config) Setup() (err error) {
	var loggerOutput io.Writer = os.Stdout
	if loc, err := time.LoadLocation("Europe/Moscow"); err == nil {
		zerolog.TimestampFunc = func() time.Time {
			return time.Now().In(loc)
		}
	}
	if conf.Human {
		loggerOutput = zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.Kitchen}
	}

	level, err := parseLevel(conf.Level)
		if err != nil {
			Logger.Err(err).Msgf("unable to parse level %q", conf.Level)
		}
	zerolog.SetGlobalLevel(level)
	if conf.Caller {
		Logger = With().Caller().Logger()
	}
	if conf.Outfile != "" {
		logFile, err := os.OpenFile(conf.Outfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0664)
		if err != nil {
			Logger.Err(err).Msgf("unable to open Log file %q", conf.Outfile)
		}
		loggerOutput = zerolog.MultiLevelWriter(loggerOutput, logFile)
	}
	Logger = Output(loggerOutput)
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
