# 目标目录
OUTPUT_DIR=output

# 创建输出目录（如果不存在）
$(OUTPUT_DIR):
	mkdir -p $(OUTPUT_DIR)

# 默认任务
all: web linux

# 构建前端
web:
	cd web && pnpm install && pnpm run build

# 编译适用于 Linux 的可执行文件
linux: web
	GOOS=linux GOARCH=amd64 go build -o $(OUTPUT_DIR)/ssh2a_linux

# 编译适用于 macOS 的可执行文件
mac: web
	GOOS=darwin GOARCH=amd64 go build -o $(OUTPUT_DIR)/ssh2a_mac

# 编译适用于 Windows 的可执行文件
windows: web
	GOOS=windows GOARCH=amd64 go build -o $(OUTPUT_DIR)/ssh2a_windows.exe

# 清理输出目录
clean:
	rm -rf $(OUTPUT_DIR)
	rm -rf web/dist
