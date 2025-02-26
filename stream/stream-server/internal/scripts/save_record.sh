#!/bin/bash
CHANNEL_NAME=$1
DATE=$(date '+%Y-%m-%d_%H:%M:%S')

mkdir /var/rec/${CHANNEL_NAME}/${DATE} -p
mv /var/hls/${CHANNEL_NAME}* /var/rec/${CHANNEL_NAME}/${DATE}/

echo "${CHANNEL_NAME}/${DATE}/${CHANNEL_NAME}.m3u8"

# #!/bin/bash
# SOURCE_DIR="/var/hls"
# DEST_DIR="/var/rec"
# while true; do
#   rsync -avur --delete "${SOURCE_DIR}/" "${DEST_DIR}/"
#   sleep 5
# done
