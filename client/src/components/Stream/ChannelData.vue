<script setup lang="ts">
import { useRoute } from 'vue-router'
import { onMounted, ref } from 'vue'
import Swal from 'sweetalert2'
import { Channel, Recording } from './types'
import VideoPlayer from '@/components/VideoPlayer.vue'
import router from '@/router'
import VideoPanel from './VideoPanel.vue'
import { FwbButton, FwbAvatar } from 'flowbite-vue'

const route = useRoute()
const id = route.params.id

const channelData = ref({} as Channel)

function getChannelData() {
  window.axios
    .get(window.gatewayURL + '/stream/api/channels/' + id)
    .then((rec) => {
      console.log(rec.data)
      channelData.value = rec.data
      setVideoOptions()
    })
    .catch((error) => {
      console.log(error)
      Swal.fire({
        title: 'Ошибка',
        text: 'Неудалось получить записи',
        icon: 'error'
      })
    })
}

const recordings = ref([] as Recording[])

function getRecordings() {
  console.log('test', channelData.value)
  window.axios
    .get(window.gatewayURL + '/stream/api/channels/' + id + '/recordings')
    .then((rec) => {
      console.log(rec.data)
      recordings.value = rec.data
    })
    .catch((error) => {
      console.log(error)
      Swal.fire({
        title: 'Ошибка',
        text: 'Неудалось получить записи',
        icon: 'error'
      })
    })
}

const videoOptions = ref({})

function setVideoOptions() {
  videoOptions.value = {
    controls: true,
    liveui: true,
    preload: 'auto',
    width: 1000,
    aspectRatio: '16:9',
    playbackRates: [0.5, 0.75, 1, 1.25, 1.5, 1.75, 2],
    controlBar: {
      skipButtons: {
        forward: 5,
        backward: 5
      }
    },
    sources: [
      {
        src: `${window.gatewayURL}/stream/hls/${channelData.value.channel_name}.m3u8`,
        type: 'application/x-mpegURL'
      }
    ]
  }
}

onMounted(() => {
  getChannelData()
  getRecordings()
})
</script>

<template>
  <div class="flex justify-between">
    <div>
      <div class="flex gap-x-2 text-gray-100">
        <div>
          <fwb-avatar
            status-position="bottom-right"
            rounded
            size="lg"
            :status="channelData.live ? 'online' : 'offline'"
          />
        </div>
        <div class="flex items-center text-xl">
          {{ channelData.channel_name }}
        </div>
      </div>
    </div>
    <div class="text-white flex items-center">
      <fwb-button
        color="light"
        @click="router.push({ name: 'channelByIdSettings', params: { id: channelData.id } })"
      >
        Настройки
      </fwb-button>
    </div>
  </div>
  <div class="mt-10">
    <div v-if="channelData.id"><VideoPlayer ref="recordingVideo" :options="videoOptions" /></div>
  </div>

  <div class="mt-10">
    <h2 class="text-3xl my-2">Все видео</h2>
    <div class="grid grid-cols-3 gap-4">
      <div v-for="(rec, idx) in recordings" :key="idx">
        <VideoPanel :recording="rec" />
      </div>
    </div>
  </div>
</template>
