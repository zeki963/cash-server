package encryption

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

// func main() {
// 	orig := "hello world"
// 	key := "123456781234567812345678"
// 	fmt.Println("原文：", orig)

// 	encryptCode := AesEncrypt(orig, key)
// 	fmt.Println("密文：", encryptCode)

// 	decryptCode := AesDecrypt(encryptCode, key)
// 	fmt.Println("解密結果：", decryptCode)
// }

//AesEncrypt  AES加密
func AesEncrypt(orig string, key string) string {
	// 轉成字節數組
	origData := []byte(orig)
	k := []byte(key)

	// 分組秘鑰
	block, _ := aes.NewCipher(k)
	// 獲取秘鑰塊的長度
	blockSize := block.BlockSize()
	// 補全碼
	origData = PKCS7Padding(origData, blockSize)
	// 加密模式
	blockMode := cipher.NewCBCEncrypter(block, k[:blockSize])
	// 創建數組
	cryted := make([]byte, len(origData))
	// 加密
	blockMode.CryptBlocks(cryted, origData)

	return base64.StdEncoding.EncodeToString(cryted)

}

//AesDecrypt AES解密
func AesDecrypt(cryted string, key string) string {
	// 轉成字節數組
	crytedByte, _ := base64.StdEncoding.DecodeString(cryted)
	k := []byte(key)

	// 分組秘鑰
	block, _ := aes.NewCipher(k)
	// 獲取秘鑰塊的長度
	blockSize := block.BlockSize()
	// 加密模式
	blockMode := cipher.NewCBCDecrypter(block, k[:blockSize])
	// 創建數組
	orig := make([]byte, len(crytedByte))
	// 解密
	blockMode.CryptBlocks(orig, crytedByte)
	// 去補全碼
	orig = PKCS7UnPadding(orig)
	return string(orig)
}

//PKCS7Padding 補碼
func PKCS7Padding(ciphertext []byte, blocksize int) []byte {
	padding := blocksize - len(ciphertext)%blocksize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// PKCS7UnPadding 去碼
func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
