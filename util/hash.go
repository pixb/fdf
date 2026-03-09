package util

import (
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"os"
)

func CalculateHash(filePath string, hasher hash.Hash) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("Open file error: %w", err)
	}
	defer file.Close()

	if _, err := io.Copy(hasher, file); err != nil {
		return "", fmt.Errorf("Read file error: %w", err)
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}
