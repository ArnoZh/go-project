/**
 * Author : kaizhongsun
 * CreateDate : 2021/1/6 17:09
 * Software : GoLand
 * Remark : 用于数据库密码加解密
 **/
package util

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
)

const key = `qXtwrpwLo7x2$aBw*n3yyUgL6il0%Th7`
const iv = `O%pHZ2*5RqAR26ny`

//todo 密码使用%%隔开
//"root:%%加密后的password%%@tcp(192.168.0.181:3306)/enterserver_db?tls=skip-verify&autocommit=true"
func MysqlDecryptPassword(source string) string {
	//密码解密
	strs := strings.Split(source, `%%`)
	if len(strs) < 3 {
		panic(errors.New(`[ERROR] This mysql password is not encrypt...`).Error())
		return source
	}
	password := strs[1]
	dp := DecryptPassword(password)
	return strings.ReplaceAll(source, `%%`+password+`%%`, dp)
}

//todo 密码使用%%隔开
//redis密码解密
func RedisDecryptPassword(password string) string {
	if !strings.HasPrefix(password, `%%`) {
		println(errors.New(`[ERROR] This redis password is not encrypt...`).Error())
		return password
	}
	strs := strings.Split(password, `%%`)
	if len(strs) < 3 {
		panic("[ERROR] Wrong password format...")
	}
	return DecryptPassword(strs[1])
}

//密码加密
func EncryptPassword(src string) string {
	return AesEncrypt(src, key)
}

//密码解密
func DecryptPassword(secret string) string {
	return AesDecrypt(secret, key)
}

//AES加密
func AesEncrypt(src, key string) string {
	plaintext := []byte(src)
	keyByte := md5.Sum([]byte(key))
	block, err := aes.NewCipher(keyByte[:])
	if err != nil {
		panic("invalid decrypt key")
	}
	blockSize := block.BlockSize()
	plaintext = PKCS5Padding(plaintext, blockSize)
	iv := []byte(iv)
	blockMode := cipher.NewCBCEncrypter(block, iv)

	ciphertext := make([]byte, len(plaintext))
	blockMode.CryptBlocks(ciphertext, plaintext)

	return fmt.Sprintf("%X", ciphertext)
}

//AES解密
func AesDecrypt(secret, key string) string {
	ciphertext, _ := hex.DecodeString(secret)
	keyByte := md5.Sum([]byte(key))
	block, err := aes.NewCipher(keyByte[:])
	if err != nil {
		panic("invalid decrypt key")
	}

	blockSize := block.BlockSize()
	if len(ciphertext) < blockSize {
		panic("ciphertext too short")
	}

	iv := []byte(iv)
	if len(ciphertext)%blockSize != 0 {
		panic("ciphertext is not a multiple of the block size")
	}

	blockModel := cipher.NewCBCDecrypter(block, iv)

	plaintext := make([]byte, len(ciphertext))
	blockModel.CryptBlocks(plaintext, ciphertext)
	plaintext = PKCS5UnPadding(plaintext)

	return string(plaintext)
}

//明文补码算法
func PKCS5Padding(src []byte, blockSize int) []byte {
	padding := blockSize - len(src)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padtext...)
}

//明文减码算法
func PKCS5UnPadding(src []byte) []byte {
	length := len(src)
	unpadding := int(src[length-1])
	return src[:(length - unpadding)]
}
