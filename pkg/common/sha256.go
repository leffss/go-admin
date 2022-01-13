package common

import (
	"crypto/sha256"
	"encoding/hex"
)

func EncodeSHA256(value string) string {
	hash := sha256.New()
	hash.Write([]byte(value))

	return hex.EncodeToString(hash.Sum(nil))
}
