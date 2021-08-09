package models

const (
	RegionScopeAll           = "all"
	RegionScopeChinaMainland = "cnml"
)

type RegionsResponse struct {
	Regions []Region `json:"data"`
	Limit   int      `json:"limit"`
	Total   int      `json:"total"`
}
type Region struct {
	City         string `json:"city"`
	CloudEnv     string `json:"cloud_env"`
	CountryCode  string `json:"country_code"`
	Enabled      bool   `json:"enabled"`
	Environment  string `json:"environment"`
	ExternalID   string `json:"external_id"`
	ID           string `json:"id"`
	Name         string `json:"name"`
	NetworkCount int    `json:"network_count"`
	Provider     string `json:"provider"`
	VpcCount     int    `json:"vpc_count"`
	ZoneCount    int    `json:"zone_count"`
}
