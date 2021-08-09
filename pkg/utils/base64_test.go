package utils

import (
	"fmt"
	"testing"
)

func TestBase64(t *testing.T) {
	enStr := Base64RawURLEncodeing("cb94d68b-3d78-44ac-a5cf-185b94b/b9c51+")
	fmt.Println(enStr)
	deStr, _ := Base64RawURLDecodeing(enStr)
	fmt.Println(deStr)
}
