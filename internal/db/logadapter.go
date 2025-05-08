package db

import (
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
)

type LogAdapter struct{}

func (l LogAdapter) Fatalf(format string, v ...interface{}) {
	log.Fatal().Msg(strings.TrimSpace(fmt.Sprintf(format, v...)))
}

func (l LogAdapter) Printf(format string, v ...interface{}) {
	log.Info().Msg(strings.TrimSpace(fmt.Sprintf(format, v...)))
}
