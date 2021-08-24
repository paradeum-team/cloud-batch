package service_cloudpods

import (
	"cloud-batch/internal/pkg/e"
	"fmt"
	"testing"
)

func TestUpdateK8sAnsibleByServers(t *testing.T) {
	ansibleHosts, err := UpdateK8sAnsibleByServers("20210823163332.102805", "solarfs.k8s")
	if err != nil {
		code := e.GetCodeByErr(err)
		if code.GetCode() != e.StatusOK.GetCode() {
			t.Errorf("%+v", err)
		}
	}
	fmt.Println(ansibleHosts)
}
