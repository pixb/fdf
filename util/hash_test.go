package util

import (
	"crypto/md5"
	"os"
	"testing"
)

func TestCalculateHash(t *testing.T) {
	// 创建一个临时文件用于测试
	tempFile, err := os.CreateTemp("", "test")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	// 写入测试内容
	testContent := "test content"
	if _, err := tempFile.WriteString(testContent); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	// 确保文件内容被刷新到磁盘
	if err := tempFile.Sync(); err != nil {
		t.Fatalf("Failed to sync temp file: %v", err)
	}
	tempFile.Close()

	// 计算哈希值
	hash, err := CalculateHash(tempFile.Name(), md5.New())
	if err != nil {
		t.Fatalf("Failed to calculate hash: %v", err)
	}

	// 打印实际哈希值，以便调试
	t.Logf("Actual hash: %s", hash)

	// 验证哈希值是否正确
	expectedHash := "9473fdd0d880a43c21b7778d34872157" // MD5 hash of "test content" (recalculated)
	if hash != expectedHash {
		t.Errorf("Expected hash %s, got %s", expectedHash, hash)
	}
}

func TestCalculateHashWithNonExistentFile(t *testing.T) {
	hasher := md5.New()
	_, err := CalculateHash("non-existent-file", hasher)
	if err == nil {
		t.Error("Expected error for non-existent file, got nil")
	}
}
