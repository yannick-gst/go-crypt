package main

import (
	"crypto/sha256"

	"golang.org/x/crypto/pbkdf2"
)

const (
	chunkSize = 64 * 1024 // 64 KB chunks,
	saltSize = 32
)

func generateKey(password string, salt []byte) []byte {
	return pbkdf2.Key([]byte(password), salt, 4096, 32, sha256.New)
}
