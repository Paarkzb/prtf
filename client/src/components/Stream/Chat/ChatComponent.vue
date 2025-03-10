<script setup lang="ts">
import { onMounted, ref, defineProps } from 'vue'
import { ChatMessage, Channel } from '../types'
import { useChannelStore } from '@/stores/store'
import moment from 'moment'
import ChatHistory from './ChatHistory.vue'
import ChatInput from './ChatInput.vue'

const ws = new WebSocket('ws://prtf.localhost:8090/stream/chat/ws')

const props = defineProps({
  channel: {
    type: Channel,
    required: true
  }
})

const channelStore = useChannelStore()

const chatMessages = ref([] as ChatMessage[])

function connectWS() {
  console.log('Подключение к серверу чата...')

  ws.onopen = (e) => {
    console.log('Успешное подключение: ', e)

    sendMessage('старт')
  }

  ws.onmessage = (e) => {
    const msg: ChatMessage = JSON.parse(e.data)
    console.log('сообщение', msg)

    chatMessages.value.push(msg)
  }

  ws.onclose = function (e) {
    console.log('Соединение разорвано: ', e)
  }

  ws.onerror = function (err) {
    console.log('Ошибка соединения: ', err)
  }
}

function sendMessage(text: string) {
  if (text) {
    const msg: ChatMessage = {
      stream_id: props.channel.rf_active_stream_id,
      text: text,
      time: moment().format(),
      channel: channelStore.channel
    }
    console.log(msg)
    ws.send(JSON.stringify(msg))
  }
}

onMounted(() => {
  connectWS()
})
</script>

<template>
  <div class="my-10 bg-zinc-800 rounded-md">
    <ChatHistory :chat-messages="chatMessages" />
    <ChatInput class="my-5" :send="sendMessage" />
  </div>
</template>
