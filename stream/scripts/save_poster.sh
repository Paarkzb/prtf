#!/bin/bash
PATH=$1
CHANNEL_NAME=$2

/usr/bin/ffmpeg -i "/var/rec/${PATH}/${CHANNEL_NAME}.m3u8" -ss 00:00:01 -vframes 1 -q:v 2 "/var/rec/${PATH}/poster.jpg"
/usr/bin/magick /var/rec/${PATH}/poster.jpg -resize 290x160 -quality 80 /var/rec/${PATH}/poster.jpg

echo -n "${PATH}/poster.jpg"
