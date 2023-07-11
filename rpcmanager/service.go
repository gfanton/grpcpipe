package rpcmanager

import (
	"context"
	"sync"

	"github.com/gfanton/grpcutil/lazy"
	"go.uber.org/zap"
)

type Service interface {
	RPCManagerServer
	ServiceClientRegister

	Close() error
}

type Options struct {
	Logger *zap.Logger
}

type service struct {
	rootCtx    context.Context
	rootCancel context.CancelFunc

	logger *zap.Logger

	muCients sync.RWMutex
	clients  map[string]*client

	streams   map[string]*lazy.Stream
	muStreams sync.RWMutex

	UnimplementedRPCManagerServer
}

func (o *Options) applyDefault() {
	if o.Logger == nil {
		o.Logger = zap.NewNop()
	}
}

func NewService(opts *Options) Service {
	opts.applyDefault()
	ctx, cancel := context.WithCancel(context.Background())
	return &service{
		rootCtx:    ctx,
		rootCancel: cancel,
		logger:     opts.Logger,
		clients:    make(map[string]*client),
		streams:    make(map[string]*lazy.Stream),
	}
}

func (s *service) Close() error {
	s.rootCancel()
	return nil
}
