package models

type Vpc struct {
	CidrBlock        string `json:"cidr_block"`
	CloudEnv         string `json:"cloud_env"`
	ExternalID       string `json:"external_id"`
	ID               string `json:"id"`
	Name             string `json:"name"`
	NetworkCount     int    `json:"network_count"`
	Provider         string `json:"provider"`
	PublicScope      string `json:"public_scope"`
	Region           string `json:"region"`
	RegionExtID      string `json:"region_ext_id"`
	RegionExternalID string `json:"region_external_id"`
	RegionID         string `json:"region_id"`
	Status           string `json:"status"`
}

type VpcCreate struct {
	CidrBlock          string `json:"cidr_block"`
	CloudproviderID    string `json:"cloudprovider_id"`
	CloudregionID      string `json:"cloudregion_id"`
	Description        string `json:"description"`
	ExternalAccessMode string `json:"external_access_mode"`
	IsPublic           bool   `json:"is_public"`
	Name               string `json:"name"`
}
