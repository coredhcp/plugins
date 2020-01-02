// Copyright 2018-present the CoreDHCP Authors. All rights reserved
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package netbox

import "encoding/json"

// netboxDCIMInterfacesReply maps to the JSON response of NetBox's
// /api/dcim/interfaces API endpoint.
type netboxDCIMInterfacesReply struct {
	Count   uint                         `json:"count"`
	Next    string                       `json:"next"`
	Prev    string                       `json:"prev"`
	Results []netboxDCIMInterfacesResult `json:"results"`
}

type netboxDCIMInterfacesResult struct {
	ID                    uint                             `json:"id"`
	Device                netboxDCIMInterfacesResultDevice `json:"device"`
	Name                  string                           `json:"name"`
	Type                  netboxValue                      `json:"type"`
	FormFactor            netboxValue                      `json:"form_factor"`
	Enabled               bool                             `json:"enabled"`
	Lag                   *json.RawMessage                 `json:"lag,omitempty"`
	MTU                   uint                             `json:"mtu"`
	MACAddress            string                           `json:"mac_address"`
	MGMTOnly              bool                             `json:"mgmt_only"`
	Description           string                           `json:"description"`
	ConnectedEndpointType *json.RawMessage                 `json:"connected_endpoint_type,omitempty"`
	ConnectedEndpoint     *json.RawMessage                 `json:"connected_endpoint,omitempty"`
	ConnectionStatus      *json.RawMessage                 `json:"connection_status,omitempty"`
	Cable                 *json.RawMessage                 `json:"cable,omitempty"`
	Mode                  *json.RawMessage                 `json:"mode,omitempty"`
	UntaggedVLAN          *json.RawMessage                 `json:"untagged_vlan,omitempty"`
	TaggedVLANs           []string                         `json:"tagged_vlans,omitempty"`
	Tags                  []string                         `json:"tags"`
	CountIPAddresses      uint                             `json:"count_ip_addresses"`
}

type netboxDCIMInterfacesResultDevice struct {
	ID          uint   `json:"id"`
	URL         string `json:"url"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
}

type netboxValue struct {
	Value uint   `json:"value"`
	Label string `json:"label"`
}

// netboxDCIMDeviceReply maps to the JSON response of NetBox's
// /api/dcim/devices/<id> API endpoint.
type netboxDCIMDeviceReply struct {
	ID               uint               `json:"id"`
	Name             string             `json:"name"`
	DisplayName      string             `json:"display_name"`
	DeviceType       *json.RawMessage   `json:"device_type,omitempty"`
	DeviceRole       *json.RawMessage   `json:"device_role,omitempty"`
	Tenant           *string            `json:"tenant,omitempty"`
	Platform         *json.RawMessage   `json:"platform,omitempty"`
	Serial           string             `json:"serial"`
	AssetTag         *json.RawMessage   `json:"asset_tag,omitempty"`
	Site             *json.RawMessage   `json:"site,omitempty"`
	Rack             *json.RawMessage   `json:"rack,omitempty"`
	Position         *json.RawMessage   `json:"position,omitempty"`
	Face             *json.RawMessage   `json:"face,omitempty"`
	ParentDevice     *json.RawMessage   `json:"parent_device,omitempty"`
	Status           netboxValue        `json:"status"`
	PrimaryIP        *netboxIP          `json:"primary_ip,omitempty"`
	PrimaryIP4       *netboxIP          `json:"primary_ip4,omitempty"`
	PrimaryIP6       *netboxIP          `json:"primary_ip6,omitempty"`
	Cluster          *json.RawMessage   `json:"cluster,omitempty"`
	VirtualChassis   *json.RawMessage   `json:"virtual_chassis,omitempty"`
	VCPosition       *json.RawMessage   `json:"vc_position,omitempty"`
	VCPriority       *json.RawMessage   `json:"vc_priority,omitempty"`
	Comments         string             `json:"comments"`
	LocalContextData *json.RawMessage   `json:"local_context_data,omitempty"`
	Tags             []*json.RawMessage `json:"tags,omitempty"`
	CustomFields     *json.RawMessage   `json:"custom_fields,omitempty"`
	ConfigContext    *json.RawMessage   `json:"config_context,omitempty"`
	Created          string             `json:"created"`
	LastUpdated      string             `json:"last_updated"`
}

type netboxIP struct {
	ID      uint   `json:"id"`
	URL     string `json:"url"`
	Family  uint   `json:"family"`
	Address string `json:"address"`
}
