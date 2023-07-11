package rpcmanager

import (
	context "context"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/interop"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/proto"

	testgrpc "google.golang.org/grpc/interop/grpc_testing"
	testpb "google.golang.org/grpc/interop/grpc_testing"
)

const messageStringTest = "Im sorry Dave, Im afraid I cant do that"
const messageTestSize = 8032

func TestNewService(t *testing.T) {
	s := NewService(&Options{})

	err := s.Close()
	assert.NoError(t, err)
}

func TestUnaryService(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cl := createBridgeTestingClient(t, ctx)

	// call `testutil.TestService/EchoTest` with empty request,
	{
		in := &testpb.SimpleRequest{
			ResponseSize: messageTestSize,
			ResponseType: testpb.PayloadType_COMPRESSABLE,
		}
		payload, err := proto.Marshal(in)
		require.NoError(t, err)

		res, err := cl.ClientInvokeUnary(ctx, &ClientInvokeUnary_Request{
			MethodDesc: &MethodDesc{
				Name: testpb.TestService_UnaryCall_FullMethodName,
			},
			Payload: payload,
		})

		require.NoError(t, err)
		require.NotNil(t, res.Error)
		require.NotNil(t, res.Payload)
		assert.Equal(t, GRPCErrCode_OK, res.Error.GrpcErrorCode)

		out := new(testpb.SimpleResponse)
		err = proto.Unmarshal(res.Payload, out)
		require.NoError(t, err)
		require.NotNil(t, res.Error)

		assert.Equal(t, messageTestSize, len(out.Payload.Body))
		// assert.Equal(t, errcode.Undefined, res.Error.ErrorCode)
		// assert.Nil(t, res.Error.ErrorDetails)
	}
}

func TestUnaryUnimplementedServiceError(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cl := createBridgeTestingClient(t, ctx)

	// call `testutil.TestService/EchoTest` with empty request,
	{
		in := &testpb.SimpleRequest{}

		payload, err := proto.Marshal(in)
		require.NoError(t, err)

		res, err := cl.ClientInvokeUnary(ctx, &ClientInvokeUnary_Request{
			MethodDesc: &MethodDesc{
				Name: testpb.TestService_UnimplementedCall_FullMethodName,
			},
			Payload: payload,
		})

		require.NoError(t, err)
		require.NotNil(t, res.Error)
		assert.Equal(t, GRPCErrCode_UNIMPLEMENTED, res.Error.GrpcErrorCode)
		// assert.Equal(t, errcode.ErrTestEcho, res.Error.ErrorCode)
		// require.NotNil(t, res.Error.ErrorDetails)
		// assert.Greater(t, len(res.Error.ErrorDetails.Codes), 0)
	}
}

func TestUnaryAdvancedServiceError(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cl := createBridgeTestingClient(t, ctx)

	// call `testutil.TestService/EchoTest` with empty request,
	{
		in := &testpb.SimpleRequest{
			ResponseStatus: &testpb.EchoStatus{
				Code:    42,
				Message: messageStringTest,
			},
		}

		payload, err := proto.Marshal(in)
		require.NoError(t, err)

		res, err := cl.ClientInvokeUnary(ctx, &ClientInvokeUnary_Request{
			MethodDesc: &MethodDesc{
				Name: testpb.TestService_UnaryCall_FullMethodName,
			},
			Payload: payload,
		})

		require.NoError(t, err)
		require.NotNil(t, res.Error)
		assert.Equal(t, GRPCErrCode(42), res.Error.GrpcErrorCode)
		assert.Equal(t, messageStringTest, res.Error.Message)

		// assert.Equal(t, errcode.ErrTestEcho, res.Error.ErrorCode)
		// require.NotNil(t, res.Error.ErrorDetails)
		// assert.Greater(t, len(res.Error.ErrorDetails.Codes), 0)
	}
}

func TestStreamService(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cl := createBridgeTestingClient(t, ctx)

	// call instance `MessengerService/EchoTest`
	var streamid string
	responses := make([]*testpb.ResponseParameters, 10)
	for i := range responses {
		responses[i] = &testpb.ResponseParameters{Size: messageTestSize}
	}

	{
		in := &testpb.StreamingOutputCallRequest{
			ResponseParameters: responses,
		}

		payload, err := proto.Marshal(in)
		require.NoError(t, err)

		res, err := cl.CreateClientStream(ctx, &ClientCreateStream_Request{
			MethodDesc: &MethodDesc{
				Name:           testgrpc.TestService_HalfDuplexCall_FullMethodName,
				IsServerStream: true,
			},
			Payload: payload,
		})

		require.NoError(t, err)
		require.NotEmpty(t, res.StreamId)

		streamid = res.StreamId
		assert.NotEmpty(t, res.StreamId)
	}

	// test stream reply
	{
		for i := 0; i < 10; i++ {
			res, err := cl.ClientStreamRecv(ctx, &ClientStreamRecv_Request{
				StreamId: streamid,
			})
			require.NoError(t, err)
			require.NotNil(t, res.Error)
			require.False(t, res.Eof)
			require.Equal(t, GRPCErrCode_OK, res.Error.GrpcErrorCode)
			// assert.Equal(t, errcode.Undefined, res.Error.ErrorCode)

			out := new(testpb.StreamingOutputCallResponse)
			err = proto.Unmarshal(res.Payload, out)
			require.NoError(t, err)

			assert.Equal(t, messageTestSize, len(out.GetPayload().GetBody()))
		}
	}

	// test close stream
	{
		res, err := cl.ClientStreamClose(ctx, &ClientStreamClose_Request{
			StreamId: streamid,
		})

		require.NoError(t, err)
		require.NotNil(t, res.Error)
		assert.Equal(t, res.Error.GrpcErrorCode, GRPCErrCode_OK)
		// assert.Equal(t, res.Error.ErrorCode, errcode.Undefined)
	}
}

func TestStreamServiceEOF(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cl := createBridgeTestingClient(t, ctx)

	// call instance `MessengerService/EchoTest`
	var streamid string
	// responses := make([]*testpb.ResponseParameters, 10)
	// for i := range responses {
	// 	responses[i] = &testpb.ResponseParameters{Size: messageTestSize}
	// }

	{
		in := &testpb.StreamingOutputCallRequest{
			// ResponseParameters: responses,
			ResponseStatus: &testpb.EchoStatus{
				Code: int32(codes.Internal),
			},
		}

		payload, err := proto.Marshal(in)
		require.NoError(t, err)

		res, err := cl.CreateClientStream(ctx, &ClientCreateStream_Request{
			MethodDesc: &MethodDesc{
				Name:           testgrpc.TestService_HalfDuplexCall_FullMethodName,
				IsServerStream: true,
			},
			Payload: payload,
		})

		require.NoError(t, err)
		require.NotEmpty(t, res.StreamId)

		streamid = res.StreamId
		assert.NotEmpty(t, res.StreamId)
	}

	// test stream reply
	{
		res, err := cl.ClientStreamRecv(ctx, &ClientStreamRecv_Request{
			StreamId: streamid,
		})
		require.NoError(t, err)
		require.NotNil(t, res.Error)
		require.True(t, res.Eof)
		require.NotEqual(t, GRPCErrCode_OK, res.Error.GrpcErrorCode)
	}
}

// func TestStreamServiceError(t *testing.T) {
// }

// @TODO:
// func TestDuplexStreamService(t *testing.T) {
// }

// @TODO:
// func TestDuplexStreamServiceError(t *testing.T) {
// }

func createBridgeTestingClient(t *testing.T, ctx context.Context) RPCManagerClient {
	t.Helper()

	srv := grpc.NewServer()
	t.Cleanup(srv.Stop)

	svc := NewService(&Options{})

	l := bufconn.Listen(4096)

	go srv.Serve(l)

	RegisterRPCManagerServer(srv, svc)

	ts := interop.NewTestServer()
	testgrpc.RegisterTestServiceServer(srv, ts)

	cc, err := grpc.DialContext(ctx, "", grpc.WithDialer(func(string, time.Duration) (net.Conn, error) {
		return l.DialContext(ctx)
	}), grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)

	for serviceName := range srv.GetServiceInfo() {
		svc.RegisterService(serviceName, cc)
	}

	return NewRPCManagerClient(cc)
}
