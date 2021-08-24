package v1

import (
	"cloud-batch/internal/pkg/app"
	"cloud-batch/internal/pkg/e"
	"cloud-batch/internal/service/service_cloudpods"
	"fmt"
	"github.com/gin-gonic/gin"
)

// BfsUpdateValues
// @Summary Update Bfs Values
// @Tags Bfs
// @Produce  json
// @Security ApiKeyAuth
// @Param batch_number query string true "batchNumber"
// @Success 200 {object} app.ResponseString
// @Failure 500 {object} app.ResponseString
// @Failure default {object} app.ResponseString
// @Router /batch/bfsUpdateValues [post]
func BfsUpdateValues(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
	)

	batchNumber := c.Query("batch_number")
	if batchNumber == "" {
		appG.ResponseI18nMsg(e.InvalidParameterValue,
			fmt.Errorf("batch_number cannot be empty"),
			map[string]interface{}{".Name": "batch_number"},
			nil,
			nil,
		)
		return
	}

	bfsValues, err := service_cloudpods.BfsUpdateValuesByServers(batchNumber)
	if err != nil {
		code := e.GetCodeByErr(err)
		appG.ResponseI18nMsgSimple(code, err, bfsValues)
		return
	}

	appG.ResponseI18nMsgSimple(e.StatusOK, nil, bfsValues)
}

// K8sAnsibleUpdateHosts
// @Summary K8s Ansible Update Hosts
// @Tags K8s
// @Produce  json
// @Security ApiKeyAuth
// @Param batch_number query string true "batchNumber"
// @Param k8s_node_suf query string true "k8s node suffix"
// @Success 200 {object} app.ResponseString
// @Failure 500 {object} app.ResponseString
// @Failure default {object} app.ResponseString
// @Router /batch/k8sAnsibleUpdateHosts [post]
func K8sAnsibleUpdateHosts(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
	)

	batchNumber := c.Query("batch_number")
	if batchNumber == "" {
		appG.ResponseI18nMsg(e.InvalidParameterValue,
			fmt.Errorf("batch_number cannot be empty"),
			map[string]interface{}{".Name": "batch_number"},
			nil,
			nil,
		)
		return
	}

	k8sNodeSuf := c.Query("k8s_node_suf")
	if k8sNodeSuf == "" {
		appG.ResponseI18nMsg(e.InvalidParameterValue,
			fmt.Errorf("k8s_node_suf cannot be empty"),
			map[string]interface{}{".Name": "k8s_node_suf"},
			nil,
			nil,
		)
		return
	}

	ansibleHosts, err := service_cloudpods.UpdateK8sAnsibleByServers(batchNumber, k8sNodeSuf)
	if err != nil {
		code := e.GetCodeByErr(err)
		appG.ResponseI18nMsgSimple(code, err, fmt.Sprintf("%v", ansibleHosts))
		return
	}

	appG.ResponseI18nMsgSimple(e.StatusOK, nil, fmt.Sprintf("%v", ansibleHosts))
}

// K8sAddNodeByBatchNumber
// @Summary K8s Ansible Update Hosts
// @Tags K8s
// @Produce  json
// @Security ApiKeyAuth
// @Param batch_number query string true "batchNumber"
// @Success 200 {object} app.ResponseString
// @Failure 500 {object} app.ResponseString
// @Failure default {object} app.ResponseString
// @Router /batch/k8sAddNodeByBatchNumber [post]
func K8sAddNodeByBatchNumber(c *gin.Context) {

	var (
		appG = app.Gin{C: c}
	)

	batchNumber := c.Query("batch_number")
	if batchNumber == "" {
		appG.ResponseI18nMsg(e.InvalidParameterValue,
			fmt.Errorf("batch_number cannot be empty"),
			map[string]interface{}{".Name": "batch_number"},
			nil,
			nil,
		)
		return
	}
	//err := service_cloudpods.K8sAddNodeByBatchNumber(batchNumber)

}
