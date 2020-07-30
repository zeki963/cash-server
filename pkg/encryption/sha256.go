package encryption

import (
	"crypto/sha256"
	"encoding/hex"
)

//Sha256encode  SHA-256
func Sha256encode(text string) string {
	h := sha256.New()
	h.Write([]byte(text))
	return hex.EncodeToString(h.Sum(nil))
}
