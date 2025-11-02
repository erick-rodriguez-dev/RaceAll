package logger

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type LogLevel int

const (
	LevelDebug LogLevel = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

type Config struct {
	Level LogLevel

	Pretty bool

	TimeFormat string

	Output io.Writer

	EnableCaller bool

	CallerSkipFrameCount int
}

func DefaultConfig() Config {
	return Config{
		Level:                LevelInfo,
		Pretty:               true,
		TimeFormat:           time.RFC3339,
		Output:               os.Stdout,
		EnableCaller:         false,
		CallerSkipFrameCount: 2,
	}
}

func ProductionConfig() Config {
	return Config{
		Level:                LevelInfo,
		Pretty:               false,
		TimeFormat:           "",
		Output:               os.Stdout,
		EnableCaller:         true,
		CallerSkipFrameCount: 2,
	}
}

func DevelopmentConfig() Config {
	return Config{
		Level:                LevelDebug,
		Pretty:               true,
		TimeFormat:           time.RFC3339,
		Output:               os.Stdout,
		EnableCaller:         true,
		CallerSkipFrameCount: 2,
	}
}

var globalLogger *zerolog.Logger

func Init(config Config) {
	var level zerolog.Level
	switch config.Level {
	case LevelDebug:
		level = zerolog.DebugLevel
	case LevelInfo:
		level = zerolog.InfoLevel
	case LevelWarn:
		level = zerolog.WarnLevel
	case LevelError:
		level = zerolog.ErrorLevel
	case LevelFatal:
		level = zerolog.FatalLevel
	default:
		level = zerolog.InfoLevel
	}

	zerolog.SetGlobalLevel(level)

	var output io.Writer = config.Output
	if output == nil {
		output = os.Stdout
	}

	if config.Pretty {
		output = zerolog.ConsoleWriter{
			Out:        output,
			TimeFormat: config.TimeFormat,
			NoColor:    false,
		}
	}

	logger := zerolog.New(output)

	if config.TimeFormat != "" {
		logger = logger.With().Timestamp().Logger()
	} else {
		logger = logger.With().Time("time", time.Now()).Logger()
	}

	if config.EnableCaller {
		logger = logger.With().CallerWithSkipFrameCount(config.CallerSkipFrameCount).Logger()
	}

	globalLogger = &logger
	log.Logger = logger
}

func Get() *zerolog.Logger {
	if globalLogger == nil {
		Init(DefaultConfig())
	}
	return globalLogger
}

func WithModule(module string) *zerolog.Logger {
	logger := Get().With().Str("module", module).Logger()
	return &logger
}

func WithFields(fields map[string]interface{}) *zerolog.Logger {
	logger := Get()
	ctx := logger.With()
	for k, v := range fields {
		ctx = ctx.Interface(k, v)
	}
	l := ctx.Logger()
	return &l
}

func Debug(msg string) {
	Get().Debug().Msg(msg)
}

func Debugf(format string, args ...interface{}) {
	Get().Debug().Msgf(format, args...)
}

func Info(msg string) {
	Get().Info().Msg(msg)
}

func Infof(format string, args ...interface{}) {
	Get().Info().Msgf(format, args...)
}

func Warn(msg string) {
	Get().Warn().Msg(msg)
}

func Warnf(format string, args ...interface{}) {
	Get().Warn().Msgf(format, args...)
}

func ErrorMsg(msg string) {
	Get().Error().Msg(msg)
}

func Errorf(format string, args ...interface{}) {
	Get().Error().Msgf(format, args...)
}

func ErrorWithErr(err error, msg string) {
	Get().Error().Err(err).Msg(msg)
}

func Fatal(msg string) {
	Get().Fatal().Msg(msg)
}

func Fatalf(format string, args ...interface{}) {
	Get().Fatal().Msgf(format, args...)
}

func SetLevel(level LogLevel) {
	var zLevel zerolog.Level
	switch level {
	case LevelDebug:
		zLevel = zerolog.DebugLevel
	case LevelInfo:
		zLevel = zerolog.InfoLevel
	case LevelWarn:
		zLevel = zerolog.WarnLevel
	case LevelError:
		zLevel = zerolog.ErrorLevel
	case LevelFatal:
		zLevel = zerolog.FatalLevel
	default:
		zLevel = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(zLevel)
}
