package lazy

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewLazyMessage(t *testing.T) {
	lm := NewMessage()
	assert.NotNil(t, lm)
}

func TestBase64(t *testing.T) {
	lm := &Message{}
	lm.FromBytes([]byte("test"))

	b64 := lm.Base64()
	assert.Equal(t, "dGVzdA==", b64)
}

func TestFromBase64(t *testing.T) {
	lm := &Message{}
	newLm, err := lm.FromBase64("dGVzdA==")

	assert.Nil(t, err)
	assert.Equal(t, []byte("test"), newLm.buf)
}

func TestBytes(t *testing.T) {
	lm := &Message{}
	lm.FromBytes([]byte("test"))

	bytes := lm.Bytes()
	assert.Equal(t, []byte("test"), bytes)
}

func TestFromBytes(t *testing.T) {
	lm := &Message{}
	newLm := lm.FromBytes([]byte("test"))

	assert.Equal(t, []byte("test"), newLm.buf)
}

func TestString(t *testing.T) {
	lm := &Message{}
	lm.FromBytes([]byte("test"))

	str := lm.String()
	assert.Equal(t, "test", str)
}

func TestReset(t *testing.T) {
	lm := &Message{}
	lm.FromBytes([]byte("test"))
	lm.Reset()

	assert.Nil(t, lm.buf)
}

func TestProtoMessage(t *testing.T) {
	lm := &Message{}
	lm.ProtoMessage()

	// No assertion needed as ProtoMessage() doesn't return anything.
	// This test is just to ensure that ProtoMessage() can be called without panicking.
}

func TestProtoReflect(t *testing.T) {
	lm := &Message{}
	lm.ProtoReflect()

	// No assertion needed as ProtoMessage() doesn't return anything.
	// This test is just to ensure that ProtoMessage() can be called without panicking.
}
