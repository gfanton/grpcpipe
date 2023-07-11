package rpcmanager

import (
	"context"
	"strings"

	"github.com/gfanton/grpcutil/lazy"
	grpc "google.golang.org/grpc"
)

type client struct {
	lc         *lazy.Client
	rootCtx    context.Context
	rootCancel context.CancelFunc
}

type ServiceClientRegister interface {
	RegisterService(name string, cc *grpc.ClientConn)
}

func (s *service) RegisterService(serviceName string, cc *grpc.ClientConn) {
	ctx, cancel := context.WithCancel(s.rootCtx)
	s.muCients.Lock()

	if c, ok := s.clients[serviceName]; ok {
		c.rootCancel()
	}

	s.clients[serviceName] = &client{
		rootCtx:    ctx,
		rootCancel: cancel,
		lc:         lazy.NewClient(cc),
	}

	s.muCients.Unlock()
}

func (s *service) getServiceClient(mdesc *MethodDesc) (c *client, ok bool) {
	if mdesc == nil {
		return
	}

	s.muCients.RLock()

	var serviceName string
	if serviceName, ok = getServiceName(mdesc); ok {
		c, ok = s.clients[serviceName]
	}

	s.muCients.RUnlock()

	return
}

func getServiceName(mdesc *MethodDesc) (string, bool) {
	names := strings.SplitN(mdesc.Name, "/", 3)
	if len(names) != 3 {
		return "", false
	}

	return names[1], true
}

type noopClient struct{}

func NewNoopServiceClientRegister() ServiceClientRegister {
	return &noopClient{}
}

func (*noopClient) RegisterService(name string, cc *grpc.ClientConn) {}
