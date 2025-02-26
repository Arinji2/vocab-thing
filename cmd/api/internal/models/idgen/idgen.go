package idgen

import (
	"crypto/rand"
	"fmt"
	"io"
	"strings"
)

// GenerateRandomID generates a cryptographically random ID string of the specified length and charset
func GenerateRandomID(length int, charset string) (string, error) {
	if length <= 0 {
		return "", fmt.Errorf("length must be positive")
	}

	if length == 0 {
		length = DefaultIDSize
	}

	if len(charset) == 0 {
		charset = URLSafeAlphanumericCharset
	}

	charSetBytes := []byte(charset)
	charSetLen := len(charSetBytes)

	randomBytes := make([]byte, length)
	_, err := io.ReadFull(rand.Reader, randomBytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	idBuilder := strings.Builder{}
	for _, b := range randomBytes {
		idBuilder.WriteByte(charSetBytes[int(b)%charSetLen])
	}

	return idBuilder.String(), nil
}
