package logger

import (
	"fmt"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type customLogger struct{}

var CustomLogger *customLogger

func Init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
}

func (cl *customLogger) Debugln(msg ...interface{}) {
	log.Debug().Msg(makeMessage(msg))
}

func (cl *customLogger) Infoln(msg ...string) {
	log.Info().Msg(makeMessage(msg))
}

func (cl *customLogger) Errorln(msg ...string) {
	log.Error().Msg(makeMessage(msg))
}

func makeMessage(msg ...interface{}) string {
	fullMessage := ""
	for _, m := range msg {
		fullMessage += fmt.Sprintf("%v ", m)
	}
	return fullMessage
}
