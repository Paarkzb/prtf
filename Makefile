deploy:
	docker cp ./auth/sso/config/ prtf_auth_server:/app
	docker cp ./quiz/configs/ prtf_quiz_server:/app
	docker cp ./quiz/.env prtf_quiz_server:/app

	goose -dir ./auth/sso/migrations postgres "host=localhost port=5435 user=postgres password=postgres dbname=postgres sslmode=disable" up
	cd quiz && soda migrate up && cd ..

migrate:
	goose -dir ./auth/sso/migrations postgres "host=localhost port=5435 user=postgres password=postgres dbname=postgres sslmode=disable" up
	cd quiz && soda migrate up && cd ..

run:
	docker-compose up -d --build

d-up-one:
	docker-compose build prtf-stream-server && docker-compose up -d --no-deps prtf-stream-server
	
backup:
	docker run --rm -v stream-platform_vod:/data -v $(pwd):/backup alpine \
    tar czf /backup/vod-$(date +%Y%m%d).tar.gz -C /data .