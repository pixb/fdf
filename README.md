# fdf - Find Duplicate Files 🗂️🔍

[![Go Version](https://img.shields.io/badge/Go-1.24+-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Build Status](https://img.shields.io/badge/Build-Passing-brightgreen.svg)](https://github.com/pixb/fdf)

A fast, efficient, and user-friendly tool to find and remove duplicate files from your system. 🚀

## ✨ Features

- **Fast Scanning**: Uses concurrent processing to quickly scan directories 📈
- **Large File Support**: Handles big files efficiently with chunked reading 📁
- **Smart Filtering**: Exclude directories and filter by file size 🎯
- **Configurable Priority**: Customize directory priority for file retention 📋
- **Dry Run Mode**: Preview changes before making them 👀
- **Progress Tracking**: Real-time progress updates during scanning 📊
- **Detailed Logs**: Clear and informative output 📝

## 📦 Installation

### From Source

```bash
# Clone the repository
git clone https://github.com/pixb/fdf.git
cd fdf

# Build the binary
make build

# Or directly with go
go build -o fdf main.go

# Install (optional)
make install
```

## 🚀 Usage

### Basic Usage

```bash
# Find duplicates in a directory
./fdf --path /path/to/search

# With custom config
./fdf --path /path/to/search --config config/fdf.example.json

# Dry run (preview only)
./fdf --path /path/to/search --dry-run
```

### Advanced Options

```bash
# Exclude directories
./fdf --path /path/to/search --exclude dir1 --exclude dir2

# Filter by file size (bytes)
./fdf --path /path/to/search --min-size 1024 --max-size 1048576

# Show version
./fdf --version
```

## ⚙️ Configuration

Create a JSON config file to set directory priorities:

```json
{
  "info": "Lower values indicate higher priority for directory retention",
  "default_priority": 99,
  "directory_priority": {
    "documents": 1,
    "photos": 2,
    "downloads": 99,
    "temp": 100
  }
}
```

## 📁 Project Structure

```
fdf/
├── config/          # Configuration files
├── scripts/         # Test and utility scripts
├── util/            # Utility functions
├── main.go          # Main application
├── Makefile         # Build instructions
└── README.md        # This file
```

## 🧪 Testing

Run the test suite to ensure everything works correctly:

```bash
# Generate test files
./scripts/generate_test_files.sh

# Run tests
./scripts/test_fdf_run.sh

# Run unit tests
go test -v ./...
```

## 🔧 Development

### Prerequisites
- Go 1.24+
- Make (optional)

### Build Commands

```bash
# Build the binary
make build

# Install the binary
make install

# Clean build artifacts
make clean

# Run tests
make test
```

## 🤝 Contributing

Contributions are welcome! Here's how you can help:

1. **Fork the repository** 🍴
2. **Create a feature branch** 🌿
3. **Make your changes** ✏️
4. **Run tests** ✅
5. **Submit a pull request** 🚀

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Cobra](https://github.com/spf13/cobra) - Command-line interface library
- [Viper](https://github.com/spf13/viper) - Configuration management
- [Go](https://golang.org/) - The programming language

## 📞 Contact

- GitHub: [@pixb](https://github.com/pixb)
- Project: [https://github.com/pixb/fdf](https://github.com/pixb/fdf)

---

Made with ❤️ by the fdf team
