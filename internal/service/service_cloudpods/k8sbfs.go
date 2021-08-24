package service_cloudpods

import (
	"cloud-batch/configs"
	"cloud-batch/internal/models"
	"cloud-batch/internal/pkg/db/gredis"
	"cloud-batch/internal/pkg/e"
	"fmt"
	"github.com/gogf/gf/encoding/gyaml"
	"github.com/gogf/gf/os/gfile"
	"github.com/pkg/errors"
	"gopkg.in/redis.v5"
	"io/ioutil"
	"os"
	"time"
)

const (
	UpdateBfsValuesStatusKeyPre  = "updateBfsValuesStatus"
	UpdateBfsValuesStatusTTL     = time.Second * 10
	UpdateBfsValuesStatusDoneTTL = time.Hour * 24
)

func GetUpdateBfsValuesError(status string) error {
	switch status {
	case e.StatusStart:
		return e.ErrStatusStart
	case e.StatusError:
		return e.ErrStatusError
	case e.StatusConflict:
		return e.ErrStatusConflict
	case e.StatusDone:
		return e.ErrStatusDone
	case "":
		return nil
	default:
		return e.ErrUnknownError
	}
}

func BfsUpdateValuesByServers(batchNumber string) (*models.BfsValues, error) {
	serverCreateServersStatus, err := gredis.Get(fmt.Sprintf("%s-%s", ServerCreateServersStatus, batchNumber)).Result()
	if err != nil {
		// 如果redis 返回的不是 nil 错误，返回
		if !errors.Is(err, redis.Nil) {
			return nil, errors.WithStack(err)
		}
	}

	// 根据状态返回相应的错误
	err = GetServerCreateError(serverCreateServersStatus)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// 生成初始备份文件
	filePath := fmt.Sprintf("%s/%s", configs.K8sBfs.BfsPath, configs.K8sBfs.BfsValuesName)
	if !gfile.Exists(fmt.Sprintf("%s.init", filePath)) {
		err = gfile.CopyFile(filePath, fmt.Sprintf("%s.init", filePath))
		if err != nil {
			return nil, errors.WithStack(err)
		}
	}

	// 备份 bfs values
	err = gfile.CopyFile(filePath, fmt.Sprintf("%s.cloud-batch-bak", filePath))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	// 打开 bfs values
	values, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	bfsValues := new(models.BfsValues)
	err = gyaml.DecodeTo(values, bfsValues)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	updateBfsValuesStatus, err := gredis.Get(fmt.Sprintf("%s-%s", UpdateBfsValuesStatusKeyPre, batchNumber)).Result()
	if err != nil {
		// 如果redis 返回的不是 nil 错误，返回
		if !errors.Is(err, redis.Nil) {
			return nil, errors.WithStack(err)
		}
	}

	err = GetUpdateBfsValuesError(updateBfsValuesStatus)
	if err != nil {
		if errors.Is(err, e.ErrStatusDone) {
			return bfsValues, errors.WithStack(err)
		}
		return nil, errors.WithStack(err)
	}

	// 设置开始更新状态
	err = gredis.Set(fmt.Sprintf("%s-%s", UpdateBfsValuesStatusKeyPre, batchNumber), e.StatusStart, UpdateBfsValuesStatusTTL)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var firstRnID int
	if bfsValues.Cluster.Rnodes == nil || len(bfsValues.Cluster.Rnodes) == 0 {
		firstRnID = 9401
	} else {
		firstRnID = bfsValues.Cluster.Rnodes[len(bfsValues.Cluster.Rnodes)-1].ID + 1
	}

	shortServersResponse, err := GetCreateServersByBatchNumber(batchNumber)
	if err != nil {
		return nil, err
	}

	nodesMap := map[string]string{}
	for _, rn := range bfsValues.Cluster.Rnodes {
		nodesMap[rn.PublicIP] = rn.PublicIP
	}

	for _, server := range shortServersResponse.Servers {

		if nodesMap[server.Eip] == server.Eip {
			return nil, errors.WithStack(e.ErrStatusConflict)
		}

		ip := server.Ips
		// 如果是公网部署 ip 为公网 ip
		if bfsValues.Cluster.PublicNetwork {
			ip = server.Eip
		}

		rnode := &models.Rnode{
			ID:       firstRnID,
			IP:       ip,
			PublicIP: server.Eip,
		}
		bfsValues.Cluster.Rnodes = append(bfsValues.Cluster.Rnodes, rnode)
		firstRnID++
	}

	bfsHelmValuesYaml, err := gyaml.Encode(bfsValues)

	f, err := gfile.OpenFile(filePath, os.O_TRUNC|os.O_EXCL|os.O_RDWR, 0664)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer f.Close()

	_, err = f.Write(bfsHelmValuesYaml)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// 设置开始更新状态
	err = gredis.Set(fmt.Sprintf("%s-%s", UpdateBfsValuesStatusKeyPre, batchNumber), e.StatusDone, UpdateBfsValuesStatusDoneTTL)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return bfsValues, nil
}
