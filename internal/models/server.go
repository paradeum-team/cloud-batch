package models

import "time"

const (
	AliyunDiskCloudEfficiency = "cloud_efficiency"
	AliyunDiskCloudSsd        = "cloud_ssd"
	AliyunDiskCloudEssd       = "cloud_essd"

	PublicIPChargeTraffic   = "traffic"
	PublicIPChargeBandwidth = "bandwidth"
)

type BatchGetServersForm struct {
	Project     string   `json:"project" binding:"required,min=2,max=20" enums:"bfs,dfs" default:"bfs" required:"true"`
	Provider    string   `json:"provider" binding:"required,min=2,max=20" enums:"aliyun" default:"aliyun" required:"true"`
	BatchNumber string   `json:"batch_number" binding:"required,len=14" minLength:"14" maxLength:"14"`
	Status      []string `json:"status" binding:"required" enums:"running,deploying,starting,deploy_fail,unknown"`
	Limit       int      `json:"limit" binding:"min=1"`
	Offset      int      `json:"offset" binding:"min=0"`
}

type BatchCreateServersForm struct {
	//   区域前缀: all 全部, cnml 中国大陆, cn 中国， ap 亚太， us 美洲， eu 欧洲, me 中东
	RegionScope string `json:"region_scope" binding:"required,min=2,max=10" enums:"all,cnml,cn,us,ap,me,eu" default:"cnml" required:"true"`
	// cpu 核数
	CpuCoreCount int `json:"cpu_core_count" binding:"required,gte=1,lte=16" enums:"1,2,4,8,16" default:"4" required:"true"`
	// 内存（单位 M)
	MemorySizeMb int `json:"memory_size_mb" binding:"required,gte=512,lte=16384" enums:"512,1024,2048,4096,8192,16384" default:"8192" required:"true"`
	// 项目标签
	Project string `json:"project" binding:"required,min=2,max=20" enums:"bfs,dfs" default:"bfs" required:"true"`
	// 云提供者
	Provider string `json:"provider" binding:"required,min=2,max=20" enums:"aliyun" default:"aliyun" required:"true"`
	// 创建数量
	Count int `json:"count" binding:"required,gte=1,lte=100" enums:"1,2,5,10,20,50,100" default:"1" required:"true"`
}

type BatchDeleteServersForm struct {
	// 主机 ids, 如果设置了 ids, 忽略其它参数
	IDs []string `json:"ids"`
	// 区域前缀: all 全部, cnml 中国大陆, cn 中国， ap 亚太， us 美洲， eu 欧洲, me 中东
	RegionScope string `json:"region_scope" binding:"min=2,max=10" enums:"all,cnml,cn,us,ap,me,eu" default:"cnml"`
	// 创建主机的批号
	BatchNumber string `json:"batch_number" binding:"min=0,max=14" minLength:"14" maxLength:"14"`
	// 项目标签
	Project string `json:"project" binding:"min=2,max=20" enums:"bfs,dfs,test" default:"test" required:"true"`
	// 云提供者
	Provider string `json:"provider" binding:"min=2,max=20" enums:"aliyun" default:"aliyun" required:"true"`
}

type BatchCreateServersResponse struct {
	Ids         []string `json:"ids"`
	Count       int      `json:"count"`
	Project     string   `json:"project"`
	Provider    string   `json:"provider"`
	BatchNumber string   `json:"batch_number"`
}

type ShortServersResponse struct {
	Servers []ShortServer `json:"data"`
	Limit   int           `json:"limit"`
	Total   int           `json:"total"`
}

type ShortServer struct {
	ID             string         `json:"id"`
	Name           string         `json:"name"`
	Eip            string         `json:"eip"`
	Ips            string         `json:"ips"`
	ZoneExtID      string         `json:"zone_ext_id"`
	ServerMetadata ServerMetadata `json:"metadata"`
	Provider       string         `json:"provider"`
	Status         string         `json:"status"`
}

type ServerMetadata struct {
	UserBatchNumber string `json:"user:batchNumber"`
	UserProject     string `json:"user:project"`
}

type ServersResponse struct {
	Servers []Server `json:"data"`
	Limit   int      `json:"limit"`
	Total   int      `json:"total"`
}
type Nics struct {
	IPAddr    string `json:"ip_addr"`
	IsExit    bool   `json:"is_exit"`
	Mac       string `json:"mac"`
	NetworkID string `json:"network_id"`
	VpcID     string `json:"vpc_id"`
}
type Secgroups struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Server struct {
	CanDelete             bool              `json:"can_delete"`
	CanRecycle            bool              `json:"can_recycle"`
	CanUpdate             bool              `json:"can_update"`
	CloudEnv              string            `json:"cloud_env"`
	CreatedAt             time.Time         `json:"created_at"`
	Deleted               bool              `json:"deleted"`
	DisableDelete         bool              `json:"disable_delete"`
	Disk                  int               `json:"disk"`
	DiskCount             int               `json:"disk_count"`
	Eip                   string            `json:"eip"`
	EipMode               string            `json:"eip_mode"`
	ExternalID            string            `json:"external_id"`
	Freezed               bool              `json:"freezed"`
	Hypervisor            string            `json:"hypervisor"`
	ID                    string            `json:"id"`
	InstanceType          string            `json:"instance_type"`
	Ips                   string            `json:"ips"`
	Keypair               string            `json:"keypair"`
	KeypairID             string            `json:"keypair_id"`
	Metadata              map[string]string `json:"metadata"`
	Name                  string            `json:"name"`
	Nics                  []Nics            `json:"nics"`
	OsType                string            `json:"os_type"`
	PendingDeleted        bool              `json:"pending_deleted"` // 资源是否处于回收站
	Provider              string            `json:"provider"`
	Region                string            `json:"region"`
	RegionExtID           string            `json:"region_ext_id"`
	RegionExternalID      string            `json:"region_external_id"`
	RegionID              string            `json:"region_id"`
	Secgroup              string            `json:"secgroup"`
	Secgroups             []Secgroups       `json:"secgroups"`
	SecgrpID              string            `json:"secgrp_id"`
	ShutdownBehavior      string            `json:"shutdown_behavior"`
	Source                string            `json:"source"`
	SrcIPCheck            bool              `json:"src_ip_check"`
	SrcMacCheck           bool              `json:"src_mac_check"`
	SshableLastState      bool              `json:"sshable_last_state"`
	Status                string            `json:"status"`
	VcpuCount             int               `json:"vcpu_count"`
	VmemSize              int               `json:"vmem_size"`
	Vpc                   string            `json:"vpc"`
	VpcExternalAccessMode string            `json:"vpc_external_access_mode"`
	VpcID                 string            `json:"vpc_id"`
	Zone                  string            `json:"zone"`
	ZoneExtID             string            `json:"zone_ext_id"`
	ZoneID                string            `json:"zone_id"`
}

type ServerCreate struct {
	Meta               map[string]string `json:"__meta__"`
	AutoStart          bool              `json:"auto_start"`
	GenerateName       string            `json:"generate_name"`
	Hypervisor         string            `json:"hypervisor"`
	Count              int               `json:"__count__"`
	DisableDelete      bool              `json:"disable_delete"`
	Disks              []*Disk           `json:"disks"`
	Nets               []*Net            `json:"nets"`
	PreferRegion       string            `json:"prefer_region"`
	PreferZone         string            `json:"prefer_zone"`
	VcpuCount          int               `json:"vcpu_count"`
	VmemSize           int               `json:"vmem_size"`
	Sku                string            `json:"sku"`
	PublicIPChargeType string            `json:"public_ip_charge_type"`
	PublicIPBw         int               `json:"public_ip_bw"`
	Keypair            string            `json:"keypair"`
	Secgroups          []string          `json:"secgroups"`
}

type Disk struct {
	DiskType string `json:"disk_type"`
	Index    int    `json:"index"`
	Backend  string `json:"backend"`
	Size     int    `json:"size"`
	ImageID  string `json:"image_id"`
}
type Net struct {
	Network string `json:"network"`
}
