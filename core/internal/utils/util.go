package utils

import "net"

func GetIPv4() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	for _, addr := range addrs {
		if ip, ok := addr.(*net.IPNet); ok && !ip.IP.IsLoopback() {
			if ip.IP.To4() == nil {
				continue
			}
			return ip.IP.String(), nil
		}
	}
	return "", err
}
