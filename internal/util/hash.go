package util

import (
	"crypto/sha256"
	"encoding/hex"
)

func HashIP(ip string) string {
  hash := sha256.Sum256([]byte(ip))
  return hex.EncodeToString(hash[:])
}