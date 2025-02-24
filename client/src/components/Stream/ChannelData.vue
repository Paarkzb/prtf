<script setup lang="ts">
import { useRoute } from 'vue-router'
import { onMounted, ref, watchEffect } from 'vue'
import Swal from 'sweetalert2'
import Hls, { Level } from 'hls.js'
import { Channel } from './types'

const route = useRoute()

const id = route.params.id

const channelData = ref({} as Channel)

function getChannelData() {
  window.axios
    .get(window.gatewayURL + '/stream/api/channels/' + id)
    .then((rec) => {
      console.log(rec.data)
      channelData.value = rec.data
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

const streamToken = ref('')

function startStream() {
  window.axios
    .post(window.gatewayURL + '/stream/api/streams/start')
    .then((resp) => {
      console.log(resp.data)
      streamToken.value = resp.data
    })
    .catch((error) => {
      console.log(error)
      Swal.fire({
        title: 'Ошибка',
        text: 'Неудалось запустить стрим',
        icon: 'error'
      })
    })
}

const hls = new Hls({
  enableWorker: true,
  autoStartLoad: true,
  capLevelToPlayerSize: true
})

const qualityLevels = ref([] as Level[])
const videoQuality = ref()

const video = ref({} as HTMLMediaElement)

function playStream() {
  const masterPlaylistUrl = `${window.gatewayURL}/stream/hls/${channelData.value.channel_name}.m3u8`

  if (Hls.isSupported()) {
    hls.attachMedia(video.value)
    hls.loadSource(masterPlaylistUrl)

    hls.on(Hls.Events.MANIFEST_PARSED, () => {
      qualityLevels.value = hls.levels

      video.value.play()
    })
  }
}

watchEffect(() => {
  if (videoQuality.value) {
    hls.currentLevel = videoQuality.value
  }
})

onMounted(() => {
  getChannelData()
})
</script>

<template>
  <div>
    <h1 class="text-6xl">{{ channelData.channel_name }}</h1>
    <h2 class="text-6xl">{{ channelData.live ? 'Онлайн' : 'Оффлайн' }}</h2>
    <h2 class="text-6xl">{{ channelData.rf_active_stream_id }}</h2>
  </div>
  <div>
    <button @click="startStream()">Запустить стрим</button>
    <div class="p-2 bg-white">{{ streamToken }}</div>
  </div>
  <div>
    <h1>Стрим</h1>
    <div>
      <button @click="playStream()">Смотреть стрим {{ channelData.channel_name }}</button>
    </div>
    <div id="qualitySelector" v-if="!!qualityLevels.length">
      <select
        v-model="videoQuality"
        @change="
          () => {
            hls.currentLevel = parseInt(videoQuality)
          }
        "
      >
        <option v-for="(level, index) in qualityLevels" :key="index" :value="index">
          {{ level.height }}p
        </option>
      </select>
    </div>
    <video ref="video" controls></video>
  </div>
</template>
