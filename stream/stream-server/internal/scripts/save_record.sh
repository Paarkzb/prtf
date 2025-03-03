#!/bin/bash
CHANNEL_NAME=$1
DATE=$(date '+%Y-%m-%d_%H:%M:%S')

mkdir /var/rec/${CHANNEL_NAME}/${DATE} -p
mv /var/hls/${CHANNEL_NAME}* /var/rec/${CHANNEL_NAME}/${DATE}/

echo "${CHANNEL_NAME}/${DATE}/${CHANNEL_NAME}.m3u8"

