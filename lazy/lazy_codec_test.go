package pipe

import (
	"reflect"
	"testing"
)

func TestNewLazyCodec(t *testing.T) {
	tests := []struct {
		name string
		want *LazyCodec
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewLazyCodec(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewLazyCodec() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLazyCodec_Marshal(t *testing.T) {
	type args struct {
		value interface{}
	}
	tests := []struct {
		name    string
		lc      *LazyCodec
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lc := &LazyCodec{}
			got, err := lc.Marshal(tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("LazyCodec.Marshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LazyCodec.Marshal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLazyCodec_Unmarshal(t *testing.T) {
	type args struct {
		buf   []byte
		value interface{}
	}
	tests := []struct {
		name    string
		l       *LazyCodec
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &LazyCodec{}
			if err := l.Unmarshal(tt.args.buf, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("LazyCodec.Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// no revelant
// func TestLazyCodec_String(t *testing.T)

// no revelant
// func TestLazyCodec_Name(t *testing.T)
