// Package encryutil 加解密模块
package encryutil

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"log"

	"GoCampusLogin/utils"
)

// CBCEncrypt AES/CBC/PKCS7Padding 加密
func CBCEncrypt(text []byte, key []byte) (string, error) {
	// iv 为随机的 16 位字符串
	iv := []byte(utils.RandomNString(aes.BlockSize))

	block, err := aes.NewCipher(key)
	if err != nil {
		log.Fatalf("AES-CBC encrypt fail: %v", err)
		return "", err
	}
	// 填充
	paddingText := PKCS7Padding(text, block.BlockSize())

	// 构建加密模式
	blockMode := cipher.NewCBCEncrypter(block, iv)
	// 加密
	result := make([]byte, len(paddingText))
	blockMode.CryptBlocks(result, paddingText)
	// 返回使用 Base64 加密后的密文
	return base64.StdEncoding.EncodeToString(result), nil
}

func PKCS7Padding(text []byte, blockSize int) []byte {
	// 计算待填充的长度
	padding := blockSize - len(text)%blockSize
	var paddingText []byte
	if padding == 0 {
		// 已对齐，填充一整块数据，每个数据为 blockSize
		paddingText = bytes.Repeat([]byte{byte(blockSize)}, blockSize)
	} else {
		// 未对齐 填充 padding 个数据，每个数据为 padding
		paddingText = bytes.Repeat([]byte{byte(padding)}, padding)
	}
	return append(text, paddingText...)
}

// CBCDecrypt decrypts the data using AES.
func CBCDecrypt(cryptoText string, key []byte) (string, error) {
	block, _ := aes.NewCipher(key)
	// base64 解密，从 ciphertext 中分离出 iv 与 密文
	ciphertext, _ := base64.StdEncoding.DecodeString(cryptoText)
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]
	// 解密并返回 plainText
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(ciphertext, ciphertext)
	ciphertext, _ = PKCSUnPadding(ciphertext)
	return string(ciphertext[48:]), nil
}

// PKCSUnPadding removes padding from the data.
func PKCSUnPadding(origData []byte) ([]byte, error) {
	length := len(origData)
	unPadding := int(origData[length-1])
	if unPadding > length {
		return nil, errors.New("unPadding error")
	}
	return origData[:(length - unPadding)], nil
}
