package net

import (
	"errors"
	"net"
)

func GetIP(ignoreLoopback bool) (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, iface := range ifaces {
		// interface down
		if iface.Flags&net.FlagUp == 0 {
			continue
		}

		// lookback - do not use
		if ignoreLoopback && iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		addrs, err := iface.Addrs()

		// interface returned error
		if err != nil {
			return "", err
		}

		// addresses given
		for _, addr := range addrs {
			var ip net.IP

			switch v := addr.(type) {
			case *net.IPNet: // network
				ip = v.IP
			case *net.IPAddr: // ip address
				ip = v.IP
			}

			// lookback again
			if ip == nil || (ignoreLoopback && ip.IsLoopback()) {
				continue
			}

			// convert to v4.
			ip = ip.To4()
			if ip == nil {
				continue
			}
			return ip.String(), nil
		}
	}

	return "", errors.New("are you connected to network? ")
}
