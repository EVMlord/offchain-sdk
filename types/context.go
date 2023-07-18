package types

import (
	"context"

	"cosmossdk.io/log"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
)

type Chain interface {
	ChainReader
	ChainWriter
	ChainSubscriber
}

type ChainWriter interface{}

type ChainReader interface {
	CurrentBlock() (*types.Block, error)
	GetBlockByNumber(number uint64) (*types.Block, error)
}

type ChainSubscriber interface {
	SubscribeFilterLogs(q ethereum.FilterQuery, ch chan<- types.Log) (ethereum.Subscription, error)
}

type Context struct {
	context.Context
	chain  Chain
	logger log.Logger
	// chain Chain
}

// UnwrapSdkContext unwraps the sdk context.
func UnwrapSdkContext(ctx context.Context) Context {
	if sdkCtx, ok := ctx.(Context); ok {
		return sdkCtx
	}
	panic("context is not sdk context")
}

func NewContext(ctx context.Context, chain Chain, logger log.Logger) *Context {
	return &Context{
		Context: ctx,
		chain:   chain,
		logger:  logger,
	}
}

func (c *Context) Chain() Chain {
	return c.chain
}

func (c *Context) Logger() log.Logger {
	return c.logger
}
