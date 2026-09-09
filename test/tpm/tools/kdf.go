package tools

import (
	"fmt"
	"crypto/rand"
	"golang.org/x/crypto/argon2"
)

const (
	SaltLen = 16
	KeyLen = 32
	argonTime    = 3         // number of passes
	argonMemory  = 64 * 1024 // 64 MiB — slows GPU/ASIC attacks dramatically
	argonThreads = 4
)
 
func GenerateSalt() ([]byte, error) {
	salt := make([]byte, SaltLen)
	if _, err := rand.Read(salt); err != nil {
		return nil, fmt.Errorf("generating salt: %w", err)
	}
	return salt, nil
}

func KDF(password, salt []byte) []byte {
	return argon2.IDKey(password, salt, argonTime, argonMemory, argonThreads, KeyLen)
}
