package encode

import (
	"crypto/md5"
	"fmt"
)

//Md5 加密
func Md5(str string) string {
	data := []byte(str)
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has)
	fmt.Println(md5str)
	return md5str
}
