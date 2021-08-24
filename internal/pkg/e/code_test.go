package e

import (
	"fmt"
	"testing"
)

func TestGetCodeByErr(t *testing.T) {
	err := GetCodeByErr(ErrServerCreating)

	fmt.Printf("%+v", err)
}
