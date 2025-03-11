<script setup lang="ts">
import { ref, watch, type PropType } from 'vue'
import { Channel } from '../types'
import VideoPlayer from '@/components/VideoPlayer.vue'

const props = defineProps({
  channel: {
    type: Object as PropType<Channel>,
    required: true
  }
})

const videoOptions = ref({})

function setVideoOptions() {
  videoOptions.value = {
    controls: true,
    liveui: true,
    preload: 'auto',
    width: 1000,
    aspectRatio: '16:9',
    playbackRates: [0.5, 0.75, 1, 1.25, 1.5, 1.75, 2],
    sources: [
      {
        src: `${window.gatewayURL}/stream/hls/${props.channel.channel_name}.m3u8`,
        type: 'application/x-mpegURL'
      }
    ]
  }
}

watch(
  () => props.channel,
  (newChannel: Channel) => {
    if (newChannel && newChannel.id) {
      setVideoOptions()
    }
  },
  { immediate: true }
)
</script>

<template>
  <div v-if="channel.live"><VideoPlayer :options="videoOptions" /></div>
  <div v-else>
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
