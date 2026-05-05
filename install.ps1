#!/usr/bin/env pwsh

# 检查当前目录是否包含必要文件
if (!(Test-Path "main.go") -or !(Test-Path "go.mod")) {
    Write-Host "错误: 请在fdf项目根目录中运行此脚本" -ForegroundColor Red
    exit 1
}

# 定义变量
$BINARY_NAME = "fdf"
$VERSION = "0.0.2"

Write-Host "正在安装 $BINARY_NAME 版本 $VERSION..." -ForegroundColor Green

# 设置环境变量
$env:CGO_ENABLED = "0"

# 执行安装命令
try {
    go install -ldflags="-X=main.Version=$VERSION"
    if ($LASTEXITCODE -eq 0) {
        Write-Host "安装成功!" -ForegroundColor Green
        Write-Host ""
        Write-Host "验证安装:"
        Write-Host "  fdf --version"
        Write-Host ""
        Write-Host "使用示例:"
        Write-Host "  fdf --path C:\path\to\search"
        Write-Host "  fdf --path C:\path\to\search --config config\fdf.yaml"
    } else {
        Write-Host "安装失败，请检查错误信息" -ForegroundColor Red
        exit $LASTEXITCODE
    }
} catch {
    Write-Host "安装过程中发生错误: $($_.Exception.Message)" -ForegroundColor Red
    exit 1
} finally {
    # 清理环境变量
    Remove-Item Env:CGO_ENABLED -ErrorAction SilentlyContinue
}
