// Copyright 2018-present the CoreDHCP Authors. All rights reserved
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

// Package netbox implements a plugin that retrieves IP lease information
// from a NetBox server and uses it in the DHCP reply.
//
// Example config:
//
// server6:
//   listen: '[::]547'
//   - example:
//   - server_id: LL aa:bb:cc:dd:ee:ff
//   - netbox: https://netbox.coredhcp.io my_api_token
//
// This will send requests to https://netbox.coredhcp.io/api/<path> using
// my_api_token for authentication.
package netbox

import (
	"fmt"
	"net"
	"net/url"

	"github.com/coredhcp/coredhcp/handler"
	"github.com/coredhcp/coredhcp/logger"
	"github.com/coredhcp/coredhcp/plugins"
	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/insomniacslk/dhcp/dhcpv6"
)

var log = logger.GetLogger("plugins/netbox")

// Plugin wraps plugin registration information
var Plugin = plugins.Plugin{
	Name:   "netbox",
	Setup6: setup6,
	Setup4: setup4,
}

var netBox *NetBox

func initNetBox(args ...string) error {
	if netBox != nil {
		// already initialized
		return nil
	}
	if len(args) != 2 {
		return fmt.Errorf("got %d arguments, want 2", len(args))
	}
	u, err := url.Parse(args[0])
	if err != nil {
		return fmt.Errorf("invalid URL '%s': %v", args[0], err)
	}
	t := args[1]
	netBox = &NetBox{
		Endpoint: u,
		APIToken: t,
	}
	return nil
}

func setup6(args ...string) (handler.Handler6, error) {
	if err := initNetBox(args...); err != nil {
		return nil, err
	}
	log.Info("Loaded netbox plugin for DHCPv6.")
	return netboxHandler6, nil
}

func setup4(args ...string) (handler.Handler4, error) {
	if err := initNetBox(args...); err != nil {
		return nil, err
	}
	log.Info("Loaded netbox plugin for DHCPv4.")
	return netboxHandler4, nil
}

func netboxHandler6(req, resp dhcpv6.DHCPv6) (dhcpv6.DHCPv6, bool) {
	log.Debugf("Received DHCPv6 packet: %s", req.Summary())
	mac, err := dhcpv6.ExtractMAC(req)
	if err != nil {
		log.Warningf("Could not find client MAC, dropping request")
		return resp, false
	}
	// extract the IA_ID from option IA_NA
	opt := req.GetOneOption(dhcpv6.OptionIANA)
	if opt == nil {
		log.Warningf("No option IA_NA found in request, dropping request")
		return resp, false
	}
	iaID := opt.(*dhcpv6.OptIANA).IaId
	log.Debugf("Retrieving IP addresses for MAC %s", mac)
	ips, err := netBox.GetIPs(mac.String())
	if err != nil {
		log.Warningf("No IPs found for MAC %s: %v", mac.String(), err)
		return resp, false
	}
	for _, addr := range ips {
		if addr.IP.To4() == nil && addr.IP.To16() != nil {
			resp.AddOption(&dhcpv6.OptIANA{
				IaId: iaID,
				Options: []dhcpv6.Option{
					&dhcpv6.OptIAAddress{
						IPv6Addr: addr.IP.To16(),
						// default lifetime, can be overridden by other plugins
						PreferredLifetime: 3600,
						ValidLifetime:     3600,
					},
				},
			})
			break
		}
	}
	log.Infof("Resp %s", resp.Summary())
	return resp, true
}

func netboxHandler4(req, resp *dhcpv4.DHCPv4) (*dhcpv4.DHCPv4, bool) {
	log.Debugf("Received DHCPv4 packet: %s", req.Summary())
	mac := req.ClientHWAddr.String()
	ips, err := netBox.GetIPs(mac)
	if err != nil {
		log.Warningf("No IPs found for MAC %s: %v", mac, err)
		return resp, false
	}
	for _, addr := range ips {
		if v := addr.IP.To4(); v != nil {
			resp.YourIPAddr = net.IP(v)
			break
		}
	}
	return resp, true
}
