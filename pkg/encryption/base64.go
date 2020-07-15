package encryption

import "encoding/base64"

const (
	base64Table = "IJjkKLMNO567PQX12RVW3YZaDEFGbcdefghiABCHlSTUmnopqrxyz04stuvw89+/"
)

var coder = base64.NewEncoding(base64Table)

//Base64Encode  加密
func Base64Encode(src []byte) []byte {
	return []byte(coder.EncodeToString(src))
}

//Base64Decode  解密
func Base64Decode(src []byte) ([]byte, error) {
	return coder.DecodeString(string(src))
}
