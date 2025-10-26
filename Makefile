.PHONY: build run

build:
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/main main.go

# TOKEN=xxx make run
run:
	@if [ -z "$$TOKEN" ]; then \
		echo "❌ 错误: 请设置 TOKEN 环境变量"; \
		echo "用法: TOKEN=your_token make run"; \
		echo "或者: export \$$(cat .env | xargs) && make run"; \
		exit 1; \
	fi
	TOKEN=$$TOKEN go run main.go
