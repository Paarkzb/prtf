<script setup lang="ts">
import { onMounted, ref } from 'vue'
import router from '@/router'
import Swal from 'sweetalert2'
import { useChannelStore } from '@/stores/store'
import { Channel } from './types'

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

// Стрим
const activeStreams = ref({})

function getActiveStreams() {
  window.axios
    .get(window.gatewayURL + '/stream/api/streams')
    .then((streams) => {
      activeStreams.value = streams.data
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
  getChannels()
})
</script>

<template>
  <div class="flex gap-x-5">
    <button @click="saveChannel()" class="p-2 border bg-blue-300">Создать канал</button>
    <button
      @click="router.push({ name: 'channelById', params: { id: channelStore.channel.id } })"
      class="p-2 border bg-blue-300"
    >
      Мой канал
    </button>
  </div>

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

  <h1>Active Streams</h1>
  <div>
    <div v-for="stream in activeStreams" :key="stream">
      <button>{{ stream }}</button>
    </div>
  </div>

  <!-- <div class="chat-box">
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
