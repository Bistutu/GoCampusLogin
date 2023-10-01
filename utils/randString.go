package utils

import (
	"math/rand"
	"strings"
)

const (
	baseString       = "ABCDEFGHJKMNPQRSTWXYZabcdefhijkmnprstwxyz0123456789" // 用于生成随机字符串
	baseStringLength = len(baseString)
)

// RandomNString 随机生成 n 位字符串
func RandomNString(n int) string {
	builder := strings.Builder{}
	for i := 0; i < n; i++ {
		builder.WriteString(string(baseString[rand.Intn(baseStringLength)]))
	}
	return builder.String()
}
