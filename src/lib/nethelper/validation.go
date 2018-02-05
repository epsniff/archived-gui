package nethelper

import (
	"fmt"
	"net"
)

func ValidateAddress(address string) (string, error) {
	host, port, err := net.SplitHostPort(address)
	if err != nil {
		return "", fmt.Errorf("bad address:%v err:%v", address, err)
	}

	if host == "" {
		//default to leting us pick the IP address for you
		tmpIp, err := BindableIP()
		if err != nil {
			return "", fmt.Errorf("error finding a bindable IP address: err:%v", err)
		}
		host = tmpIp
		address = fmt.Sprintf("%s:%s", host, port)
	} else {
		if ip := net.ParseIP(host); ip == nil {
			return "", fmt.Errorf("bad hostname in address (grid only supports ip address): parsed_host:%v address:%v", host, address)
		}
	}

	return address, nil
}
