package logger

import (
	"errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"strconv"
	"sync"
	"time"
)

var once sync.Once

var log zerolog.Logger

func LoadLogger() zerolog.Logger {
	once.Do(func() {
		zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
		zerolog.TimeFieldFormat = time.RFC3339Nano

		logLevel, err := strconv.Atoi(os.Getenv("LOG_LEVEL"))
		if err != nil {
			logLevel = int(zerolog.InfoLevel) // default to INFO
		}

		var output io.Writer = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		}

		if os.Getenv("APP_ENV") != "development" {
			fileLogger := &lumberjack.Logger{
				Filename:   "slack.log",
				MaxSize:    5, //
				MaxBackups: 10,
				MaxAge:     14,
				Compress:   true,
			}

			output = zerolog.MultiLevelWriter(os.Stderr, fileLogger)
		}

		log = zerolog.New(output).
			Level(zerolog.Level(logLevel)).
			With().
			Timestamp().
			Logger()
	})

	return log
}

type Logger struct {
	instance zerolog.Logger
}

func GetLogger() *zerolog.Logger {
	return &log
}

func NewLogger(zeroLog zerolog.Logger) *Logger {
	return &Logger{instance: zeroLog}
}

func (l *Logger) Info(msg string) {
	l.instance.Info().Msg(msg)
}

func (l *Logger) Fatal(err error, msg string) {
	l.instance.Fatal().Err(err).Msg(msg)
}

func (l *Logger) Error(err error, msg string) {
	if err != nil {
		l.instance.Error().Err(err).Msg(msg)
	} else {
		l.instance.Error().Err(errors.New("Unknown Error")).Msg(msg)
	}
}
