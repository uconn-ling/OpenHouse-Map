package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

func HashString(ipa string) string {
	h := sha256.New()
	h.Write([]byte(ipa))
	hash := hex.EncodeToString(h.Sum(nil))
	return hash
}
