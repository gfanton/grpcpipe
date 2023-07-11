// Package pipe provides a net.Conn implemented by net.Pipe and related dialing and listening functionality.
// For a buffered connection pipe use: https://pkg.go.dev/google.golang.org/grpc/test/bufconna

package pipe

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"
)

type Listener struct {
	cancel context.CancelFunc
	ctx    context.Context
	cconn  chan net.Conn
	once   sync.Once
}

func NewListener() *Listener {
	ctx, cancel := context.WithCancel(context.Background())
	return &Listener{
		cancel: cancel,
		ctx:    ctx,
		cconn:  make(chan net.Conn, 1),
	}
}

// Add conn forward the given conn to the listener
func (pl *Listener) AddConn(c net.Conn) {
	select {
	case <-pl.ctx.Done():
	case pl.cconn <- c:
	}
}

func (pl *Listener) Dialer(addr string, _ time.Duration) (net.Conn, error) {
	return pl.ContextDialer(context.Background(), addr)
}

func (pl *Listener) ContextDialer(_ context.Context, _ string) (cclient net.Conn, _ error) {
	var cserver net.Conn
	cclient, cserver = net.Pipe()
	pl.AddConn(cserver)
	return
}

// Listener
var _ net.Listener = (*Listener)(nil)

func (pl *Listener) Addr() net.Addr { return pl }
func (pl *Listener) Accept() (net.Conn, error) {
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
func (pl *Listener) Close() error {
	pl.cancel()
	pl.once.Do(func() { close(pl.cconn) })
	return nil
}

// Addr
var _ net.Addr = (*Listener)(nil)

func (pl *Listener) Network() string { return "pipe_network" }
func (pl *Listener) String() string  { return "pipe" }
