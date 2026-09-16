package server

import (
	"context"
	"net"

	"github.com/mrdude/pcc-common"
	"go.uber.org/zap"
)

func attachServerNameToLogger(logger *zap.Logger, sv Server) *zap.Logger {
	return logger.With(zap.String("server.name", sv.Name()))
}

type DeathAction string

const (
	CrashProcessOnServerDeath DeathAction = "crash"
	LogErrorOnServerDeath     DeathAction = "log"
)

// Start starts the server in a new goroutine
func Start(ctx context.Context, sock net.Listener, sv Server, deathAction DeathAction) {
	logger := pcommon.GetLogger(ctx)
	logger = attachServerNameToLogger(logger, sv)

	go func() {
		logger.Info("Server started")
		err := sv.Run(ctx, sock)

		switch deathAction {
		case CrashProcessOnServerDeath:
			logger.Fatal("server died",
				zap.Error(err))
		case LogErrorOnServerDeath:
			logger.Error("server died",
				zap.Error(err))
		}
	}()
}
