// Package pipe provides a net.Conn implemented by net.Pipe and related dialing and listening functionality.
// For a buffered connection pipe use: https://pkg.go.dev/google.golang.org/grpc/test/bufconna

package pipe

import (
	"context"
	"fmt"
	"net"
	"sync"
)

type Net struct {
	cancel context.CancelFunc
	ctx    context.Context
	cconn  chan net.Conn
	once   sync.Once
}

func NewNet() *Net {
	ctx, cancel := context.WithCancel(context.Background())
	return &Net{
		cancel: cancel,
		ctx:    ctx,
		cconn:  make(chan net.Conn, 1),
	}
}

// Add conn forward the given conn to the listener
func (pl *Net) AddConn(c net.Conn) {
	select {
	case <-pl.ctx.Done():
	case pl.cconn <- c:
	}
}

func (pl *Net) Dial() (net.Conn, error) {
	return pl.DialContext(context.Background())
}

func (pl *Net) DialContext(_ context.Context) (cclient net.Conn, _ error) {
	var cserver net.Conn
	cclient, cserver = net.Pipe()
	pl.AddConn(cserver)
	return
}

// Listener
var _ net.Listener = (*Net)(nil)

func (pl *Net) Addr() net.Addr { return pl }
func (pl *Net) Accept() (net.Conn, error) {
	select {
	case conn := <-pl.cconn:
		if conn != nil {
			return conn, nil
		}
	case <-pl.ctx.Done():
		return nil, pl.ctx.Err()
	}

	return nil, fmt.Errorf("pipe listener is closing")
}
func (pl *Net) Close() error {
	pl.cancel()
	pl.once.Do(func() { close(pl.cconn) })
	return nil
}

// Addr
var _ net.Addr = (*Net)(nil)

func (pl *Net) Network() string { return "pipe_network" }
func (pl *Net) String() string  { return "pipe" }
