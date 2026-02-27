.PHONY: build install clean test help

# 构建可执行文件
build:
	go build -o bin/gfrd-gen main.go

# 安装到 GOPATH
install:
	go install .

# 运行测试
test:
	go test -v ./...

# 清理构建文件
clean:
	rm -rf bin/
	rm -rf vendor/

# 格式化代码
fmt:
	go fmt ./...

# 代码检查
lint:
	go vet ./...

# 整理依赖
tidy:
	go mod tidy

# 生成示例代码（测试用）
demo:
	./bin/gfrd-gen crud \
		--table="sys_user" \
		--db="mysql:root:123456@tcp(127.0.0.1:3306)/gfrd" \
		--output="./testdata/server" \
		--web-output="./testdata/web" \
		--module="sys" \
		--features="add,edit,delete,view,list" \
		--with-test

# 预览生成结果
preview:
	./bin/gfrd-gen preview \
		--table="sys_user" \
		--db="mysql:root:123456@tcp(127.0.0.1:3306)/gfrd" \
		--module="sys"

# 显示帮助
help:
	@echo "GFRD Code Generator - Makefile"
	@echo ""
	@echo "Usage:"
	@echo "  make build     - 构建可执行文件"
	@echo "  make install   - 安装到 GOPATH"
	@echo "  make test      - 运行测试"
	@echo "  make clean     - 清理构建文件"
	@echo "  make fmt       - 格式化代码"
	@echo "  make lint      - 代码检查"
	@echo "  make tidy      - 整理依赖"
	@echo "  make demo      - 生成示例代码"
	@echo "  make preview   - 预览生成结果"
	@echo ""
