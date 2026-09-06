APP_NAME := server
BUILD_DIR := build
MAIN := main.go

# Windows
BINARY_WIN := $(BUILD_DIR)\$(APP_NAME).exe

# Linux
BINARY_LINUX := $(BUILD_DIR)/$(APP_NAME)

.PHONY: swag run build clean air linux

# 生成 Swagger 文档
swag:
	swag init

# 本地运行
run: swag
	go run $(MAIN)

# 热更新
air: swag
	air



linux:
	set CGO_ENABLED=0&& set GOOS=linux&& set GOARCH=amd64&& go build -o ../build/server/server main.go

# 清理
clean:
	if exist $(BUILD_DIR) rmdir /s /q $(BUILD_DIR)
	mkdir $(BUILD_DIR)