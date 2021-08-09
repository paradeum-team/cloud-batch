package models

type ServerSkusResponse struct {
	ServerSkus []ServerSku `json:"data"`
	Limit      int         `json:"limit"`
	Total      int         `json:"total"`
}
type ServerSku struct {
	CloudEnv             string `json:"cloud_env"`
	CPUCoreCount         int    `json:"cpu_core_count"`
	DataDiskTypes        string `json:"data_disk_types"`
	ExternalID           string `json:"external_id"`
	ID                   string `json:"id"`
	InstanceTypeCategory string `json:"instance_type_category"`
	InstanceTypeFamily   string `json:"instance_type_family"`
	MemorySizeMb         int    `json:"memory_size_mb"`
	Name                 string `json:"name"`
	Provider             string `json:"provider"`
	Region               string `json:"region"`
	RegionExtID          string `json:"region_ext_id"`
	RegionExternalID     string `json:"region_external_id"`
	RegionID             string `json:"region_id"`
	Status               string `json:"status"`
	SysDiskType          string `json:"sys_disk_type"`
	Zone                 string `json:"zone"`
	ZoneExtID            string `json:"zone_ext_id"`
	ZoneID               string `json:"zone_id"`
}
