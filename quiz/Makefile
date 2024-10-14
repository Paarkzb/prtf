# go build -o ./dist/serverV2.exe ./cmd | .\dist\serverV2.exe

build:
	@echo building server
	go build -o ./dist/serverV2.exe ./cmd
	@echo server built

start: build
	@echo starting server
	.\dist\serverV2.exe
	@echo server started

migrate-up:
	goose -dir .\migrations postgres "host=localhost port=5436 user=postgres password=postgres dbname=postgres sslmode=disable" up
migrate-down:
	goose -dir .\migrations postgres "host=localhost port=5436 user=postgres password=postgres dbname=postgres sslmode=disable" down