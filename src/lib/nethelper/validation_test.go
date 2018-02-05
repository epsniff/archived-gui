package nethelper

import (
	"strings"
	"testing"
)

func TestValidateAddress(t *testing.T) {
	t.Parallel()
	type args struct {
		address string
	}
	tests := []struct {
		name     string
		args     args
		validate func(t *testing.T, got string, err error)
	}{
		{"default bind value`:5503`", args{":5503"}, func(t *testing.T, got string, err error) {
			if err != nil {
				t.Fatalf("no error expected but got:%v", err)
			}
			if got == "" || got == ":5503" || got == "127.0.0.1:5503" || got == "localhost:5503" {
				t.Fatalf("got an unexpected value back:%v", got)
			}
		}},
		{"localhost`:5503`", args{"localhost:5503"}, func(t *testing.T, got string, err error) {
			if err == nil {
				t.Fatalf("expected an error but got nil")
			}
			e := err.Error()
			if !strings.Contains(e, "bad hostname in address") {
				t.Fatalf("expected an error for bad hostname, but got: %v", e)
			}
		}},
		//
		{"default bind value`192.111.111.111:80`", args{"192.111.111.111:80"}, func(t *testing.T, got string, err error) {
			if err != nil {
				t.Fatalf("no error expected but got:%v", err)
			}
			if got == "" || got == ":80" || got == "127.0.0.1:80" || got != "192.111.111.111:80" {
				t.Fatalf("got an unexpected value back:%v", got)
			}
		}},
		//
		{"default bind value`[2001:db8:85a3:8d3:1319:8a2e:370:7348]:1234`", args{"[2001:db8:85a3:8d3:1319:8a2e:370:7348]:1234"}, func(t *testing.T, got string, err error) {
			if err != nil {
				t.Fatalf("no error expected but got:%v", err)
			}
			if got == "" || got != "[2001:db8:85a3:8d3:1319:8a2e:370:7348]:1234" {
				t.Fatalf("got an unexpected value back:%v", got)
			}
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ValidateAddress(tt.args.address)
			tt.validate(t, got, err)
			t.Logf("testcase `%v` got:`%v`", tt.name, got)
		})
	}
}
