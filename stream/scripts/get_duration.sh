#!/bin/bash
PATH=$1
CHANNEL_NAME=$2

DURATION=$(/usr/bin/ffprobe -i "/var/rec/${PATH}/${CHANNEL_NAME}.m3u8" -show_entries format=duration -v quiet -of csv=p=0)

echo -n "$DURATION"