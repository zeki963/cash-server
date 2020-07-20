package encryption

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

const (
	sKey        = "asaismyhero66666"
	ivParameter = "cqig7777cqig7777"
)

// //DEMO
// O := "我的雞雞很大"
// fmt.Println("原文：", O)
// A := encryption.PswEncrypt(O, "asaismyhero66666")
// fmt.Println("密文：", A)
// B := encryption.PswDecrypt(A, "asaismyhero66666")
// fmt.Println("解密結果：", B)

//PswEncrypt 加密
func PswEncrypt(src string, myKey string) string {
	key := []byte(myKey)
	iv := []byte(ivParameter)
	result, err := Aes128Encrypt([]byte(src), key, iv)
	if err != nil {
		panic(err)
	}
	return base64.RawStdEncoding.EncodeToString(result)
}

//PswDecrypt 解密
func PswDecrypt(src string, myKey string) string {
	key := []byte(myKey)
	iv := []byte(ivParameter)
	var result []byte
	var err error
	result, err = base64.RawStdEncoding.DecodeString(src)
	if err != nil {
		panic(err)
	}
	origData, err := Aes128Decrypt(result, key, iv)
	if err != nil {
		panic(err)
	}
	return string(origData)
}

//Aes128Encrypt Aes128Encrypt
func Aes128Encrypt(origData, key []byte, IV []byte) ([]byte, error) {
	if key == nil || len(key) != 16 {
		return nil, nil
	}
	if IV != nil && len(IV) != 16 {
		return nil, nil
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	origData = PKCS5Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, IV[:blockSize])
	crypted := make([]byte, len(origData))
	// 根据CryptBlocks方法的说明，如下方式初始化crypted也可以
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

//Aes128Decrypt Aes128Decrypt
func Aes128Decrypt(crypted, key []byte, IV []byte) ([]byte, error) {
	if key == nil || len(key) != 16 {
		return nil, nil
	}
	if IV != nil && len(IV) != 16 {
		return nil, nil
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, IV[:blockSize])
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS5UnPadding(origData)
	return origData, nil
}

//PKCS5Padding PKCS5Padding
func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

//PKCS5UnPadding PKCS5UnPadding
func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	// 去掉最后一个字节 unpadding 次
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
