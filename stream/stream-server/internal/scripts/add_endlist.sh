#!/bin/bash
CHANNEL_NAME=$1

echo -e "#EXT-X-PLAYLIST-TYPE:VOD\n#EXT-X-ENDLIST" >> /var/hls/${CHANNEL_NAME}_360p/index.m3u8
echo -e "#EXT-X-PLAYLIST-TYPE:VOD\n#EXT-X-ENDLIST" >> /var/hls/${CHANNEL_NAME}_480p/index.m3u8
echo -e "#EXT-X-PLAYLIST-TYPE:VOD\n#EXT-X-ENDLIST" >> /var/hls/${CHANNEL_NAME}_720p/index.m3u8
