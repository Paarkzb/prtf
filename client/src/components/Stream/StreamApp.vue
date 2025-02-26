<script setup lang="ts">
import { onMounted, ref, watchEffect } from 'vue'
import Hls, { Level } from 'hls.js'
import router from '@/router'
import Swal from 'sweetalert2'
import { useChannelStore } from '@/stores/store'
import { Channel, Recording } from './types'

// Чат
const ws = new WebSocket('ws://prtf.localhost:8090/stream/chat/ws')

const chatMessages = ref([] as ChatMessage[])
const chatInput = ref('')

declare interface ChatMessage {
  stream_id: string
  text: string
  time: number
  username: string
}

ws.onmessage = (event) => {
  const msg = JSON.parse(event.data)

  chatMessages.value.push(msg)
}

function sendMessage() {
  const text = chatInput.value.trim()
  if (text) {
    const msg = {
      stream_id: '',
      username: 'Guest',
      text: text
    }
    ws.send(JSON.stringify(msg))
    chatInput.value = ''
  }
}

const hls = new Hls({
  enableWorker: true,
  autoStartLoad: true,
  capLevelToPlayerSize: true
})

// Стрим

const activeStreams = ref({})
const recordings = ref([] as Recording[])
const qualityLevels = ref([] as Level[])
const videoQuality = ref()

function getActiveStreams() {
  window.axios
    .get(window.gatewayURL + '/stream/api/streams')
    .then((streams) => {
      activeStreams.value = Object.keys(streams.data)
    })
    .catch((error) => {
      console.log(error)
      Swal.fire({
        title: 'Ошибка',
        text: 'Неудалось получить список активных стримов',
        icon: 'error'
      })
    })
}

const video = ref({} as HTMLMediaElement)

function playStream(streamKey: string) {
  const masterPlaylistUrl = `${window.gatewayURL}/stream/hls/${streamKey}.m3u8`

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

function getRecordings() {
  window.axios
    .get(window.gatewayURL + '/stream/api/recordings')
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

const recordingVideo = ref({} as HTMLMediaElement)
function playRecording(path: string) {
  const recordingUrl = `${window.gatewayURL}/stream/vod/${path}`
  recordingVideo.value.src = recordingUrl
  recordingVideo.value.load()
  recordingVideo.value.play()
}

// Канал

const channelStore = useChannelStore()

function getMyChannel() {
  window.axios
    .get(window.gatewayURL + '/stream/api/channels/user')
    .then((resp) => {
      console.log(resp.data)
      channelStore.login(resp.data)
    })
    .catch((error) => {
      console.log(error)
      Swal.fire({
        title: 'Ошибка',
        text: 'Неудалось получить данные канала',
        icon: 'error'
      })
    })
}

function saveChannel() {
  window.axios
    .post(window.gatewayURL + '/stream/api/channels')
    .then((resp) => {
      console.log(resp.data)
    })
    .catch((error) => {
      console.log(error)
      Swal.fire({
        title: 'Ошибка',
        text: 'Неудалось создать канал',
        icon: 'error'
      })
    })
}

const channels = ref([] as Channel[])

function getChannels() {
  window.axios
    .get(window.gatewayURL + '/stream/api/channels')
    .then((rec) => {
      console.log(rec.data)
      channels.value = rec.data
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

onMounted(() => {
  getMyChannel()
  getActiveStreams()
  getRecordings()
  getChannels()
})
</script>

<template>
  <button @click="saveChannel()" class="p-2 border bg-blue-300">Создать канал</button>
  <button
    @click="router.push({ name: 'channelById', params: { id: channelStore.channel.id } })"
    class="p-2 border bg-blue-300"
  >
    Мой канал
  </button>

  <h1 class="text-4xl text-center">Список каналов</h1>
  <div>
    <div v-for="(channel, idx) in channels" :key="idx">
      <button
        class="p-2 bg-red-300 border"
        @click="router.push({ name: 'channelById', params: { id: channel.id } })"
      >
        {{ channel.channel_name }} {{ channel.live ? 'Онлайн' : 'Оффлайн' }}
      </button>
    </div>
  </div>

  <!-- <h1>Active Streams</h1>
  <div>
    <div v-for="stream in activeStreams" :key="stream">
      <button @click="playStream(stream)">{{ stream }}</button>
    </div>
  </div>
  <div id="qualitySelector" v-if="!!qualityLevels.length">
    {{ videoQuality }}
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

  <div>
    <div v-for="(rec, idx) in recordings" :key="idx">
      {{ rec.name }} <button @click="playRecording(rec.path)">play</button>
    </div>
  </div>
  <video ref="recordingVideo" controls width="640"></video>

  <div class="chat-box">
    <div>
      <div v-for="(msg, idx) in chatMessages" :key="idx">
        <b>{{ msg.username }}</b
        >: {{ msg.text }}
        <small>{{ new Date(msg.time * 1000).toLocaleTimeString() }}</small>
      </div>
    </div>
    <input
      type="text"
      v-model="chatInput"
      v-on:keypress="
        (e) => {
          if (e.key === 'Enter') sendMessage()
        }
      "
      placeholder="Напишите сообщение ..."
    />
    <button @click="sendMessage()">Отправить</button>
  </div> -->
</template>
