package lazy

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/interop"
	testgrpc "google.golang.org/grpc/interop/grpc_testing"
	testpb "google.golang.org/grpc/interop/grpc_testing"
	"google.golang.org/grpc/test/bufconn"
)

func testingService(t *testing.T) *grpc.ClientConn {
	srv := grpc.NewServer()
	t.Cleanup(srv.Stop)

	l := bufconn.Listen(4096)
	service := interop.NewTestServer()
	testgrpc.RegisterTestServiceServer(srv, service)

	go srv.Serve(l)

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	cc, err := grpc.DialContext(ctx, "", grpc.WithDialer(
		func(_ string, _ time.Duration) (net.Conn, error) {
			return l.DialContext(ctx)
		}), grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)

	return cc
}

func TestNewLazyClient(t *testing.T) {
	cc := testingService(t)

	lc := NewClient(cc)
	assert.NotNil(t, lc)
	assert.Equal(t, cc, lc.cc)
}

func TestInvokeUnary(t *testing.T) {
	cc := testingService(t)
	defer cc.Close()

	lc := NewClient(cc)

	desc := &MethodDesc{
		Name: testpb.TestService_UnaryCall_FullMethodName,
	}

	msg := new(testpb.SimpleRequest)

	in, err := NewMessage().FromMessage(msg)
	assert.NoError(t, err)

	out, err := lc.InvokeUnary(context.Background(), desc, in)

	assert.NotNil(t, out)
	assert.Nil(t, err)
}

func TestInvokeStream(t *testing.T) {
	const testCount = 10

	cc := testingService(t)
	defer cc.Close()

	lc := NewClient(cc)

	desc := &MethodDesc{
		Name:          testpb.TestService_FullDuplexCall_FullMethodName,
		ServerStreams: true,
		ClientStreams: true,
	}

	ls, err := lc.InvokeStream(context.Background(), desc, nil)

	require.NoError(t, err)

	cout := make(chan *Message, testCount)
	go func() {
		defer close(cout)
		for i := 0; i < testCount; i++ {
			out := NewMessage()
			err := ls.RecvMsg(out)
			assert.NoError(t, err)
			cout <- out
		}
	}()

	for i := 0; i < testCount; i++ {
		hello := fmt.Sprintf("hello[%d]", i)
		msg := &testpb.StreamingOutputCallRequest{
			ResponseParameters: []*testpb.ResponseParameters{
				{Size: 8000},
			},
			Payload: &testpb.Payload{
				Body: []byte(hello),
			},
		}

		in, err := NewMessage().FromMessage(msg)
		require.NoError(t, err)

		err = ls.SendMsg(in)
		assert.NoError(t, err)
	}

	for out := range cout {
		require.NotNil(t, out)
		assert.Greater(t, len(out.Bytes()), 0)
	}

}
