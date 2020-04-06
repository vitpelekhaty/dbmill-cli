package commands

import (
	"testing"

	"github.com/vitpelekhaty/dbmill-cli/internal/pkg/log"
)

var logLevelCases = []struct {
	level     string
	want      log.Level
	withError bool
}{
	{
		level:     "trace",
		want:      log.TraceLevel,
		withError: false,
	},
	{
		level:     "debug",
		want:      log.DebugLevel,
		withError: false,
	},
	{
		level:     "info",
		want:      log.InfoLevel,
		withError: false,
	},
	{
		level:     "warning",
		want:      log.WarningLevel,
		withError: false,
	},
	{
		level:     "error",
		want:      log.ErrorLevel,
		withError: false,
	},
	{
		level:     "fatal",
		want:      log.FatalLevel,
		withError: false,
	},
	{
		level:     "panic",
		want:      log.PanicLevel,
		withError: false,
	},
	{
		level:     "test",
		want:      log.InfoLevel,
		withError: true,
	},
	{
		level:     "nfo",
		want:      log.InfoLevel,
		withError: true,
	},
}

func TestParseLogLevel(t *testing.T) {
	var done bool

	for _, test := range logLevelCases {
		have, err := ParseLogLevel(test.level)
		withError := err != nil

		done = have == test.want && withError == test.withError

		if !done {
			t.Errorf(`ParseLogLevel("%s") failed!`, test.level)
		}
	}
}
