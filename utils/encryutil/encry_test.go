package encryutil

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"GoCampusLogin/utils"
)

func TestEncry(t *testing.T) {
	key := []byte(utils.RandomNString(16)) // 生成 16 位随机密钥
	raw := "123456"
	text := utils.RandomNString(64) + raw // 64 随机字符 + 明文
	// 加密
	cipherText, err := CBCEncrypt([]byte(text), key)
	assert.NoError(t, err)
	// 解密
	plainText, err := CBCDecrypt(cipherText, key)
	assert.NoError(t, err)
	assert.Equal(t, raw, plainText)
}

// 加密性能测试
func BenchmarkEncry(b *testing.B) {
	key := []byte(utils.RandomNString(16)) // 生成 16 位随机密钥
	text := []byte("123456")
	for i := 0; i < b.N; i++ {
		en, err := CBCEncrypt(text, key)
		assert.NoError(b, err)
		assert.NotEmpty(b, en)
	}
}
