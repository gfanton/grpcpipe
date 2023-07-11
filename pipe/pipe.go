package pipe

import (
	"context"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

type Listener interface {
	net.Listener

	Dial() (net.Conn, error)
	DialContext(ctx context.Context) (net.Conn, error)
}

type Pipe struct {
	Listener
}

func NewNetPipe(sz int) *Pipe {
	return &Pipe{
		Listener: NewNet(),
	}
}

func NewBufferPipe(sz int) *Pipe {
	return &Pipe{Listener: bufconn.Listen(sz)}
}

func (bl *Pipe) dialer(context.Context, string) (net.Conn, error) {
	return bl.Dial()
}

func (bl *Pipe) ClientConn(ctx context.Context, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	mendatoryOpts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(bl.dialer), // set pipe dialer
	}

	return grpc.DialContext(ctx, "buf", append(opts, mendatoryOpts...)...)
}
