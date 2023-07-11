package lazy

import (
	"fmt"

	"google.golang.org/grpc/encoding"
)

// Codec is basically a no-op grpc.Codec use to pass LazyMessage through
// grpc
type Codec struct{}

func NewCodec() *Codec { return &Codec{} }

// Codec
var _ encoding.Codec = (*Codec)(nil)

func (lc *Codec) Marshal(value interface{}) ([]byte, error) {
	if lm, ok := value.(*Message); ok {
		return lm.buf, nil
	}

	return nil, fmt.Errorf("lazy-codec marshal: message is not lazy")
}

func (*Codec) Unmarshal(buf []byte, value interface{}) error {
	if lm, ok := value.(*Message); ok {
		lm.buf = buf
		return nil
	}

	return fmt.Errorf("lazy-codec unmarshal: message is not lazy")
}

func (lc *Codec) String() string { return "lazy-codec" }
func (lc *Codec) Name() string   { return lc.String() }
