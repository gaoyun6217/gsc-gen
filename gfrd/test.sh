#!/bin/bash

# GFRD Framework Test Script
# 测试 gen 模块集成

set -e

echo "========================================"
echo "GFRD Framework - Gen Module Test"
echo "========================================"
echo ""

# 进入项目根目录
cd "$(dirname "$0")"

echo "Step 1: 检查 gen 模块..."
if [ -f "gen/main.go" ]; then
    echo "  ✓ gen 模块存在"
else
    echo "  ✗ gen 模块不存在"
    exit 1
fi

echo ""
echo "Step 2: 构建 gen 模块..."
cd gen
go build -o bin/gfrd-gen main.go 2>&1 || {
    echo "  ✗ gen 模块构建失败"
    exit 1
}
echo "  ✓ gen 模块构建成功"
cd ..

echo ""
echo "Step 3: 测试 gen CLI 帮助..."
./gen/bin/gfrd-gen --help > /dev/null 2>&1 || {
    echo "  ✗ gen CLI 执行失败"
    exit 1
}
echo "  ✓ gen CLI 执行成功"

echo ""
echo "Step 4: 检查 server 目录..."
if [ -f "server/go.mod" ]; then
    echo "  ✓ server 目录存在"
else
    echo "  ✗ server 目录不存在"
    exit 1
fi

echo ""
echo "Step 5: 检查 web 目录..."
if [ -f "web/package.json" ]; then
    echo "  ✓ web 目录存在"
else
    echo "  ✗ web 目录不存在"
    exit 1
fi

echo ""
echo "Step 6: 测试从根目录运行 gen..."
go run ./gen preview --help > /dev/null 2>&1 || {
    echo "  ✗ 从根目录运行 gen 失败"
    exit 1
}
echo "  ✓ 从根目录运行 gen 成功"

echo ""
echo "========================================"
echo "所有测试通过！✓"
echo "========================================"
echo ""
echo "使用方法:"
echo "  # 生成 CRUD 代码"
echo "  go run ./gen crud \\"
echo "    --table=\"sys_user\" \\"
echo "    --db=\"mysql:root:123456@tcp(127.0.0.1:3306)/gfrd\" \\"
echo "    --output=\"./server\" \\"
echo "    --web-output=\"./web/src\" \\"
echo "    --module=\"sys\""
echo ""
