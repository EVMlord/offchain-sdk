package baseapp

import (
	"context"
	"os"

	"github.com/berachain/offchain-sdk/client/eth"
	"github.com/berachain/offchain-sdk/job"
	"github.com/berachain/offchain-sdk/log"
	sdk "github.com/berachain/offchain-sdk/types"
)

// BaseApp is the base application.
type BaseApp struct {
	// name is the name of the application
	name string

	// logger is the logger for the baseapp.
	logger log.Logger

	// contains filtered or unexported fields
	ethCfg eth.Config

	// jobMgr
	jobMgr *JobManager

	// ethClient is the client for communicating with the chain
	ethClient eth.Client

	// nonceManager is the manager for managing nonces
	nonceManager eth.NonceManager
}

// New creates a new baseapp.
func New(
	name string,
	logger log.Logger,
	ethCfg *eth.Config,
	jobs []job.Basic,
) *BaseApp {
	return &BaseApp{
		name:   name,
		logger: log.NewBlankLogger(os.Stdout),
		ethCfg: *ethCfg,
		ethClient: eth.NewClient(
			ethCfg,
		),
		jobMgr: NewJobManager(
			name,
			logger,
			jobs,
		),
	}
}

// Logger returns the logger for the baseapp.
func (b *BaseApp) Logger() log.Logger {
	return b.logger.With("namespace", b.name+"-app")
}

// Start starts the baseapp.
func (b *BaseApp) Start() {
	b.Logger().Info("starting app")

	// TODO: create a new context for every job request / creation.
	ctx := sdk.NewContext(
		context.Background(),
		eth.NewContextualClient(
			context.Background(),
			eth.NewClient(&b.ethCfg),
		),
		b.Logger(),
	)
	b.jobMgr.executionPool.Start()
	b.jobMgr.Start(*ctx)

	b.nonceManager = eth.NewNonceManager()
}

// Stop stops the baseapp.
func (b *BaseApp) Stop() {
	b.Logger().Info("stopping app")
	b.jobMgr.executionPool.Stop()
}
