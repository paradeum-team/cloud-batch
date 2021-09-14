package service_cloudpods

import (
	"cloud-batch/configs"
	"cloud-batch/internal/models"
	"cloud-batch/internal/pkg/db/gredis"
	"cloud-batch/internal/pkg/e"
	"cloud-batch/internal/pkg/logging"
	"cloud-batch/pkg/utils"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/gogf/gf/encoding/gjson"
	"github.com/pkg/errors"
	"gopkg.in/redis.v5"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	authKey                          = "cloud-auth-cookie"
	loginPath                        = "/api/v1/auth/login"
	cloudRegionsPath                 = "/api/v2/cloudregions"
	zonesPath                        = "/api/v2/zones"
	vpcsPath                         = "/api/v2/vpcs"
	networksPath                     = "/api/v2/networks"
	serverSkusPath                   = "/api/v2/serverskus"
	serversPath                      = "/api/v2/servers"
	ServerCreateServersStatus        = "batchCreateingServers-status"
	ServerCreateServersStatusTTL     = time.Hour * 1
	ServerCreateServersStatusDoneTTL = time.Hour * 24
	ServerCreating                   = "creating"
	ServerCreateDone                 = "done"
	ServerCreateTimeout              = "timeout"
	ServerCreateError                = "error"
	ShortServersResponseSaveKeyPre   = "shortServersResponse"
)

func GetServerCreateError(status string) error {
	switch status {
	case ServerCreating:
		return e.ErrServerCreating
	case ServerCreateTimeout:
		return e.ErrServerCreateTimeout
	case ServerCreateError:
		return e.ErrServerCreateError
	case ServerCreateDone:
		return nil
	case "":
		return e.ErrServerCreateStatusEmpty
	default:
		return e.ErrUnknownError
	}
}

func NewRestyClient(insecureSkipVerify bool) *resty.Client {
	client := resty.New()
	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: insecureSkipVerify})
	client.R().
		SetHeader("Content-Type", "application/json")
	return client
}

func NewRestyAuthClient() (*resty.Client, error) {
	auth, err := GetAuth()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			auth, err = Login()
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	client := NewRestyClient(true)
	client.SetHeader("cookie", auth)
	return client, nil
}

// Login cloudpods
func Login() (string, error) {
	client := NewRestyClient(true)

	resp, err := client.R().
		SetBody(fmt.Sprintf(`{"username":"%s", "password":"%s", "domain":"%s"}`,
			configs.Cloudpods.Username, configs.Cloudpods.Password, configs.Cloudpods.Domain)).
		Post(fmt.Sprintf("%s%s", configs.Cloudpods.BaseUrl, loginPath))
	if err != nil {
		return "", err
	}
	if resp.StatusCode() != http.StatusOK {
		return "", fmt.Errorf(string(resp.Body()))
	}

	var authName = "yunionauth"
	var cookie *http.Cookie

	for _, item := range resp.Cookies() {
		if item.Name == authName {
			cookie = item
		}
	}

	if cookie == nil {
		return "", errors.New("response cookie not found yunionauth")
	}

	cookiesJson, err := json.Marshal(cookie)
	if err != nil {
		return "", err
	}

	err = gredis.Set(authKey, string(cookiesJson), time.Hour*24)
	if err != nil {
		return "", err
	}

	return cookie.Raw, nil
}

func GetAuth() (string, error) {
	cookieByte, err := gredis.Get(authKey).Bytes()
	if err != nil {
		return "", err
	}

	cookie := new(http.Cookie)

	err = json.Unmarshal(cookieByte, cookie)
	if err != nil {
		return "", err
	}

	// cookie 过期重新登录, 留出1分钟缓冲时间
	if cookie.Expires.Sub(time.Now()) < time.Minute {
		authRaw, err := Login()
		if err != nil {
			return "", err
		}
		return authRaw, nil
	}

	return cookie.Raw, nil
}

// ListRegions /api/v2/cloudregions
func ListRegions(queryParams map[string]string, urlValues url.Values) (*models.RegionsResponse, *models.CloudError, error) {
	client, err := NewRestyAuthClient()
	if err != nil {
		return nil, nil, err
	}

	defaultQueryParams := map[string]string{
		"enabled":   "true",
		"cloud_env": "public",
		"status":    "inservice",
		"filter":    "external_id.notequals('Aliyun/cn-nanjing')",
		"order":     "asc",
	}
	resp, err := client.R().
		SetQueryParams(defaultQueryParams).
		SetQueryParams(queryParams).
		SetQueryParamsFromValues(urlValues).
		Get(fmt.Sprintf("%s/%s", configs.Cloudpods.BaseUrl, cloudRegionsPath))
	if err != nil {
		return nil, nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		cloudErr, err := generateCloudError(resp.Body())
		return nil, cloudErr, err
	}
	regionResp := new(models.RegionsResponse)

	err = json.Unmarshal(resp.Body(), regionResp)
	if err != nil {
		return nil, nil, err
	}

	return regionResp, nil, nil
}

// ListZones /api/v2/zones
func ListZones(queryParams map[string]string, urlValues url.Values) (*models.ZonesResponse, *models.CloudError, error) {
	client, err := NewRestyAuthClient()
	if err != nil {
		return nil, nil, err
	}

	defaultQueryParams := map[string]string{
		"enabled":   "true",
		"cloud_env": "public",
		"status":    "enable",
	}
	resp, err := client.R().
		SetQueryParams(defaultQueryParams).
		SetQueryParams(queryParams).
		SetQueryParamsFromValues(urlValues).
		Get(fmt.Sprintf("%s/%s", configs.Cloudpods.BaseUrl, zonesPath))
	if err != nil {
		return nil, nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		cloudErr, err := generateCloudError(resp.Body())
		return nil, cloudErr, err
	}

	zonesResp := new(models.ZonesResponse)

	err = json.Unmarshal(resp.Body(), zonesResp)
	if err != nil {
		return nil, nil, err
	}

	return zonesResp, nil, nil
}

func generateCloudError(body []byte) (*models.CloudError, error) {
	cloudErr := new(models.CloudError)
	err := json.Unmarshal(body, cloudErr)
	if err != nil {
		return nil, err
	}

	return cloudErr, errors.New(cloudErr.Details)
}

func GetVpc(idOrName string, queryParams map[string]string) (*models.Vpc, *models.CloudError, error) {
	client, err := NewRestyAuthClient()
	if err != nil {
		return nil, nil, err
	}

	defaultQueryParams := map[string]string{
		"cloud_env": "public",
	}

	resp, err := client.R().
		SetQueryParams(defaultQueryParams).
		SetQueryParams(queryParams).
		Get(fmt.Sprintf("%s/%s/%s", configs.Cloudpods.BaseUrl, vpcsPath, idOrName))
	if err != nil {
		return nil, nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		cloudErr, err := generateCloudError(resp.Body())
		return nil, cloudErr, err
	}

	vpc := new(models.Vpc)

	err = json.Unmarshal(resp.Body(), vpc)
	if err != nil {
		return nil, nil, err
	}

	return vpc, nil, nil
}

func CreateVpc(vpcCreate *models.VpcCreate) (*models.Vpc, *models.CloudError, error) {
	client, err := NewRestyAuthClient()
	if err != nil {
		return nil, nil, err
	}

	vpcCreateJson, err := json.Marshal(vpcCreate)
	if err != nil {
		return nil, nil, err
	}

	resp, err := client.R().
		SetBody(vpcCreateJson).
		Post(fmt.Sprintf("%s/%s", configs.Cloudpods.BaseUrl, vpcsPath))

	if err != nil {
		return nil, nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		cloudErr, err := generateCloudError(resp.Body())
		return nil, cloudErr, err
	}

	vpc := new(models.Vpc)
	err = json.Unmarshal(resp.Body(), vpc)
	if err != nil {
		return nil, nil, err
	}

	logging.Logger.Infof("create vpc %s in region %s done.", vpc.Name, vpc.Region)
	return vpc, nil, nil
}

func CreateRegionsVpc(provider string) (*models.CloudError, error) {
	regionResp, _, err := ListRegions(
		map[string]string{
			"usable":   "false",
			"limit":    "2048",
			"offset":   "0",
			"provider": provider,
		}, nil)
	if err != nil {
		return nil, err
	}

	for _, region := range regionResp.Regions {
		// 没有 vpc 创建 vpc
		vpcName := fmt.Sprintf("%s-%s", configs.Cloudpods.VpcPreName, region.ID)
		vpcCreate := &models.VpcCreate{
			CidrBlock:          configs.Cloudpods.VpcCidrBlock,
			CloudproviderID:    region.Provider,
			CloudregionID:      region.ID,
			Description:        "cloud batch default network",
			ExternalAccessMode: "eip",
			IsPublic:           true,
			Name:               vpcName,
		}
		// 没有 vpc 创建
		if region.VpcCount == 0 {
			if _, _, err := CreateVpc(vpcCreate); err != nil {
				return nil, err
			}
			continue
		}

		// 有 vpc ，没有 cloudBatch vpc, 创建
		_, cloudErr, err := GetVpc(vpcName, map[string]string{
			"region_id": region.ID,
			"provider":  region.Provider,
		})
		if cloudErr != nil && cloudErr.Code == 404 {
			if _, cloudErr, err := CreateVpc(vpcCreate); err != nil {
				return cloudErr, err
			}
			continue
		}

		if err != nil {
			return nil, err
		}
		logging.Logger.Infof("region %s has a VPC %s", region.Name, vpcCreate.Name)

	}
	return nil, nil
}

func DeleteVpc(idOrName string) (*models.CloudError, error) {
	client, err := NewRestyAuthClient()
	if err != nil {
		return nil, err
	}

	resp, err := client.R().Delete(fmt.Sprintf("%s/%s/%s", configs.Cloudpods.BaseUrl, vpcsPath, idOrName))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		cloudErr, err := generateCloudError(resp.Body())
		return cloudErr, err
	}

	return nil, nil

}

func GetNetwork(idOrName string, queryParams map[string]string) (*models.Network, *models.CloudError, error) {
	client, err := NewRestyAuthClient()
	if err != nil {
		return nil, nil, err
	}

	defaultQueryParams := map[string]string{
		"cloud_env": "public",
	}

	resp, err := client.R().
		SetQueryParams(defaultQueryParams).
		SetQueryParams(queryParams).
		Get(fmt.Sprintf("%s/%s/%s", configs.Cloudpods.BaseUrl, networksPath, idOrName))
	if err != nil {
		return nil, nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		cloudErr, err := generateCloudError(resp.Body())
		return nil, cloudErr, err
	}

	network := new(models.Network)

	err = json.Unmarshal(resp.Body(), network)
	if err != nil {
		return nil, nil, err
	}

	return network, nil, nil
}

func CreateNetwork(networkCreate *models.NetworkCreate) (*models.Network, *models.CloudError, error) {
	client, err := NewRestyAuthClient()
	if err != nil {
		return nil, nil, err
	}

	networkCreateJson, err := json.Marshal(networkCreate)
	if err != nil {
		return nil, nil, err
	}

	resp, err := client.R().
		SetBody(networkCreateJson).
		Post(fmt.Sprintf("%s/%s", configs.Cloudpods.BaseUrl, networksPath))

	if err != nil {
		return nil, nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		cloudErr, err := generateCloudError(resp.Body())
		return nil, cloudErr, err
	}

	network := new(models.Network)
	err = json.Unmarshal(resp.Body(), network)
	if err != nil {
		return nil, nil, err
	}

	logging.Logger.Infof("create network %s in zone %s done.", network.Name, network.Zone)
	return network, nil, nil
}

func CreateZonesNetworkByRegionID(regionID, provider string) (*models.CloudError, error) {
	zonesResp, cloudErr, err := ListZones(
		map[string]string{
			"usable":    "false",
			"limit":     "2048",
			"offset":    "0",
			"provider":  provider,
			"region_id": regionID,
			"order":     "asc",
		}, nil)
	if err != nil {
		return cloudErr, err
	}

	for i, zone := range zonesResp.Zones {
		networkName := fmt.Sprintf("%s-%s", configs.Cloudpods.NetworkPreName, strings.ReplaceAll(zone.ExternalID, "/", "-"))
		networkCreate := &models.NetworkCreate{
			GuestIpPrefix: fmt.Sprintf("%s.%d.0/24", configs.Cloudpods.NetworkGuestIpPrefix, i),
			Name:          networkName,
			Vpc:           fmt.Sprintf("%s-%s", configs.Cloudpods.VpcPreName, regionID),
			Zone:          zone.ID,
		}

		// 删除
		//cloudErr,err := DeleteNetwork(fmt.Sprintf("%s-%s",configs.Cloudpods.NetworkPreName, zone.ID))
		//if err != nil {
		//	logging.Logger.Warnf("zone %s delete a Network %s: %v", zone.Name, networkName, err)
		//}else{
		//	logging.Logger.Infof("zone %s delete a Network %s", zone.Name, networkName)
		//}
		//continue

		// 如果没有 Network 创建
		//if zone.Networks == 0 {
		//	if _, cloudErr, err := CreateNetwork(networkCreate); err != nil {
		//		return cloudErr, fmt.Errorf("zone: %s, zone_id: %s, network: %s , err: %v", zone.Name, zone.ID, networkName, err)
		//	}
		//	continue
		//}

		var network *models.Network
		// 有 Network ，没有 cloud-batch 创建
		network, cloudErr, err = GetNetwork(networkName, map[string]string{
			"zone_id":  zone.ID,
			"provider": zone.Provider,
		})

		if cloudErr != nil && cloudErr.Code == http.StatusNotFound {
			if _, cloudErr, err := CreateNetwork(networkCreate); err != nil {
				return cloudErr, fmt.Errorf("zone: %s, zone_id: %s, network: %s , err: %v", zone.Name, zone.ID, networkName, err)
			}
			continue
		}

		if err != nil {
			return nil, err
		}

		// 异常状态的 Network 删除重建
		if network != nil && network.Status != models.NetworkStatusAvailable {
			cloudErr, err := DeleteNetwork(network.ID)
			if err != nil {
				return nil, fmt.Errorf("DeleteNetwork cloudErr: %v, err: %v", cloudErr, err)
			}
			if _, cloudErr, err := CreateNetwork(networkCreate); err != nil {
				return cloudErr, errors.Errorf("zone: %s, zone_id: %s, network: %s , err: %v", zone.Name, zone.ID, networkName, err)
			}
			continue
		}
		logging.Logger.Infof("zone %s has a Network %s", zone.Name, networkName)

	}

	return nil, nil
}

func DeleteNetwork(idOrName string) (*models.CloudError, error) {
	client, err := NewRestyAuthClient()
	if err != nil {
		return nil, err
	}

	resp, err := client.R().Delete(fmt.Sprintf("%s/%s/%s", configs.Cloudpods.BaseUrl, networksPath, idOrName))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		cloudErr, err := generateCloudError(resp.Body())
		return cloudErr, err
	}

	return nil, nil

}

func ListServerSkus(queryParams map[string]string, urlValues url.Values) (*models.ServerSkusResponse, *models.CloudError, error) {
	client, err := NewRestyAuthClient()
	if err != nil {
		return nil, nil, err
	}

	defaultQueryParams := map[string]string{
		"enabled":   "true",
		"cloud_env": "public",
		"status":    "ready",
	}

	resp, err := client.R().
		SetQueryParams(defaultQueryParams).
		SetQueryParams(queryParams).
		SetQueryParamsFromValues(urlValues). // 重复的quere param value
		Get(fmt.Sprintf("%s/%s", configs.Cloudpods.BaseUrl, serverSkusPath))
	if err != nil {
		return nil, nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		cloudErr, err := generateCloudError(resp.Body())
		return nil, cloudErr, err
	}

	serverSkusResp := new(models.ServerSkusResponse)

	err = json.Unmarshal(resp.Body(), serverSkusResp)
	if err != nil {
		return nil, nil, err
	}

	return serverSkusResp, nil, nil
}

func ListServers(queryParams map[string]string, urlValues url.Values) ([]byte, *models.CloudError, error) {
	client, err := NewRestyAuthClient()
	if err != nil {
		return nil, nil, err
	}

	defaultQueryParams := map[string]string{
		"enabled":      "true",
		"with_meta":    "true",
		"tags.0.key":   "user:cloudBatch",
		"tags.0.value": "true",
	}

	resp, err := client.R().
		SetQueryParams(defaultQueryParams).
		SetQueryParams(queryParams).
		SetQueryParamsFromValues(urlValues). // 重复的quere param value
		Get(fmt.Sprintf("%s/%s", configs.Cloudpods.BaseUrl, serversPath))
	if err != nil {
		return nil, nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		cloudErr, err := generateCloudError(resp.Body())
		return nil, cloudErr, err
	}

	return resp.Body(), nil, nil
}

func CreateServer(serverCreate *models.ServerCreate) (*models.Server, *models.CloudError, error) {
	client, err := NewRestyAuthClient()
	if err != nil {
		return nil, nil, err
	}

	serverCreateJson, err := json.Marshal(serverCreate)
	if err != nil {
		return nil, nil, err
	}

	resp, err := client.R().
		SetBody(serverCreateJson).
		Post(fmt.Sprintf("%s/%s", configs.Cloudpods.BaseUrl, serversPath))

	if err != nil {
		return nil, nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		cloudErr, err := generateCloudError(resp.Body())
		return nil, cloudErr, err
	}

	server := new(models.Server)
	err = json.Unmarshal(resp.Body(), server)
	if err != nil {
		return nil, nil, err
	}

	logging.Logger.Infof("create server %s in zone %s done.", server.Name, server.ZoneExtID)
	return server, nil, nil
}

func DeleteServer(idOrName string, queryParams map[string]string) ([]byte, *models.CloudError, error) {
	client, err := NewRestyAuthClient()
	if err != nil {
		return nil, nil, err
	}

	resp, err := client.R().
		SetQueryParams(queryParams).
		Delete(fmt.Sprintf("%s/%s/%s", configs.Cloudpods.BaseUrl, serversPath, idOrName))
	if err != nil {
		return nil, nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		cloudErr, err := generateCloudError(resp.Body())
		return nil, cloudErr, err
	}

	return resp.Body(), nil, nil
}

func BatchDeleteServers(deleteServersForm models.BatchDeleteServersForm) (serverCount, doneCount int, errIDs []string, err error) {
	deleteQueryParams := map[string]string{
		"OverridePendingDelete": "true",
	}

	// 如果 有 ids 数据，忽略其它参数
	if deleteServersForm.IDs != nil && len(deleteServersForm.IDs) > 0 && deleteServersForm.IDs[0] != "" {
		for _, id := range deleteServersForm.IDs {
			_, cloudErr, err := DeleteServer(id, deleteQueryParams)
			if err != nil {
				logging.Logger.Errorf("DeleteServer id: %s, cloudErr: %v, err: %v", id, cloudErr, err)
				errIDs = append(errIDs, id)
			}
			doneCount++
		}
		return len(deleteServersForm.IDs), doneCount, errIDs, nil
	}

	// 如果 ids 数据 ，按 其它参数 查询 servers
	urlValues := url.Values{}

	// 查询过滤字段
	urlValues.Set("provider", deleteServersForm.Provider)
	urlValues.Set("tags.1.key", "user:project")
	if deleteServersForm.Project != "" {
		urlValues.Set("tags.1.value", deleteServersForm.Project)
	}
	if deleteServersForm.BatchNumber != "" {
		urlValues.Set("tags.2.key", "user:batchNumber")
		urlValues.Set("tags.2.value", deleteServersForm.BatchNumber)
	}

	resp, _, err := ListServers(nil, urlValues)
	if err != nil {
		return 0, 0, nil, err
	}

	shortServerResponse := new(models.ShortServersResponse)
	err = json.Unmarshal(resp, shortServerResponse)
	if err != nil {
		return 0, 0, nil, err
	}

	if shortServerResponse.Total == 0 {
		return shortServerResponse.Total, 0, nil, nil
	}

	for _, server := range shortServerResponse.Servers {
		_, cloudErr, err := DeleteServer(server.ID, deleteQueryParams)
		if err != nil {
			logging.Logger.Errorf("DeleteServer  id: %s, cloudErr: %v, err: %v", server.ID, cloudErr, err)
			errIDs = append(errIDs, server.ID)
		}
		doneCount++
	}
	return shortServerResponse.Total, doneCount, errIDs, nil
}

// BatchCreateServers 批量创建主机
func BatchCreateServers(batchCreateServersForm models.BatchCreateServersForm) (batchCreateServersResponse *models.BatchCreateServersResponse, cloudErr *models.CloudError, err error) {

	zonesQueryParams := map[string]string{
		"usable":   "false",
		"limit":    "2048",
		"offset":   "0",
		"provider": batchCreateServersForm.Provider,
	}

	zonesUrlValues := url.Values{}
	// 过滤主机类型
	instanceTypeFamilys := "'c7','c6','c6a','c5','c4'"
	zonesUrlValues.Add("joint_filter", fmt.Sprintf("serverskus.zone_id(id).instance_type_family.in(%s)", instanceTypeFamilys))
	zonesUrlValues.Add("joint_filter", fmt.Sprintf("serverskus.zone_id(id).sys_disk_type.contains.('%s')", models.AliyunDiskCloudEfficiency))
	zonesUrlValues.Add("joint_filter", fmt.Sprintf("serverskus.zone_id(id).sys_disk_type.cpu_core_count.equals.('%d')", batchCreateServersForm.CpuCoreCount))
	zonesUrlValues.Add("joint_filter", fmt.Sprintf("serverskus.zone_id(id).sys_disk_type.memory_size_mb.equals.('%d')", batchCreateServersForm.MemorySizeMb))
	switch batchCreateServersForm.RegionScope {
	case models.RegionScopeChinaMainland:
		zonesUrlValues.Add("filter", fmt.Sprintf("external_id.like('%s/cn-%s')", batchCreateServersForm.Provider, "%"))
	default:
		zonesUrlValues.Add("filter", fmt.Sprintf("external_id.like('%s/%s-%s')", batchCreateServersForm.Provider, batchCreateServersForm.RegionScope, "%"))
	}

	zonesResp, cloudErr, err := ListZones(zonesQueryParams, zonesUrlValues)
	if err != nil {
		return nil, cloudErr, err
	}

	if zonesResp.Zones == nil || len(zonesResp.Zones) == 0 {
		return nil, nil, errors.New("not fount zones")
	}

	// 最终可用区列表
	zones := map[string]models.Zone{}
	// zones  数组转为 map, 乱序
	for _, zone := range zonesResp.Zones {
		// 因为 cloudpods 查询 zones 列表不支持排除region_ext_id, 所以如果只查询 中国大陆 这里排除 cn-hongkong 得到 zones
		if batchCreateServersForm.RegionScope == models.RegionScopeChinaMainland && zone.RegionExtID == "cn-hongkong" {
			continue
		}
		zones[zone.ID] = zone
	}

	disks := []*models.Disk{{
		DiskType: "sys",
		Index:    0,
		Backend:  models.AliyunDiskCloudEfficiency,
		Size:     102400,
		ImageID:  configs.Cloudpods.AliyunDefaultImageId,
	}}

	// 结果 ids
	var ids []string
	batchNumber := utils.NowNanoTimeStamp()

	err = gredis.Set(fmt.Sprintf("%s-%s", ServerCreateServersStatus, batchNumber), ServerCreating, ServerCreateServersStatusTTL)
	if err != nil {
		return nil, nil, err
	}
	// 余数
	remainder := batchCreateServersForm.Count % len(zones)
	// 商
	quotient := batchCreateServersForm.Count / len(zones)

	i := 0
	for _, zone := range zones {
		logging.Logger.Debugf("zone: %s start", zone.Name)
		var serverCount int
		if i < remainder {
			// 不能整除的多余的主机前面的zone挨个添加1
			serverCount = quotient + 1
		} else if quotient == 0 {
			// 总数小于zones 数量， i 大于余数退出循环
			break
		} else {
			serverCount = quotient
		}
		i++

		skusQueryParams := map[string]string{
			"zone_id":        zone.ID,
			"cpu_core_count": utils.IntToString(batchCreateServersForm.CpuCoreCount),
			"memory_size_mb": utils.IntToString(batchCreateServersForm.MemorySizeMb),
			"provider":       batchCreateServersForm.Provider,
			"order_by":       "created_at",
			"order":          "desc",
		}
		urlValues := url.Values{}
		urlValues.Add("filter", fmt.Sprintf("sys_disk_type.contains('%s')", models.AliyunDiskCloudEfficiency))
		urlValues.Add("filter", fmt.Sprintf("instance_type_family.in(%s)", instanceTypeFamilys))
		skusResp, cloudErr, err := ListServerSkus(skusQueryParams, urlValues)

		if err != nil {
			logging.Logger.Errorf("BatchCreateServers ListServerSkus cloudErr: %+v, err: %+v", cloudErr, err)
			continue
		}
		if skusResp.Total == 0 {
			logging.Logger.Errorf("BatchCreateServers ListServerSkus len 0, skusQueryParams: %v, urlValues: %v", skusQueryParams, urlValues)
			continue
		}

		networkName := fmt.Sprintf("%s-%s", configs.Cloudpods.NetworkPreName, strings.ReplaceAll(zone.ExternalID, "/", "-"))
		nets := []*models.Net{{
			Network: networkName,
		}}

		for j := 1; j <= serverCount; j++ {
			serverCreate := &models.ServerCreate{
				Meta: map[string]string{
					"user:project":     batchCreateServersForm.Project,
					"user:batchNumber": batchNumber,
					"user:cloudBatch":  "true",
				},
				AutoStart:          true,
				GenerateName:       fmt.Sprintf("%s-%s-%s-%d", batchCreateServersForm.Project, batchNumber, zone.ID, j),
				Hypervisor:         batchCreateServersForm.Provider,
				Count:              1,
				DisableDelete:      false,
				Disks:              disks,
				Nets:               nets,
				PreferZone:         zone.ID,
				Sku:                skusResp.ServerSkus[0].Name,
				PublicIPChargeType: models.PublicIPChargeTraffic,
				PublicIPBw:         100,
				Keypair:            configs.Cloudpods.DefaultKeypair,
				Secgroups:          []string{configs.Cloudpods.DefaultSecgroup},
			}
			server, cloudErr, err := CreateServer(serverCreate)
			if err != nil {
				logging.Logger.Errorf("BatchCreateServers CreateServer: %v cloudErr: %+v, err: %+v", serverCreate, cloudErr, err)
				continue
			}
			ids = append(ids, server.ID)
			logging.Logger.Infof("zone: %s, create server: %s, sku: %s, serverCount: %d", zone.Name, serverCreate.GenerateName, skusResp.ServerSkus[0].InstanceTypeFamily, serverCount)
		}
		logging.Logger.Debugf("zone: %s end", zone.Name)
	}

	batchCreateServersResponse = &models.BatchCreateServersResponse{
		Ids:         ids,
		Count:       len(ids),
		Project:     batchCreateServersForm.Project,
		Provider:    batchCreateServersForm.Provider,
		BatchNumber: batchNumber,
	}

	go updateCreateServerStatus(batchNumber, batchCreateServersForm.Count)
	return batchCreateServersResponse, nil, err
}

func updateCreateServerStatus(batchNumber string, createCount int) {
	var (
		shortServersResponse *models.ShortServersResponse
		err                  error
	)
	// 按批号查询  running、deploy_fail、disk_file 两种状态的主机，一直到 与createCount相同退出
	for i := 0; i < 60; i++ {
		if i >= 59 {
			err = gredis.Set(fmt.Sprintf("%s-%s", ServerCreateServersStatus, batchNumber), ServerCreateTimeout, ServerCreateServersStatusTTL)
			if err != nil {
				logging.Logger.Errorf("gredis.Set %s", err)
			}
			logging.Logger.Errorf("%s 创建主机查询状态超时", batchNumber)
			break
		}
		time.Sleep(time.Second * 5)

		shortServersResponse, err = QueryCreateServersTotal(batchNumber, []string{"running", "deploy_fail", "disk_fail", "sched_fail"})

		if err != nil {
			logging.Logger.Errorf("QueryCreateServersTotal err: %v", err)
			continue
		}

		if shortServersResponse == nil || shortServersResponse.Total == 0 {
			logging.Logger.Error("shortServersResponse is nil")
			continue
		}

		if shortServersResponse.Total == createCount {
			isErr := false
			for _, s := range shortServersResponse.Servers {
				if s.Eip == "" {
					isErr = true
					break
				}
			}
			if isErr {
				continue
			}
			break
		}
	}

	serversJson, err := json.Marshal(shortServersResponse)
	if err != nil {
		logging.Logger.Errorf("json.Marshal(shortServersResponse) err: %v", err)
		return
	}

	err = gredis.Set(fmt.Sprintf("%s-%s", ShortServersResponseSaveKeyPre, batchNumber), string(serversJson), ServerCreateServersStatusDoneTTL)
	if err != nil {
		logging.Logger.Errorf("json.Marshal(shortServersResponse) err: %v", err)
		return
	}

	err = gredis.Set(fmt.Sprintf("%s-%s", ServerCreateServersStatus, batchNumber), ServerCreateDone, ServerCreateServersStatusDoneTTL)
	if err != nil {
		logging.Logger.Errorf("redis.Set %s", err)
	}
}

func QueryCreateServersTotal(batchNumber string, status []string) (*models.ShortServersResponse, error) {
	urlValues := url.Values{}
	// 查询过滤字段
	// 默认提供 tags.0.key tags.0.values 所以这里从1开始
	if batchNumber != "" {
		urlValues.Set("tags.1.key", "user:batchNumber")
		urlValues.Set("tags.1.value", batchNumber)
	}

	// eip 不能为空
	urlValues.Add("filter", "eip.isnullorempty('')")

	for _, s := range status {
		urlValues.Add("status", s)
		urlValues.Add("status", s)
	}
	resp, _, err := ListServers(nil, urlValues)
	if err != nil {
		return nil, err
	}
	shortServersResponse := new(models.ShortServersResponse)
	err = json.Unmarshal(resp, shortServersResponse)
	if err != nil {
		return nil, err
	}

	return shortServersResponse, nil
}

func DeleteFailServer() error {

	// 查询部署失败的主机
	shortServersResponse, err := QueryCreateServersTotal("", []string{"deploy_fail"})
	if err != nil {
		return err
	}

	// 删除部署失败主机
	for _, server := range shortServersResponse.Servers {
		_, _, err := DeleteServer(server.ID, nil)
		if err != nil {
			logging.Logger.Errorf("DeleteServer %s err :%v ", server.ID, err)
		}
	}
	return nil
}

func GetCreateServersByBatchNumber(batchNumber string) (*models.ShortServersResponse, error) {
	shortServersResponseJson, err := gredis.Get(fmt.Sprintf("%s-%s", ShortServersResponseSaveKeyPre, batchNumber)).Result()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	shortServersResponse := new(models.ShortServersResponse)
	err = gjson.DecodeTo(shortServersResponseJson, shortServersResponse)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if shortServersResponse.Total == 0 {
		return nil, errors.New("BfsUpdateValuesByServers shortServersResponse total is 0")
	}
	return shortServersResponse, nil
}
