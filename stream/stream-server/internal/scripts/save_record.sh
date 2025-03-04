#!/bin/bash
CHANNEL_NAME=$1
DATE=$(date '+%Y-%m-%d_%H:%M:%S')

echo -e "#EXT-X-PLAYLIST-TYPE:VOD\n#EXT-X-ENDLIST" >> /var/hls/${CHANNEL_NAME}_360p/index.m3u8
echo -e "#EXT-X-PLAYLIST-TYPE:VOD\n#EXT-X-ENDLIST" >> /var/hls/${CHANNEL_NAME}_480p/index.m3u8
echo -e "#EXT-X-PLAYLIST-TYPE:VOD\n#EXT-X-ENDLIST" >> /var/hls/${CHANNEL_NAME}_720p/index.m3u8

mkdir /var/rec/${CHANNEL_NAME}/${DATE} -p
mv /var/hls/${CHANNEL_NAME}* /var/rec/${CHANNEL_NAME}/${DATE}/

echo "${CHANNEL_NAME}/${DATE}/${CHANNEL_NAME}.m3u8"

