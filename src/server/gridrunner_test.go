package server

import (
	"net"
	"testing"
)

func TestCreateAddressListner(t *testing.T) {
	t.Parallel()
	type args struct {
		netType     string
		bindAddress string
	}
	tests := []struct {
		name string
		args args
		val  func(t *testing.T, n net.Listener, err error)
	}{
		{"localhost", args{"tcp", "127.0.0.1:9234"}, func(t *testing.T, n net.Listener, err error) {
			if err != nil {
				t.Fatalf("got an error:%v", err)
			}
			if n == nil {
				t.Fatalf("expected to get a listner not nil?")
			}
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := createAddressListner(tt.args.netType, tt.args.bindAddress)
			tt.val(t, got, err)
		})
	}
}
