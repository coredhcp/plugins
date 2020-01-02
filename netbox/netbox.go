// Copyright 2018-present the CoreDHCP Authors. All rights reserved
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package netbox

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
)

// NetBox wraps endpoint and API token.
type NetBox struct {
	Endpoint *url.URL
	APIToken string
}

// Request sends an API request to the specified NetBox API method.
func (nb *NetBox) Request(path string) ([]byte, error) {
	// quick and dirty copy
	u, _ := url.Parse(nb.Endpoint.String())
	u.Path += "/api/" + path
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", nb.APIToken))
	req.Header.Set("Accept", "application/json; indent=4")
	cl := http.Client{}
	resp, err := cl.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// GetIPs returns the IP addresses associated to a given MAC address,
// if any.
func (nb *NetBox) GetIPs(mac string) ([]net.IPNet, error) {
	mac = strings.ToLower(mac)
	resp, err := nb.Request("dcim/interfaces")
	if err != nil {
		return nil, err
	}
	var interfacesReply netboxDCIMInterfacesReply
	if err := json.Unmarshal(resp, &interfacesReply); err != nil {
		return nil, err
	}
	found := false
	enabled := false
	var deviceID uint
	for _, result := range interfacesReply.Results {
		if strings.ToLower(result.MACAddress) == mac {
			found = true
			enabled = result.Enabled
			deviceID = result.Device.ID
			break
		}
	}
	if !found {
		return nil, fmt.Errorf("mac %s not found", mac)
	}
	if !enabled {
		return nil, fmt.Errorf("device with mac %s is not enabled", mac)
	}
	resp, err = nb.Request(fmt.Sprintf("dcim/devices/%d", deviceID))
	if err != nil {
		return nil, err
	}
	var deviceReply netboxDCIMDeviceReply
	if err := json.Unmarshal(resp, &deviceReply); err != nil {
		return nil, err
	}
	ips := make([]net.IPNet, 0)
	if deviceReply.PrimaryIP4 != nil {
		ip, mask, err := parseAddr(deviceReply.PrimaryIP4.Address)
		if err != nil {
			return nil, fmt.Errorf("failed to parse address '%s': %v", deviceReply.PrimaryIP4.Address, err)
		}
		ips = append(ips, net.IPNet{IP: ip, Mask: mask})
	}
	if deviceReply.PrimaryIP6 != nil {
		ip, mask, err := parseAddr(deviceReply.PrimaryIP6.Address)
		if err != nil {
			return nil, fmt.Errorf("failed to parse address '%s': %v", deviceReply.PrimaryIP6.Address, err)
		}
		ips = append(ips, net.IPNet{IP: ip, Mask: mask})
	}
	log.Debugf("Found %d IP(s) for mac %s: %v", len(ips), mac, ips)
	return ips, nil
}
