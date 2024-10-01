migrate-up:
	goose -dir .\migrations postgres "host=localhost port=5435 user=postgres password=postgres dbname=postgres sslmode=disable" up
migrate-down:
	goose -dir .\migrations postgres "host=localhost port=5435 user=postgres password=postgres dbname=postgres sslmode=disable" down

test-migrate-up:
	goose -dir .\tests\migrations postgres "host=localhost port=5435 user=postgres password=postgres dbname=postgres sslmode=disable" up

test:
	go test ./tests -count=1 -v