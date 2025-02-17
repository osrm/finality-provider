package e2etest

import (
	"testing"

	"github.com/lightningnetwork/lnd/signal"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/babylonlabs-io/finality-provider/eotsmanager"
	"github.com/babylonlabs-io/finality-provider/eotsmanager/config"
	"github.com/babylonlabs-io/finality-provider/eotsmanager/service"
)

type EOTSServerHandler struct {
	t           *testing.T
	interceptor *signal.Interceptor
	eotsServer  *service.Server
}

func NewEOTSServerHandler(t *testing.T, cfg *config.Config, eotsHomeDir string, shutdownInterceptor signal.Interceptor) *EOTSServerHandler {
	dbBackend, err := cfg.DatabaseConfig.GetDBBackend()
	require.NoError(t, err)
	logger := zap.NewNop()
	eotsManager, err := eotsmanager.NewLocalEOTSManager(eotsHomeDir, cfg.KeyringBackend, dbBackend, logger)
	require.NoError(t, err)

	eotsServer := service.NewEOTSManagerServer(cfg, logger, eotsManager, dbBackend, shutdownInterceptor)

	return &EOTSServerHandler{
		t:           t,
		eotsServer:  eotsServer,
		interceptor: &shutdownInterceptor,
	}
}

func (eh *EOTSServerHandler) Start() {
	go eh.startServer()
}

func (eh *EOTSServerHandler) startServer() {
	err := eh.eotsServer.RunUntilShutdown()
	require.NoError(eh.t, err)
}

func (eh *EOTSServerHandler) Stop() {
	eh.interceptor.RequestShutdown()
}
