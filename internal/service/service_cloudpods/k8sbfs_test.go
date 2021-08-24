package service_cloudpods

import (
	"fmt"
	"testing"
)

func TestAddK8sNodes(t *testing.T) {
	bfsValues, err := BfsUpdateValuesByServers("20210820163855.141268")
	if err != nil {
		t.Errorf("%+v", err)
	}
	fmt.Println(bfsValues)
}
