package encryption

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
)

// func main() {
// 	//定義密鑰，必須是24byte
// 	key := []byte("123456789012345678901234")
// 	//定義明文
// 	origData := []byte("hello world")

// 	//加密
// 	en := ThriDESEnCrypt(origData, key)
// 	//解密
// 	de := ThriDESDeCrypt(en, key)
// 	fmt.Println(string(de))
// }

//ThriDESDeCrypt 解密
func ThriDESDeCrypt(crypted, key []byte) []byte {
	//獲取block塊
	block, _ := des.NewTripleDESCipher(key)
	//創建切片
	context := make([]byte, len(crypted))
	//設置解密方式
	blockMode := cipher.NewCBCDecrypter(block, key[:8])
	//解密密文到數組
	blockMode.CryptBlocks(context, crypted)
	//去補碼
	context = PKCSUnPadding(context)
	return context
}

//PKCSUnPadding 去補碼
func PKCSUnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:length-unpadding]
}

//ThriDESEnCrypt 加密
func ThriDESEnCrypt(origData, key []byte) []byte {
	//獲取block塊
	block, _ := des.NewTripleDESCipher(key)
	//補碼
	origData = PKCSPadding(origData, block.BlockSize())
	//設置加密方式為 3DES  使用3條56位的密鑰對數據進行三次加密
	blockMode := cipher.NewCBCEncrypter(block, key[:8])

	//創建明文長度的數組
	crypted := make([]byte, len(origData))

	//加密明文
	blockMode.CryptBlocks(crypted, origData)

	return crypted

}

//PKCSPadding 補碼
func PKCSPadding(origData []byte, blockSize int) []byte {
	//計算需要補幾位數
	padding := blockSize - len(origData)%blockSize
	//在切片後面追加char數量的byte(char)
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)

	return append(origData, padtext...)
}
