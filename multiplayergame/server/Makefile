build: 
	@echo "building server"
	@go build -o ./dist/server.exe ./cmd/
	@echo "server built"

start: build
	@echo "starting server"
	cmd /c start /B ./dist/server.exe
	@echo "server running"

stop:
	@echo "stopping server"
	@taskkill /IM server.exe /F
	@echo "server stopped"