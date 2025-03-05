#!/bin/bash
CHANNEL_NAME=$1
DATE=$(date '+%Y-%m-%d_%H:%M:%S')

mkdir /var/rec/${CHANNEL_NAME}/${DATE} -p
cp -r /var/hls/${CHANNEL_NAME}* /var/rec/${CHANNEL_NAME}/${DATE}/
rm -r /var/hls/${CHANNEL_NAME}*

echo -e "#EXT-X-PLAYLIST-TYPE:VOD\n#EXT-X-ENDLIST" >> /var/rec/${CHANNEL_NAME}/${DATE}/${CHANNEL_NAME}_360p/index.m3u8 
echo -e "#EXT-X-PLAYLIST-TYPE:VOD\n#EXT-X-ENDLIST" >> /var/rec/${CHANNEL_NAME}/${DATE}/${CHANNEL_NAME}_480p/index.m3u8 
echo -e "#EXT-X-PLAYLIST-TYPE:VOD\n#EXT-X-ENDLIST" >> /var/rec/${CHANNEL_NAME}/${DATE}/${CHANNEL_NAME}_720p/index.m3u8 

ffmpeg -i "/var/rec/${CHANNEL_NAME}/${DATE}/${CHANNEL_NAME}.m3u8" -ss 00:00:01 -vframes 1 -q:v 2 "/var/rec/${CHANNEL_NAME}/${DATE}/poster.jpg"

echo -n "${CHANNEL_NAME}/${DATE}"

