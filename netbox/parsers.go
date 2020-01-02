// Copyright 2018-present the CoreDHCP Authors. All rights reserved
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package netbox

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

func parseAddr(a string) (net.IP, net.IPMask, error) {
	parts := strings.SplitN(a, "/", 2)
	if len(parts) != 2 {
		return nil, nil, fmt.Errorf("invalid address: got %d components, want 2", len(parts))
	}
	ip := net.ParseIP(parts[0])
	if ip == nil {
		return nil, nil, fmt.Errorf("invalid IP address '%s'", parts[0])
	}
	ones, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid netmask length '%s': %v", parts[1], err)
	}
	bits := net.IPv6len * 8
	if v := ip.To4(); v != nil {
		bits = net.IPv4len * 8
	}
	return ip, net.CIDRMask(int(ones), bits), nil
}
