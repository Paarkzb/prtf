migrate-up:
	goose -dir .\migrations postgres "host=localhost port=5454 user=postgres password=postgres dbname=postgres sslmode=disable" up
migrate-down:
	goose -dir .\migrations postgres "host=localhost port=5454 user=postgres password=postgres dbname=postgres sslmode=disable" down