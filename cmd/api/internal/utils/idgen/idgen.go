package idgen

import (
	"crypto/rand"
	"io"
	"strings"
)

// GenerateRandomID generates a cryptographically random ID string of the specified length and charset
func GenerateRandomID(length int, charset string) string {
	if length == 0 {
		length = DefaultIDSize
	}

	if len(charset) == 0 {
		charset = URLSafeAlphanumericCharset
	}

	charSetBytes := []byte(charset)
	charSetLen := len(charSetBytes)

	randomBytes := make([]byte, length)
	io.ReadFull(rand.Reader, randomBytes)

	idBuilder := strings.Builder{}
	for _, b := range randomBytes {
		idBuilder.WriteByte(charSetBytes[int(b)%charSetLen])
	}

	return idBuilder.String()
}
