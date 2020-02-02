package pipe

import (
	"reflect"
	"testing"
)

func TestNewLazyMessage(t *testing.T) {
	tests := []struct {
		name string
		want *LazyMessage
	}{
		{"dummy", &LazyMessage{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewLazyMessage(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewLazyMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLazyMessage_Base64(t *testing.T) {
	type fields struct {
		buf []byte
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"empty", fields{[]byte{}}, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &LazyMessage{
				buf: tt.fields.buf,
			}
			if got := m.Base64(); got != tt.want {
				t.Errorf("LazyMessage.Base64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLazyMessage_FromBase64(t *testing.T) {
	type fields struct {
		buf []byte
	}
	type args struct {
		b64 string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantLm  *LazyMessage
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &LazyMessage{
				buf: tt.fields.buf,
			}
			gotLm, err := m.FromBase64(tt.args.b64)
			if (err != nil) != tt.wantErr {
				t.Errorf("LazyMessage.FromBase64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotLm, tt.wantLm) {
				t.Errorf("LazyMessage.FromBase64() = %v, want %v", gotLm, tt.wantLm)
			}
		})
	}
}

func TestLazyMessage_Bytes(t *testing.T) {
	type fields struct {
		buf []byte
	}
	tests := []struct {
		name   string
		fields fields
		want   []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &LazyMessage{
				buf: tt.fields.buf,
			}
			if got := m.Bytes(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LazyMessage.Bytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLazyMessage_FromBytes(t *testing.T) {
	type fields struct {
		buf []byte
	}
	type args struct {
		buf []byte
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *LazyMessage
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &LazyMessage{
				buf: tt.fields.buf,
			}
			if got := m.FromBytes(tt.args.buf); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LazyMessage.FromBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLazyMessage_Reset(t *testing.T) {
	type fields struct {
		buf []byte
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &LazyMessage{
				buf: tt.fields.buf,
			}
			m.Reset()
		})
	}
}

func TestLazyMessage_String(t *testing.T) {
	type fields struct {
		buf []byte
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &LazyMessage{
				buf: tt.fields.buf,
			}
			if got := m.String(); got != tt.want {
				t.Errorf("LazyMessage.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLazyMessage_ProtoMessage(t *testing.T) {
	type fields struct {
		buf []byte
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &LazyMessage{
				buf: tt.fields.buf,
			}
			m.ProtoMessage()
		})
	}
}
