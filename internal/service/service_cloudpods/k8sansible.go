package service_cloudpods

import (
	"bytes"
	"cloud-batch/configs"
	"cloud-batch/internal/pkg/db/gredis"
	"cloud-batch/internal/pkg/e"
	"cloud-batch/internal/pkg/logging"
	"fmt"
	"github.com/gogf/gf/os/gfile"
	"github.com/pkg/errors"
	"github.com/relex/aini"
	"golang.org/x/net/context"
	"gopkg.in/redis.v5"
	"os"
	"os/exec"
	"strings"
	"time"
)

const (
	UpdateK8sAnsibleHostsStatusKeyPre = "updateK8sAnsibleHostsStatus"
	UpdateK8sAnsibleHostsTTL          = time.Second * 5
	UpdateK8sAnsibleHostsDoneTTL      = time.Hour * 24
	AddK8sNodeStatusKey               = "addK8sNodeStatus"
	AddK8sNodeTTL                     = time.Minute * 30
	AddK8sNodeDoneTTL                 = time.Hour * 24
)

func UpdateK8sAnsibleByServers(batchNumber string, k8sNodeSuf string) (*aini.InventoryData, error) {

	// 获取主机创建状态
	serverCreateServersStatus, err := gredis.Get(fmt.Sprintf("%s-%s", ServerCreateServersStatus, batchNumber)).Result()
	if err != nil {
		// 如果redis 返回的不是 nil 错误，返回
		if !errors.Is(err, redis.Nil) {
			return nil, errors.WithStack(err)
		}
	}

	// 不是完成状态就是创建中，不处理中间状态
	if serverCreateServersStatus != e.StatusDone {
		return nil, e.ErrServerCreating
	}

	// cpoy初始 hosts.init
	filePath := fmt.Sprintf("%s/%s", configs.K8sBfs.AnsibleConfigPath, configs.K8sBfs.AnsibleHostsName)
	if !gfile.Exists(fmt.Sprintf("%s.init", filePath)) {
		err = gfile.CopyFile(filePath, fmt.Sprintf("%s.init", filePath))
		if err != nil {
			return nil, errors.WithStack(err)
		}
	}

	// 备份
	err = gfile.CopyFile(filePath, fmt.Sprintf("%s.cloud-batch-bak", filePath))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	ansibleHosts, err := aini.ParseFile(filePath)
	if err != nil {
		return nil, err
	}

	if len(ansibleHosts.Groups["new_nodes"].Hosts) > 0 {
		return nil, e.ErrAnsibleHostsNewNodesNotEmpty
	}

	updateK8sAnsibleHostsStatus, err := gredis.Get(fmt.Sprintf("%s-%s", UpdateK8sAnsibleHostsStatusKeyPre, batchNumber)).Result()
	if err != nil {
		// 如果redis 返回的不是 nil 错误，返回
		if !errors.Is(err, redis.Nil) {
			return nil, errors.WithStack(err)
		}
	}

	if updateK8sAnsibleHostsStatus == e.StatusDone {
		return ansibleHosts, nil
	}

	// 设置开始更新状态
	err = gredis.Set(fmt.Sprintf("%s-%s", UpdateK8sAnsibleHostsStatusKeyPre, batchNumber), e.StatusStart, UpdateK8sAnsibleHostsTTL)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	shortServersResponse, err := GetCreateServersByBatchNumber(batchNumber)
	if err != nil {
		return nil, err
	}

	for _, server := range shortServersResponse.Servers {

		ip := server.Ips
		// 如果是公网部署 ip 为公网 ip
		if ansibleHosts.Groups["k8sCluster"].Vars["public_network_node"] == "True" {
			ip = server.Eip
		}

		hostName := fmt.Sprintf("node%s.%s", strings.ReplaceAll(ip, ".", "-"), k8sNodeSuf)
		if _, ok := ansibleHosts.Hosts[hostName]; ok {
			return nil, errors.WithStack(e.ErrStatusConflict)
		}

		host := &aini.Host{
			Name: hostName,
			Vars: map[string]string{},
		}

		host.Vars["ansible_host"] = ip
		host.Vars["advertise_address"] = server.Eip

		ansibleHosts.Groups["new_nodes"].Hosts[hostName] = host
		ansibleHosts.Hosts[hostName] = host
	}

	f, err := gfile.OpenFile(filePath, os.O_APPEND|os.O_EXCL|os.O_RDWR, 0664)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer f.Close()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	for _, h := range ansibleHosts.Groups["new_nodes"].Hosts {
		host := fmt.Sprintf("%s ansible_host=%s advertise_address=%s\n", h.Name, h.Vars["ansible_host"], h.Vars["advertise_address"])
		_, err = f.WriteString(host)
	}

	// 设置开始更新状态
	err = gredis.Set(fmt.Sprintf("%s-%s", UpdateK8sAnsibleHostsStatusKeyPre, batchNumber), e.StatusDone, UpdateK8sAnsibleHostsDoneTTL)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return ansibleHosts, nil
}

func K8sAddNode(batchNumber string) error {
	addK8sNodeStatus, err := gredis.Get(AddK8sNodeStatusKey).Result()
	if err != nil {
		return errors.WithStack(err)
	}

	if addK8sNodeStatus == e.StatusStart {
		return e.ErrStatusStart
	}

	if addK8sNodeStatus == e.StatusDone {
		return nil
	}

	// 判断 ansible hosts 修改状态
	updateK8sAnsibleHostsStatus, err := gredis.Get(fmt.Sprintf("%s-%s", UpdateK8sAnsibleHostsStatusKeyPre, batchNumber)).Result()
	if err != nil {
		// 如果redis 返回的不是 nil 错误，返回
		if !errors.Is(err, redis.Nil) {
			return errors.WithStack(err)
		}
	}

	if updateK8sAnsibleHostsStatus != e.StatusDone {
		return errors.Errorf("batchNumber %s ansible hosts 更新未完成", batchNumber)
	}

	err = gredis.Set(AddK8sNodeStatusKey, e.StatusStart, AddK8sNodeTTL)
	if err != nil {
		logging.Logger.Errorf("addK8sNodes set AddK8sNodeStatusKey err: %v", err)
		return nil
	}
	go addK8sNodes()

	return nil
}

func addK8sNodes() {
	// 超时机制
	ctxt, cancel := context.WithTimeout(context.Background(), AddK8sNodeTTL)
	defer cancel()

	cmd := exec.CommandContext(ctxt, "%s%s", configs.K8sBfs.AnsibleConfigPath, configs.K8sBfs.K8sAddNodeScript)
	var buf bytes.Buffer
	cmd.Stderr = &buf

	if err := cmd.Start(); err != nil {
		if err := gredis.Set(AddK8sNodeStatusKey, e.StatusError, AddK8sNodeTTL); err != nil {
			logging.Logger.Errorf("Set AddK8sNodeStatus err: %+v", err)
		}
		logging.Logger.Errorf("addK8sNodes start err: %+v", err)
	}

	if err := cmd.Wait(); err != nil {
		if err := gredis.Set(AddK8sNodeStatusKey, e.StatusError, AddK8sNodeTTL); err != nil {
			logging.Logger.Errorf("Set AddK8sNodeStatus err: %+v", err)
		}
		logging.Logger.Errorf("addK8sNodes wait err: %+v", err)
	}

	if err := gredis.Set(AddK8sNodeStatusKey, e.StatusError, AddK8sNodeDoneTTL); err != nil {
		logging.Logger.Errorf("Set AddK8sNodeStatus err: %+v", err)
	}

	if err := gredis.Del(AddK8sNodeStatusKey); err != nil {
		logging.Logger.Errorf("Del AddK8sNodeStatus err: %+v", err)
	}
}
