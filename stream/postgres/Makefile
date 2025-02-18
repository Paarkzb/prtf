migrate-up:
	goose -dir ./migrations postgres postgres://postgres:postgres@localhost:5490/postgres?sslmode=disable up

create-migration:
	goose -dir .\migrations\ create [name] sql