package servie_cloudpods

import (
	"cloud-batch/configs"
	"cloud-batch/internal/models"
	"cloud-batch/internal/pkg/db/gleveldb"
	"fmt"
	"net/url"
	"testing"
)

func TestLoginCloudPods(t *testing.T) {
	auth, err := Login()
	if err != nil {
		t.Errorf("%v", err)
		return
	}

	cookieJson, err := gleveldb.Get(authKey)
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println(string(cookieJson))
	fmt.Println(auth)
}

func TestGetAuth(t *testing.T) {
	auth, err := GetAuth()
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(auth)
}

func TestListRegions(t *testing.T) {
	regionsResponse, cloudErr, err := ListRegions(
		map[string]string{
			"filter":   "country_code.equals('cn')",
			"usable":   "false",
			"limit":    "2048",
			"offset":   "0",
			"provider": models.ProviderAliyun,
		}, nil)
	if err != nil {
		t.Error(cloudErr)
		t.Error(err)
		return
	}

	fmt.Println(regionsResponse.Regions)
}

func TestListZones(t *testing.T) {
	zonesResponse, _, err := ListZones(
		map[string]string{
			"filter":   "external_id.like('%/cn-%')",
			"usable":   "false",
			"limit":    "2048",
			"offset":   "0",
			"provider": models.ProviderAliyun,
			"order":    "asc",
		}, nil)
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println(zonesResponse.Zones)
}

func TestGetVpc(t *testing.T) {
	vpc, cloudErr, err := GetVpc(
		"cloudBatch-e0587cb2-b07b-43e0-8f7e-611d50f892ad",
		map[string]string{
			//"region_id": "d1897432-7cb8-4896-8278-7ca0c6c5eb95",
			"provider": models.ProviderAliyun,
		})
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(cloudErr)
	fmt.Println(vpc)
}

func TestDeleteVpc(t *testing.T) {
	cloudErr, err := DeleteVpc("cloudBatch")
	if err != nil {
		t.Error(cloudErr)
		t.Error(err)
	}
}

func TestCreateRegionsVpc(t *testing.T) {
	cloudErr, err := CreateRegionsVpc(models.ProviderAliyun)
	if err != nil {
		t.Error(cloudErr)
		t.Error(err)
	}
}

func TestCreateVpc(t *testing.T) {
	vpcName := fmt.Sprintf("%s", configs.Cloudpods.VpcPreName)
	vpcCreate := &models.VpcCreate{
		CidrBlock:          configs.Cloudpods.VpcCidrBlock,
		CloudproviderID:    models.ProviderAliyun,
		CloudregionID:      "d1897432-7cb8-4896-8278-7ca0c6c5eb95",
		Description:        "cloud batch default network",
		ExternalAccessMode: "eip",
		IsPublic:           true,
		Name:               vpcName,
	}
	vpc, cloudErr, err := CreateVpc(vpcCreate)
	if err != nil {
		t.Error(cloudErr)
		t.Error(err)
	}

	fmt.Println(vpc)
}

func TestCreateZonesNetwork(t *testing.T) {
	regionsResponse, cloudErr, err := ListRegions(
		map[string]string{
			"usable":   "false",
			"limit":    "2048",
			"offset":   "0",
			"provider": models.ProviderAliyun,
			"order":    "asc",
		}, nil)
	if err != nil {
		t.Error(cloudErr)
		t.Error(err)
		return
	}

	for _, region := range regionsResponse.Regions {
		cloudErr, err := CreateZonesNetworkByRegionID(region.ID, models.ProviderAliyun)
		if err != nil {
			t.Error(cloudErr)
			t.Error(err)
		}
	}
}

func TestDeleteNetwork(t *testing.T) {
	cloudErr, err := DeleteNetwork("cloudBatch")
	if err != nil {
		t.Error(cloudErr)
		t.Error(err)
	}
}

func TestListServerSkus(t *testing.T) {
	queryParams := map[string]string{
		"zone_id":        "a5170b98-b705-4813-82c1-62d3f4291341",
		"cpu_core_count": "4",
		"memory_size_mb": "8192",
		"provider":       models.ProviderAliyun,
		"order_by":       "created_at",
		"order":          "desc",
	}
	urlValues := url.Values{}

	urlValues.Add("filter", "instance_type_family.startswith('c')")
	urlValues.Add("filter", "sys_disk_type.contains('cloud_efficiency')")
	urlValues.Add("filter", "instance_type_family.notequals('c6r')")
	urlValues.Add("filter", "instance_type_family.notequals('c7r')")
	resp, cloudErr, err := ListServerSkus(queryParams, urlValues)

	if err != nil {
		t.Error(cloudErr)
		t.Error(err)
	}
	fmt.Println(resp)
}

func TestListServers(t *testing.T) {
	queryParams := map[string]string{
		"tags.0.key": "user:project",
		//"tags.0.value": "bfs",
	}

	urlValues := url.Values{}
	urlValues.Add("field", "id")
	urlValues.Add("field", "name")
	urlValues.Add("field", "eip")
	urlValues.Add("field", "ips")
	urlValues.Add("field", "zone_ext_id")
	urlValues.Add("field", "provider")
	urlValues.Add("field", "status")
	respBody, cloudErr, err := ListServers(queryParams, urlValues)
	if err != nil {
		t.Error(cloudErr)
		t.Error(err)
	}
	fmt.Printf("%+s\n", respBody)
}

func TestCreateServer(t *testing.T) {

	project := "bfs"

	disks := []*models.Disk{{
		DiskType: "sys",
		Index:    0,
		Backend:  models.AliyunDiskCloudEfficiency,
		Size:     102400,
		ImageID:  "6851cdd4-95a4-4484-8b98-a61c6e61164a",
	}}

	nets := []*models.Net{{
		Network: "d3d0184e-cfaf-4250-83e2-873e59200eaa",
	}}

	serverCreate := &models.ServerCreate{
		Meta:               map[string]string{"user:project": project},
		AutoStart:          true,
		GenerateName:       "cloudBatch-test",
		Hypervisor:         models.ProviderAliyun,
		Count:              1,
		DisableDelete:      false,
		Disks:              disks,
		Nets:               nets,
		PreferZone:         "bdc32422-e090-4052-821b-47d42ff21394",
		Sku:                "ecs.c5.xlarge",
		PublicIPChargeType: models.PublicIPChargeTraffic,
		PublicIPBw:         100,
		Keypair:            "masterkey",
		Secgroups:          []string{"sg-cloudbatch-all"},
	}

	server, cloudErr, err := CreateServer(serverCreate)
	if err != nil {
		t.Error(cloudErr)
		t.Error(err)
	}
	fmt.Println(server)
}

func TestDeleteServer(t *testing.T) {
	queryParams := map[string]string{
		"OverridePendingDelete": "true",
	}
	cloudErr, err := DeleteServer("13137a0b-eb33-4ac9-8f12-7ff9d17671a2", queryParams)
	if err != nil {
		t.Error(cloudErr)
		t.Error(err)
	}
}

func TestCreateServers(t *testing.T) {
	createServersForm := models.BatchCreateServersForm{
		RegionScope:  models.RegionScopeChinaMainland,
		CpuCoreCount: 4,
		MemorySizeMb: 8192,
		Project:      "bfs",
		Provider:     models.ProviderAliyun,
		Count:        1,
	}
	resp, cloudErr, err := BatchCreateServers(createServersForm)
	if err != nil {
		t.Error(cloudErr)
		t.Error(err)
	}
	fmt.Printf("ids len: %d\n", len(resp.Ids))
	fmt.Printf("%+v\n", resp)
}
