package pipe

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"
)

type PipeListener struct {
	cancel context.CancelFunc
	ctx    context.Context
	cconn  chan net.Conn
	once   sync.Once
}

func NewPipeListener() *PipeListener {
	ctx, cancel := context.WithCancel(context.Background())
	return &PipeListener{
		cancel: cancel,
		ctx:    ctx,
		cconn:  make(chan net.Conn, 1),
	}
}

// Add conn forward the given conn to the listener
func (pl *PipeListener) AddConn(c net.Conn) {
	select {
	case <-pl.ctx.Done():
	case pl.cconn <- c:
	}
}

func (pl *PipeListener) Dialer(addr string, _ time.Duration) (net.Conn, error) {
	return pl.ContextDialer(context.Background(), addr)
}

func (pl *PipeListener) ContextDialer(ctx context.Context, addr string) (cclient net.Conn, _ error) {
	var cserver net.Conn
	cclient, cserver = net.Pipe()
	pl.AddConn(cserver)
	return
}

// Listener
var _ net.Listener = (*PipeListener)(nil)

func (pl *PipeListener) Addr() net.Addr { return pl }
func (pl *PipeListener) Accept() (net.Conn, error) {
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
func (pl *PipeListener) Close() error {
	pl.cancel()
	pl.once.Do(func() { close(pl.cconn) })
	return nil
}

// Addr
var _ net.Addr = (*PipeListener)(nil)

func (pl *PipeListener) Network() string { return "pipe_network" }
func (pl *PipeListener) String() string  { return "pipe" }
