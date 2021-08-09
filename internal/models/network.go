package models

const (
	NetworkStatusAvailable = "available"
)

type NetworkCreate struct {
	GuestIpPrefix string `json:"guest_ip_prefix"`
	Name          string `json:"name"`
	Vpc           string `json:"vpc"`  // vpc id or name
	Zone          string `json:"zone"` // zone id or name
}

type Network struct {
	CloudEnv         string `json:"cloud_env"`
	ExternalID       string `json:"external_id"`
	GuestGateway     string `json:"guest_gateway"`
	GuestIPEnd       string `json:"guest_ip_end"`
	GuestIPMask      int    `json:"guest_ip_mask"`
	GuestIPStart     string `json:"guest_ip_start"`
	ID               string `json:"id"`
	Name             string `json:"name"`
	Provider         string `json:"provider"`
	Region           string `json:"region"`
	RegionExtID      string `json:"region_ext_id"`
	RegionExternalID string `json:"region_external_id"`
	RegionID         string `json:"region_id"`
	Status           string `json:"status"`
	Vpc              string `json:"vpc"`
	VpcExtID         string `json:"vpc_ext_id"`
	VpcID            string `json:"vpc_id"`
	Wire             string `json:"wire"`
	WireID           string `json:"wire_id"`
	Zone             string `json:"zone"`
	ZoneID           string `json:"zone_id"`
}
