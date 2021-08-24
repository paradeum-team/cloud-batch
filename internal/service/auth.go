package service

import (
	"cloud-batch/api/middleware"
	"cloud-batch/internal/models"
	"cloud-batch/internal/pkg/db/gredis"
	"cloud-batch/internal/pkg/e"
	"cloud-batch/internal/pkg/logging"
	"fmt"
	"github.com/pkg/errors"
)

func init() {
	value, err := gredis.Get("user-admin").Result()
	if err != nil {
		logging.Logger.Warnf("没有admin用户，需要初始化admin: %v", err)
		// 初始化admin用户
		err := gredis.Set("user-admin", "admin", -1)
		if err != nil {
			logging.Logger.Fatalf("初始化admin用户失败：%v", err)
		}
	} else {
		logging.Logger.Info("已经存在admin用户")
		if value == "admin" {
			logging.Logger.Warnf("admin用户密码为默认密码，请尽快修改")
		}
	}
}

func Login(auth models.Auth) (string, e.ErrorCode, error) {
	code, err := CheckAuth(auth.Username, auth.Password)
	if err != nil {
		return "", code, err
	}
	token, err := middleware.GenerateToken(auth.Username, auth.Password)
	if err != nil {
		return "", e.InternalError, err
	}
	return token, e.StatusOK, nil
}

func CheckAuth(username, password string) (e.ErrorCode, error) {
	value, err := gredis.Get(fmt.Sprintf("user-%s", username)).Result()
	if err != nil {
		return e.AuthError, err
	}

	if value != password {
		return e.AuthError, errors.New("password error")
	}

	return e.StatusOK, nil
}

func UpdateAuth(updateAuth models.UpdateAuth) (e.ErrorCode, error) {
	code, err := CheckAuth(updateAuth.Username, updateAuth.OldPassword)
	if err != nil {
		return code, err
	}

	err = gredis.Set(fmt.Sprintf("user-%s", updateAuth.Username), updateAuth.Password, -1)
	if err != nil {
		return e.InternalError, err
	}

	// 使用旧密码重新生成 token 覆盖旧 token, 间接使旧 token 失效
	middleware.GenerateToken(updateAuth.Username, updateAuth.OldPassword)

	return e.StatusOK, nil
}
