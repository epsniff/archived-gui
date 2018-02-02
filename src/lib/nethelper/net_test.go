package nethelper

import "testing"

func TestBindableIP(t *testing.T) {
	got, err := BindableIP()
	if err != nil {
		t.Errorf("BindableIP() error = %v", err)
		return
	}
	if got == "" || got == "127.0.0.1" {
		t.Errorf("BindableIP() = %v", got)
	}
}
