package lazy

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewLazyCodec(t *testing.T) {
	lc := NewCodec()
	assert.NotNil(t, lc)
}

func TestMarshal(t *testing.T) {
	lc := NewCodec()
	lm := &Message{}
	lm.FromBytes([]byte("test"))

	bytes, err := lc.Marshal(lm)
	assert.Nil(t, err)
	assert.Equal(t, []byte("test"), bytes)

	_, err = lc.Marshal("test")
	assert.NotNil(t, err)
}

func TestUnmarshal(t *testing.T) {
	lc := NewCodec()
	lm := &Message{}

	err := lc.Unmarshal([]byte("test"), lm)
	assert.Nil(t, err)
	assert.Equal(t, []byte("test"), lm.buf)

	err = lc.Unmarshal([]byte("test"), "test")
	assert.NotNil(t, err)
}
