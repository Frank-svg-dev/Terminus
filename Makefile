# ==============================================================================
# 🛠️ Terminus Build System
# ==============================================================================

# --- 变量定义 ---
BINARY_NAME=terminus-enforcer
EXPORTER_BIN_NAME=terminus-exporter
CMD_PATH=./cmd/terminus-enforcer
EXPORTER_PATH=./cmd/terminus-exporter
BIN_DIR=./bin
DOCKER_IMAGE=terminus-enforcer
VERSION?=v0.1.0

# 获取 Git Commit Hash 和 构建时间 (用于注入版本信息)
GIT_COMMIT=$(shell git rev-parse --short HEAD || echo "unknown")
BUILD_TIME=$(shell date "+%F %T")

# 编译参数：
# -s -w: 去掉调试符号，减小二进制体积
# -X: 注入变量 (如果你在代码里定义了 Version 变量)
LDFLAGS=-ldflags "-s -w -X 'main.Version=${VERSION}' -X 'main.GitCommit=${GIT_COMMIT}' -X 'main.BuildTime=${BUILD_TIME}'"

# --- 默认任务 ---
.PHONY: all
all: build

# ==============================================================================
# 📦 编译与构建 (Build)
# ==============================================================================

.PHONY: build
build: ## 编译当前平台的二进制文件
	@echo "🚀 Building ${BINARY_NAME}..."
	@mkdir -p ${BIN_DIR}
	go build ${LDFLAGS} -o ${BIN_DIR}/${BINARY_NAME} ${CMD_PATH}
	@echo "✅ Build success: ${BIN_DIR}/${BINARY_NAME}"

.PHONY: build-exporter
build-exporter: ## 交叉编译 Linux (amd64) 版本 (适合生产环境/Docker)
	@echo "🐧 Building Linux amd64  exporter static binary..."
	@mkdir -p ${BIN_DIR}
	# 关键参数解释：
	# CGO_ENABLED=0 : 禁用 CGO，强制使用 Go 原生实现（如 DNS 解析），切断对 libc 的依赖
	# -a            : 强制重新编译所有包
	# -ldflags      : -s -w 去掉符号表减小体积，-extldflags "-static" 确保完全静态
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags '-s -w -extldflags "-static"' -o ${BIN_DIR}/${EXPORTER_BIN_NAME}-linux ${EXPORTER_PATH}
	@echo "✅ Linux binary ready: ${BIN_DIR}/${EXPORTER_BIN_NAME}-linux"


.PHONY: build-linux
build-linux: ## 交叉编译 Linux (amd64) 版本 (适合生产环境/Docker)
	@echo "🐧 Building Linux amd64  static binary..."
	@mkdir -p ${BIN_DIR}
	# 关键参数解释：
	# CGO_ENABLED=0 : 禁用 CGO，强制使用 Go 原生实现（如 DNS 解析），切断对 libc 的依赖
	# -a            : 强制重新编译所有包
	# -ldflags      : -s -w 去掉符号表减小体积，-extldflags "-static" 确保完全静态
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags '-s -w -extldflags "-static"' -o ${BIN_DIR}/${BINARY_NAME}-linux ${CMD_PATH}
	@echo "✅ Linux binary ready: ${BIN_DIR}/${BINARY_NAME}-linux"

.PHONY: run
run: build ## 编译并在本机运行 (需要 sudo 权限连 socket)
	@echo "🏃 Running ${BINARY_NAME}..."
	sudo ${BIN_DIR}/${BINARY_NAME} --v=2

# ==============================================================================
# 🧹 代码质量与清理 (Quality & Clean)
# ==============================================================================

.PHONY: clean
clean: ## 清理构建产物
	@echo "🧹 Cleaning up..."
	@rm -rf ${BIN_DIR}
	@echo "✅ Done."

.PHONY: fmt
fmt: ## 格式化代码 (go fmt)
	@go fmt ./...

.PHONY: vet
vet: ## 静态检查 (go vet)
	@go vet ./...

.PHONY: lint
lint: ## 运行 golangci-lint (需要先安装)
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "⚠️ golangci-lint not installed. Skipping."; \
	fi

.PHONY: test
test: ## 运行单元测试
	@go test -v ./...

# ==============================================================================
# 🐳 Docker 相关
# ==============================================================================

.PHONY: docker
docker: build-linux ## 构建 Docker 镜像
	@echo "🐳 Building Docker image: ${DOCKER_IMAGE}:${VERSION}"
	docker build -t ${DOCKER_IMAGE}:${VERSION} .

# ==============================================================================
# ❓ 帮助信息
# ==============================================================================

.PHONY: help
help: ## 显示帮助信息
	@echo "Terminus Makefile Commands:"
	@awk 'BEGIN {FS = ":.*##"; printf "\033[36m\033[0m"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)