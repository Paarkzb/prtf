<script setup lang="ts">
import { onMounted, ref, watchEffect } from 'vue'
import Hls, { Level } from 'hls.js'
import Swal from 'sweetalert2'

const chatServer = 'prtf.localhost:8090/stream/chat'
const streamServer = 'http://prtf.localhost:8090/stream'
const apiServer = 'http://prtf.localhost:8090/stream/api'

// Чат
const ws = new WebSocket('ws://' + chatServer + 'ws')

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
    .get(apiServer + '/streams')
    .then((streams) => {
      activeStreams.value = Object.keys(streams)
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
  const masterPlaylistUrl = `${streamServer}/hls/${streamKey}.m3u8`

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

declare interface Recording {
  date: Date
  duration: string
  name: string
  path: string
}

function getRecordings() {
  window.axios
    .get(apiServer + '/recordings')
    .then((rec) => {
      console.log(rec)
      // recordings.value = rec
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
  const recordingUrl = `${streamServer}/vod/${path}`
  recordingVideo.value.src = recordingUrl
  recordingVideo.value.load()
  recordingVideo.value.play()
}

function saveChannel() {
  window.axios
    .post(streamServer + '/api/api/channel')
    .then((resp) => {
      console.log(resp)
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

onMounted(() => {
  getActiveStreams()
  getRecordings()
})
</script>

<template>
  <button @click="saveChannel()">Создать канал</button>
  <h1>Active Streams</h1>
  <div id="streams">
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
  </div>
</template>
