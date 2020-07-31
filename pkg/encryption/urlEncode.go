package encryption

import "net/url"

//Urlencode URL 函式編碼
func Urlencode(text string) string {
	return url.QueryEscape(text)
}

//Urldecrypt  URL 函式解碼
func Urldecrypt(text string) (string, error) {
	return url.QueryUnescape(text)
}