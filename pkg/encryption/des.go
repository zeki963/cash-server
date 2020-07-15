package encryption

import (
	"bytes"
	"crypto/des"
	"encoding/hex"
	"errors"
)

func zeroPadding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{0}, padding)
	return append(ciphertext, padtext...)
}

func zeroUnPadding(origData []byte) []byte {
	return bytes.TrimFunc(origData,
		func(r rune) bool {
			return r == rune(0)
		})
}

//DESEncrypt  DES 加密
func DESEncrypt(text string, key []byte) (string, error) {
	src := []byte(text)
	block, err := des.NewCipher(key)
	if err != nil {
		return "", err
	}
	bs := block.BlockSize()
	src = zeroPadding(src, bs)
	if len(src)%bs != 0 {
		return "", errors.New("Need a multiple of the blocksize")
	}
	out := make([]byte, len(src))
	dst := out
	for len(src) > 0 {
		block.Encrypt(dst, src[:bs])
		src = src[bs:]
		dst = dst[bs:]
	}
	return hex.EncodeToString(out), nil
}

//DESDecrypt  DES 解密
func DESDecrypt(decrypted string, key []byte) (string, error) {
	src, err := hex.DecodeString(decrypted)
	if err != nil {
		return "", err
	}
	block, err := des.NewCipher(key)
	if err != nil {
		return "", err
	}
	out := make([]byte, len(src))
	dst := out
	bs := block.BlockSize()
	if len(src)%bs != 0 {
		return "", errors.New("crypto/cipher: input not full blocks")
	}
	for len(src) > 0 {
		block.Decrypt(dst, src[:bs])
		src = src[bs:]
		dst = dst[bs:]
	}
	out = zeroUnPadding(out)
	return string(out), nil
}

// func main() {
//     key := []byte("2fa6c1e9")
//     str :="I love this beautiful world!"
//     strEncrypted, err := Encrypt(str, key)
//     if err != nil {
//         log.Fatal(err)
//     }
//     fmt.Println("Encrypted:", strEncrypted)
//     strDecrypted, err := Decrypt(strEncrypted, key)
//     if err != nil {80 67 69 103 81 67 80 114 69 87 65 61
//         log.Fatal(err)
//     }
//     fmt.Println("Decrypted:", strDecrypted)
// }
// //Output:
// //Encrypted: 5d2333b9fbbe5892379e6bcc25ffd1f3a51b6ffe4dc7af62beb28e1270d5daa1
// //Decrypted: I love this beautiful world!
