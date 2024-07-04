package logger

import (
	"fmt"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type CustomLoggerType struct{}

// deprecated
var CustomLogger *CustomLoggerType

func NewLogger() *CustomLoggerType {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	return &CustomLoggerType{}
}

func Init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
}

func (cl *CustomLoggerType) Debugln(msg ...interface{}) {
	log.Debug().Msg(makeMessage(msg))
}

func (cl *CustomLoggerType) Infoln(msg ...string) {
	log.Info().Msg(makeMessage(msg))
}

func (cl *CustomLoggerType) Errorln(msg ...string) {
	log.Error().Msg(makeMessage(msg))
}

func makeMessage(msg ...interface{}) string {
	fullMessage := ""
	for _, m := range msg {
		fullMessage += fmt.Sprintf("%v ", m)
	}
	return fullMessage
}
