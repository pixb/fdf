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

	// 重置哈希对象，以便重复使用
	hasher.Reset()

	// 分块读取文件，避免一次性加载大文件到内存
	buffer := make([]byte, 64*1024) // 64KB 缓冲区
	for {
		n, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			return "", fmt.Errorf("Read file error: %w", err)
		}
		if n == 0 {
			break
		}
		if _, err := hasher.Write(buffer[:n]); err != nil {
			return "", fmt.Errorf("Write to hasher error: %w", err)
		}
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}
