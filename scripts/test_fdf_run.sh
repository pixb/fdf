#!/usr/bin/env bash

cd "$(dirname "$0")/../"

# 测试基本重复文件
echo "Testing basic duplicate files..."
./fdf --path test-path/scenario1 --config config/fdf.example.json

# 测试大文件
echo "Testing large files..."
./fdf --path test-path/scenario2 --config config/fdf.example.json

# 测试多重复文件
echo "Testing multiple duplicate files..."
./fdf --path test-path/scenario3 --config config/fdf.example.json

# 测试排除目录
echo "Testing excluding directories..."
./fdf --path test-path --exclude scenario2 --config config/fdf.example.json

# 测试文件大小过滤
echo "Testing file size filtering..."
./fdf --path test-path --min-size 1024 --max-size 1048576 --config config/fdf.example.json

echo "All tests completed successfully!"