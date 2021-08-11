package v1

import (
	"cloud-batch/internal/models"
	"cloud-batch/internal/pkg/app"
	"cloud-batch/internal/pkg/e"
	"cloud-batch/internal/pkg/logging"
	"cloud-batch/internal/service/servie_cloudpods"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	"net/url"
)

// BatchCreateServers
// @Summary Batch Create Servers
// @Tags Servers
// @Produce  json
// @Security ApiKeyAuth
// @Param batchCreateServers body models.BatchCreateServersForm true "batchCreateServers"
// @Success 200 {object} app.ResponseString
// @Failure 500 {object} app.ResponseString
// @Failure default {object} app.ResponseString
// @Router /batch/servers [post]
func BatchCreateServers(c *gin.Context) {
	var (
		appG              = app.Gin{C: c}
		createServersForm = models.BatchCreateServersForm{}
	)

	err := c.BindJSON(&createServersForm)
	if err != nil {
		code := app.VerdictJsonErr(err)
		appG.ResponseI18nMsgSimple(code, err, nil)
		return
	}
	resp, cloudErr, err := servie_cloudpods.BatchCreateServers(createServersForm)
	if err != nil {
		if cloudErr != nil {
			logging.Logger.Errorf("BatchCreateServers cloudErr: %v", cloudErr)
		}
		appG.ResponseI18nMsgSimple(e.InternalError, err, nil)
		return
	}

	appG.ResponseI18nMsgSimple(e.StatusOK, nil, resp)
}

// BatchGetServers
// @Summary Batch Get Servers
// @Tags Servers
// @Produce  json
// @Param project query string false "project" default(bfs)
// @Param provider query string true "provider" default(aliyun)
// @Param status query []string false "status"
// @Param batch_number query string false "batch_number"
// @Success 200 {object} app.ResponseString
// @Failure 500 {object} app.ResponseString
// @Failure default {object} app.ResponseString
// @Router /servers [get]
func BatchGetServers(c *gin.Context) {
	var (
		appG  = app.Gin{C: c}
		valid = validator.New()
	)

	project := c.Query("project")
	provider := c.Query("provider")
	status := c.QueryArray("status")

	batchNumber := c.Query("batch_number")

	getServersForm := models.BatchGetServersForm{
		Project:     project,
		Provider:    provider,
		BatchNumber: batchNumber,
		Status:      status,
	}

	err := valid.Struct(getServersForm)
	if err != nil {
		appG.ResponseI18nMsgSimple(e.InvalidParameter, err, nil)
		return
	}

	urlValues := url.Values{}

	// 结果返回字段, 因为不能指定 metadata 字段返回，弃用
	//urlValues.Add("field", "id")
	//urlValues.Add("field", "name")
	//urlValues.Add("field", "eip")
	//urlValues.Add("field", "ips")
	//urlValues.Add("field", "zone_ext_id")
	//urlValues.Add("field", "metadata")
	//urlValues.Add("field", "provider")
	//urlValues.Add("field", "status")

	// 查询状态数组
	if status != nil && len(status) > 0 {
		for _, item := range status {
			urlValues.Add("status", item)
		}
	}

	// 查询过滤字段
	urlValues.Set("provider", provider)
	urlValues.Set("tags.0.key", "user:project")
	if project != "" {
		urlValues.Set("tags.0.value", project)
	}
	if batchNumber != "" {
		urlValues.Set("tags.1.key", "user:batchNumber")
		urlValues.Set("tags.1.value", batchNumber)
	}

	resp, _, err := servie_cloudpods.ListServers(nil, urlValues)
	if err != nil {
		appG.ResponseI18nMsgSimple(e.InternalError, err, nil)
		return
	}

	shortServersResponse := new(models.ShortServersResponse)
	err = json.Unmarshal(resp, shortServersResponse)
	if err != nil {
		appG.ResponseI18nMsgSimple(e.InternalError, err, nil)
	}

	appG.ResponseI18nMsgSimple(e.StatusOK, nil, shortServersResponse)
}

// BatchDeleteServers
// @Summary Batch Get Servers
// @Tags Servers
// @Produce  json
// @Security ApiKeyAuth
// @Param batchDeleteServersForm body models.BatchDeleteServersForm true "batchDeleteServersForm"
// @Success 200 {object} app.ResponseString
// @Failure 500 {object} app.ResponseString
// @Failure default {object} app.ResponseString
// @Router /batch/servers [delete]
func BatchDeleteServers(c *gin.Context) {
	var (
		appG              = app.Gin{C: c}
		deleteServersForm = models.BatchDeleteServersForm{}
	)

	err := c.BindJSON(&deleteServersForm)
	if err != nil {
		code := app.VerdictJsonErr(err)
		appG.ResponseI18nMsgSimple(code, err, nil)
		return
	}

	if (deleteServersForm.IDs == nil || len(deleteServersForm.IDs) == 0) && (deleteServersForm.Provider == "" || deleteServersForm.Project == "") {
		appG.ResponseI18nMsgSimple(e.InvalidParameter, errors.New("如果IDs 为空，则必须填写 Provider 和 project 参数"), nil)
		return
	}

	serverCount, doneCount, errIDs, err := servie_cloudpods.BatchDeleteServers(deleteServersForm)
	if err != nil {
		appG.ResponseI18nMsgSimple(e.InternalError, err, nil)
		return
	}

	if serverCount == 0 {
		appG.ResponseI18nMsgSimple(e.NotFound, errors.New("serverCount is 0"),nil)
		return
	}

	appG.ResponseI18nMsgSimple(e.StatusOK, nil, gin.H{"server_count": serverCount, "done_count": doneCount, "errIds": errIDs})
}
