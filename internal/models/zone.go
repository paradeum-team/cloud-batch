package models

type ZonesResponse struct {
	Zones []Zone `json:"data"`
	Limit int    `json:"limit"`
	Total int    `json:"total"`
}

type Zone struct {
	ID               string `json:"id"`
	CloudEnv         string `json:"cloud_env"`
	ExternalID       string `json:"external_id"`
	Name             string `json:"name"`
	Networks         int    `json:"networks"`
	Provider         string `json:"provider"`
	Region           string `json:"region"`
	RegionExtID      string `json:"region_ext_id"`
	RegionExternalID string `json:"region_external_id"`
	RegionID         string `json:"region_id"`
}
