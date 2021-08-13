package utils

import (
	"encoding/base64"
)

func Base64RawURLEncodeing(d string) string {
	return base64.RawURLEncoding.EncodeToString([]byte(d))
}

func Base64RawURLDecodeing(d string) (string, error) {
	b, err := base64.RawURLEncoding.DecodeString(d)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
func Base64DecodeStringToString(d string) (string, error) {
	b, err := base64.StdEncoding.DecodeString(d)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

