package lazy

import (
	"context"
	"sync/atomic"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

var (
	lazyCodec = NewCodec()
	streamids uint64
)

type Client struct {
	cc *grpc.ClientConn
}

type MethodDesc struct {
	Name          string
	ClientStreams bool
	ServerStreams bool
}

func NewClient(cc *grpc.ClientConn) *Client {
	return &Client{cc: cc}
}

func (lc *Client) InvokeUnary(ctx context.Context, desc *MethodDesc, in *Message, copts ...grpc.CallOption) (out *Message, err error) {
	out = NewMessage()
	err = grpc.Invoke(ctx, desc.Name, in, out, lc.cc, append(copts, grpc.ForceCodec(lazyCodec))...)
	return
}

func (lc *Client) InvokeStream(ctx context.Context, desc *MethodDesc, in *Message, copts ...grpc.CallOption) (*Stream, error) {
	sdesc := &grpc.StreamDesc{
		StreamName:    desc.Name,
		ServerStreams: desc.ServerStreams,
		ClientStreams: desc.ClientStreams,
	}

	sctx, cancel := context.WithCancel(ctx)
	cstream, err := grpc.NewClientStream(sctx, sdesc, lc.cc, desc.Name, append(copts, grpc.ForceCodec(lazyCodec))...)
	if err != nil {
		cancel()
		return nil, err
	}

	if !desc.ClientStreams && desc.ServerStreams {
		if err := cstream.SendMsg(in); err != nil {
			cancel()
			return nil, err
		}

		if err := cstream.CloseSend(); err != nil {
			cancel()
			return nil, err
		}
	}

	return &Stream{
		id:           atomic.AddUint64(&streamids, 1),
		ClientStream: cstream,
		CancelFunc:   cancel,
	}, nil
}

type Stream struct {
	// used to close the stream
	context.CancelFunc
	grpc.ClientStream

	id uint64
}

func (s *Stream) SendMsg(in proto.Message) (err error) {
	if err = s.ClientStream.SendMsg(in); err != nil {
		s.CancelFunc()
	}

	return
}

func (s *Stream) RecvMsg(out proto.Message) (err error) {
	if err = s.ClientStream.RecvMsg(out); err != nil {
		s.CancelFunc()
	}

	return
}

func (s *Stream) Close() (err error) {
	err = s.CloseSend()
	s.CancelFunc()
	return
}

func (s *Stream) ID() uint64 {
	return s.id
}
