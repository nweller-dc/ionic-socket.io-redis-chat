// Copyright 2012 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ipv4_test

import (
	"net"
	"testing"
)

func isUnicast(ip net.IP) bool {
	return ip.To4() != nil && (ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsGlobalUnicast())
}

// LoopbackInterface returns a logical network interface for loopback
// tests.
func loopbackInterface() *net.Interface {
	ift, err := net.Interfaces()
	if err != nil {
		return nil
	}
	for _, ifi := range ift {
		if ifi.Flags&net.FlagLoopback == 0 || ifi.Flags&net.FlagUp == 0 {
			continue
		}
		ifat, err := ifi.Addrs()
		if err != nil {
			continue
		}
		for _, ifa := range ifat {
			switch ifa := ifa.(type) {
			case *net.IPAddr:
				if isUnicast(ifa.IP) {
					return &ifi
				}
			case *net.IPNet:
				if isUnicast(ifa.IP) {
					return &ifi
				}
			}
		}
	}
	return nil
}

// isMulticastAvailable returns true if ifi is a multicast access
// enabled network interface.  It also returns a unicast IPv4 address
// that can be used for listening on ifi.
func isMulticastAvailable(ifi *net.Interface) (net.IP, bool) {
	if ifi == nil || ifi.Flags&net.FlagUp == 0 || ifi.Flags&net.FlagMulticast == 0 {
		return nil, false
	}
	ifat, err := ifi.Addrs()
	if err != nil {
		return nil, false
	}
	for _, ifa := range ifat {
		switch ifa := ifa.(type) {
		case *net.IPAddr:
			if isUnicast(ifa.IP) {
				return ifa.IP, true
			}
		case *net.IPNet:
			if isUnicast(ifa.IP) {
				return ifa.IP, true
			}
		}
	}
	return nil, false
}

func acceptor(t *testing.T, ln net.Listener, done chan<- bool) {
	defer func() { done <- true }()

	c, err := ln.Accept()
	if err != nil {
		t.Errorf("net.Listener.Accept failed: %v", err)
		return
	}
	c.Close()
}
