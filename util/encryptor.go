package util

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
)

func HashMD5(input string) string {
	hashedBytes := md5.Sum([]byte(input))
	return hex.EncodeToString(hashedBytes[:])
}

func HashSHA256(input string) string {
	hashedBytes := sha256.Sum256([]byte(input))
	return hex.EncodeToString(hashedBytes[:])
}
