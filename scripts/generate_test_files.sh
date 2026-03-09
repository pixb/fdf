#!/bin/bash

set -e
cd "$(dirname "$0")/../"

# 创建测试目录结构
mkdir -p test-path/scenario1/dir1 test-path/scenario1/dir2
mkdir -p test-path/scenario2/dir1 test-path/scenario2/dir2
mkdir -p test-path/scenario3/dir1 test-path/scenario3/dir2 test-path/scenario3/dir3

# 生成基本重复文件
echo "Hello, World!" > test-path/scenario1/dir1/file1.txt
echo "Different content" > test-path/scenario1/dir1/file2.txt
cp test-path/scenario1/dir1/file1.txt test-path/scenario1/dir2/file1.txt
echo "Another file" > test-path/scenario1/dir2/file3.txt

# 生成大文件（100MB）
dd if=/dev/zero of=test-path/scenario2/dir1/large_file.bin bs=1M count=100
cp test-path/scenario2/dir1/large_file.bin test-path/scenario2/dir2/large_file.bin

# 生成多重复文件
echo "Common content" > test-path/scenario3/dir1/file1.txt
echo "Another common content" > test-path/scenario3/dir1/file2.txt
echo "Unique content" > test-path/scenario3/dir1/file3.txt
cp test-path/scenario3/dir1/file1.txt test-path/scenario3/dir2/file1.txt
cp test-path/scenario3/dir1/file2.txt test-path/scenario3/dir2/file2.txt
cp test-path/scenario3/dir1/file1.txt test-path/scenario3/dir3/file1.txt

echo "Test files generated successfully!"