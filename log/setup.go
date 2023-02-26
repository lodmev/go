package log

import (
	"flag"
	"fmt"
	"io"
	stdlog "log"
	"os"
	"time"
"encoding/json"

	"github.com/rs/zerolog"
)

var Logger = zerolog.New(os.Stderr).With().Timestamp().Logger()
type LogLevel struct {
	zerolog.Level
}
func (l *LogLevel) Set(s string) error {
	level, err:= zerolog.ParseLevel(s)
	l.Level = level
	return err
}

type Config struct {
	Lvl   LogLevel
	Human   bool
	Caller  bool
	Outfile string
}

func ExportConf(fs *flag.FlagSet) *Config {
	conf := &Config{}
	fs.Var(&conf.Lvl, "level","level of logging, may be: trace, debug, info, warn, error")
	fs.BoolVar(&conf.Human, "human", false, "set pretty output of logging")
	fs.BoolVar(&conf.Caller, "caller", false, "add info about code, that call log message")
	fs.StringVar(&conf.Outfile, "outfile", "", "path to log file")
	return conf

}

func (conf *Config) Setup() (err error) {
	loc, locerr := time.LoadLocation("Europe/Moscow")
	if locerr != nil {
		Errf("can't load TZ locaton: %w", err)
	} else {
		zerolog.TimestampFunc = func () time.Time {
			return time.Now().In(loc)
		}
	}
	Logger = Level(conf.Lvl.Level)
	var loggerOutput io.Writer = os.Stderr
	if conf.Human {
		loggerOutput = getConsoleWriter(loggerOutput)
	}
	if conf.Caller {
		Logger = With().Caller().Logger()
	}
	if conf.Outfile != "" {
		logFile, err := os.OpenFile(conf.Outfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0664)
		if err != nil {
			Logger.Err(err).Msgf("unable to open Log file %q", conf.Outfile)
		} else {
		loggerOutput = zerolog.MultiLevelWriter(loggerOutput, logFile)
	}
}
	Logger = zerolog.New(loggerOutput).With().Timestamp().Logger().Level(conf.Lvl.Level)
	stdlog.SetFlags(0)
	stdlog.SetOutput(Logger)
	return
}

func getConsoleWriter (out io.Writer) zerolog.ConsoleWriter{
// colorize returns the string s wrapped in ANSI code c, unless disabled is true.
colorize :=func (s interface{}, c int, out io.Writer) string {
	if out != os.Stderr && out != os.Stdout {
		return fmt.Sprintf("%s", s)
	}
	return fmt.Sprintf("\x1b[%dm%v\x1b[0m", c, s)
}
	return zerolog.ConsoleWriter{
	Out: out,
	FormatTimestamp: func(i interface{}) string {
t := "<nil>"
		switch tt := i.(type) {
		case string:
			ts, err := time.Parse(time.RFC3339, tt)
			if err != nil {
				t = tt
			} else {
				t = fmt.Sprintf("%02d:%02d:%02d",ts.Hour(),ts.Minute(), ts.Second())			}
		case json.Number:
			i, err := tt.Int64()
			if err != nil {
				t = tt.String()
			} else {
				var sec, nsec int64

				switch zerolog.TimeFieldFormat {
				case zerolog.TimeFormatUnixNano:
					sec, nsec = 0, i
				case zerolog.TimeFormatUnixMicro:
					sec, nsec = 0, int64(time.Duration(i)*time.Microsecond)
				case zerolog.TimeFormatUnixMs:
					sec, nsec = 0, int64(time.Duration(i)*time.Millisecond)
				default:
					sec, nsec = i, 0
				}

				ts := time.Unix(sec, nsec)
				t = ts.String()
			}
		}
		return colorize(t, 90, out)
	},

}
}


