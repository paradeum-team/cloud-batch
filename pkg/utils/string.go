package utils

import (
	"crypto/rand"
	"fmt"
	"hash/crc32"
	"math/big"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

func GetBetweenStr(str, end string) string {
	n := 7
	str = string([]byte(str)[n:])
	m := strings.LastIndex(str, end)
	if m == -1 {
		m = len(str)
	}
	str = string([]byte(str)[:m])
	return str
}

func ConvertRelativePath(relativePath string) string {
	relativePath = filepath.Dir(fmt.Sprintf("/%s", relativePath))
	if relativePath != "" && relativePath != "/" {
		return relativePath
	}
	return ""
}

func ReplaceSpecialChar(str string) string {
	if str == "" {
		return str
	}
	replacer := strings.NewReplacer(
		"/", "-",
		":", "-",
	)
	return replacer.Replace(str)
}

func StringToInt64(numStr string) (int64, error) {
	return strconv.ParseInt(numStr, 10, 64)
}

func StringMustToInt64(numStr string) int64 {
	num, _ := strconv.ParseInt(numStr, 10, 64)
	return num
}

func StringToInt(numStr string) (int, error) {
	return strconv.Atoi(numStr)
}

func StringMustToInt(numStr string) int {
	num, _ := strconv.Atoi(numStr)
	return num
}

func IntToString(num int) string {
	return strconv.Itoa(num)
}

func Int64ToString(num int64) string {
	return strconv.FormatInt(num, 10)
}

func SplitWords(name string) string {
	var gatherRegexp = regexp.MustCompile("([^A-Z]+|[A-Z]+[^A-Z]+|[A-Z]+)")
	var acronymRegexp = regexp.MustCompile("([A-Z]+)([A-Z][^A-Z]+)")

	words := gatherRegexp.FindAllStringSubmatch(name, -1)
	if len(words) > 0 {
		var name []string
		for _, words := range words {
			if m := acronymRegexp.FindStringSubmatch(words[0]); len(m) == 3 {
				name = append(name, m[1], m[2])
			} else {
				name = append(name, words[0])
			}
		}

		return strings.Join(name, "_")
	}
	return ""
}

// GenerateRandomBytes returns securely generated random bytes.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
// reference: https://gist.github.com/dopey/c69559607800d2f2f90b1b1ed4e550fb
func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

// GenerateRandomString returns a securely generated random string.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
// reference: https://gist.github.com/dopey/c69559607800d2f2f90b1b1ed4e550fb
func GenerateRandomString(n int, letters string) (string, error) {
	if letters == "" {
		letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"
	}

	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		ret[i] = letters[num.Int64()]
	}

	return string(ret), nil
}

// 生成随机小写字符及数字
func GenerateRandomStringLower(n int) (string, error) {
	return GenerateRandomString(n, "0123456789abcdefghijklmnopqrstuvwxyz")
}

func StringToHashCode(str string) int {
	v := int(crc32.ChecksumIEEE([]byte(str)))
	if v >= 0 {
		return v
	}
	if -v >= 0 {
		return -v
	}
	// v == MinInt
	return 0
}
