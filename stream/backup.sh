#!/bin/bash
DATE=$(date +%Y%m%d)
mysqldump -h postgres -u stream -psecret streams > /backups/db-$DATE.sql
aws s3 sync /var/vod s3://backup-bucket/vod-$DATE