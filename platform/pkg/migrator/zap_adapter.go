package migrator

import (
	"context"
	"fmt"

	"github.com/pressly/goose/v3"
	"go.uber.org/zap"

	"github.com/aleksandr-mv/school_schedule/platform/pkg/logger"
)

type gooseZapAdapter struct{}

// Print вызывается для простых сообщений.
func (a *gooseZapAdapter) Print(v ...interface{}) {
	msg := fmt.Sprint(v...)
	logger.Info(
		context.Background(),
		"🐣 [goose] "+msg,
		zap.String("component", "goose"),
	)
}

// Printf вызывается для форматированных сообщений.
func (a *gooseZapAdapter) Printf(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	logger.Info(
		context.Background(),
		"🐣 [goose] "+msg,
		zap.String("component", "goose"),
	)
}

// Fatalf вызывается при фатальных ошибках goose.
func (a *gooseZapAdapter) Fatalf(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	logger.Error(
		context.Background(),
		"💥 [goose] "+msg,
		zap.String("component", "goose"),
	)
	panic(msg)
}

func init() {
	goose.SetLogger(&gooseZapAdapter{})
}
