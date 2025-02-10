d-up-one:
	docker-compose build go-server && docker-compose up -d --no-deps go-server
backup:
	docker run --rm -v stream-platform_vod:/data -v $(pwd):/backup alpine \
    tar czf /backup/vod-$(date +%Y%m%d).tar.gz -C /data .