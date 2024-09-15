# 目标目录
OUTPUT_DIR=output

# 创建输出目录（如果不存在）
$(OUTPUT_DIR):
	mkdir -p $(OUTPUT_DIR)

# 默认任务，编译所有平台
all: linux mac windows

# 编译适用于 Linux 的可执行文件
linux:
	GOOS=linux GOARCH=amd64 go build -o $(OUTPUT_DIR)/ssh2a_linux

# 编译适用于 macOS 的可执行文件
mac:
	GOOS=darwin GOARCH=amd64 go build -o $(OUTPUT_DIR)/ssh2a_mac

# 编译适用于 Windows 的可执行文件
windows:
	GOOS=windows GOARCH=amd64 go build -o $(OUTPUT_DIR)/ssh2a_windows.exe

# 清理输出目录
clean:
	rm -rf $(OUTPUT_DIR)
