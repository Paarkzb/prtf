# go build -o ./dist/serverV2.exe ./cmd | .\dist\serverV2.exe

build:
	@echo building server
	go build -o ./dist/serverV2.exe ./cmd
	@echo server built

start: build
	@echo starting server
	.\dist\serverV2.exe
	@echo server started