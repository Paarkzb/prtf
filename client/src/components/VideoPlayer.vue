<template>
  <div :class="available ? 'block' : 'hidden'">
    <video ref="videoPlayer" class="video-js"></video>
  </div>
  <div :class="available ? 'hidden' : 'block'">
    <div class="flex justify-center">
      <div class="w-[40%] border border-white rounded-md bg-zinc-800 min-h-[200px]">
        <div class="p-5">
          <div class="bg-white rounded-md px-2 inline-block">НЕ В СЕТИ</div>
          <p class="font-normal text-gray-700 dark:text-gray-400">
            Lorem ipsum dolor sit amet, consectetur adipisicing elit. Possimus, id.
          </p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import videojs from 'video.js'
import 'video.js/dist/video-js.css'
import 'videojs-hls-quality-selector/src/plugin'
import { defineProps, onMounted, onUnmounted, ref, useTemplateRef } from 'vue'

const props = defineProps({
  options: {
    type: Object,
    required: true
  }
})

const videoPlayer = useTemplateRef('videoPlayer')
const player = ref()
const available = ref(false)

function initializePlayer() {
  checkVideoAvailability(props.options.sources[0].src).then((isAvailable) => {
    available.value = isAvailable
    if (isAvailable) {
      player.value = videojs(videoPlayer.value, props.options, () => {
        player.value.log('onPlayerReady', this)
      })
      player.value.hlsQualitySelector()
      console.log('video player initialized.')
    } else {
      console.log('video source is not available, player not initialized.')
    }
  })
}

function checkVideoAvailability(url) {
  return window.axios
    .head(url)
    .then(() => {
      return true
    })
    .catch((error) => {
      console.error('error while checking video availability:', error)
      return false
    })
}

onMounted(() => {
  initializePlayer()
})

onUnmounted(() => {
  if (player.value) {
    player.value.dispose()
  }
})
</script>
