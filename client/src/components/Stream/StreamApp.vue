<script setup lang="ts">
import { onMounted, ref } from 'vue'
import router from '@/router'
import Swal from 'sweetalert2'
import { useChannelStore } from '@/stores/store'
import { Channel } from './types'
import { FwbButton, FwbHeading } from 'flowbite-vue'
import ChannelAvatar from './ChannelAvatar.vue'
import { RouterLink } from 'vue-router'
import ChannelData from './ChannelData.vue'
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
const activeChannels = ref([] as Channel[])

function getActiveStreams() {
  window.axios
    .get(window.gatewayURL + '/stream/api/streams')
    .then((streams) => {
      activeChannels.value = streams.data
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
  <div>
    <fwb-button
      v-if="channelStore.channel.id"
      @click="router.push({ name: 'channelById', params: { id: channelStore.channel.id } })"
      color="green"
    >
      Мой канал
    </fwb-button>
    <fwb-button v-else color="yellow" @click="saveChannel()">Создать канал</fwb-button>
  </div>

  <fwb-heading tag="h1" class="text-center">Список каналов</fwb-heading>
  <div class="flex flex-wrap gap-x-4">
    <div v-for="(channel, idx) in channels" :key="idx">
      <router-link :to="{ name: 'channelById', params: { id: channel.id } }"
        ><ChannelAvatar :channel="channel" imgSize="md"
      /></router-link>
    </div>
  </div>

  <fwb-heading tag="h2" class="text-center mt-4">Активные каналы</fwb-heading>
  <div>
    <div v-for="channel in activeChannels" :key="channel.id">
      <router-link :to="{ name: 'channelById', params: { id: channel.id } }"
        ><ChannelAvatar :channel="channel" imgSize="md"
      /></router-link>
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
