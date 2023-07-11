package lazy

import (
	"encoding/base64"
	"fmt"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// Message is basically a no-op `proto.Message` used to pass
// serialized message through grpc
type Message struct {
	buf []byte
}

func NewMessage() *Message {
	return &Message{}
}

func (m *Message) FromMessage(msg proto.Message) (*Message, error) {
	buf, err := proto.Marshal(msg)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal message: %w", err)
	}

	m.buf = buf
	return m, nil
}

func (m *Message) Base64() string {
	return base64.StdEncoding.EncodeToString(m.buf)
}

func (m *Message) FromBase64(b64 string) (lm *Message, err error) {
	m.buf, err = base64.StdEncoding.DecodeString(b64)
	lm = m
	return
}

func (m *Message) Bytes() []byte {
	return m.buf
}

func (m *Message) FromBytes(buf []byte) *Message {
	m.buf = make([]byte, len(buf))
	copy(m.buf, buf)
	return m
}

// ProtoReflect is not usefull with LazyMessage
func (m *Message) ProtoReflect() protoreflect.Message {
	return nil
}

type LazyMessageReflect *Message

// Message
var _ proto.Message = (*Message)(nil)

func (m *Message) Reset()         { *m = Message{} }
func (m *Message) String() string { return string(m.buf[:]) }
func (m *Message) ProtoMessage()  {}
