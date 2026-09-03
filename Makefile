.PHONY: build-s build-s-linux build-w

build-s:
	@echo Building backend server for Linux...
	@if not exist build\server mkdir build\server
	cd server && set GOOS=linux&& set GOARCH=amd64&& set CGO_ENABLED=0&& go build -o ..\build\server\server main.go
	@echo Backend server built successfully: build\server\server
	@echo Backend server built successfully: build\server\server.exe


build-w:
	@echo Building frontend web...
	@if not exist build\web mkdir build\web
	cd web && pnpm install && pnpm build
	@if exist build\web (rmdir /s /q build\web)
	xcopy /e /i /y web\dist build\web
	@echo Frontend web built successfully: build\web
