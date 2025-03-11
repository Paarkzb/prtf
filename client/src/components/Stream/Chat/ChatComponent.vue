<script setup lang="ts">
import { onMounted, ref, defineProps, watch, type PropType } from 'vue'
import { ChatMessage, Channel } from '../types'
import { useChannelStore } from '@/stores/store'
import moment from 'moment'
import ChatHistory from './ChatHistory.vue'
import ChatInput from './ChatInput.vue'

const ws = ref({} as WebSocket)

const props = defineProps({
  channel: {
    type: Object as PropType<Channel>,
    required: true
  }
})

const channelStore = useChannelStore()

const chatMessages = ref([] as ChatMessage[])

const chatState = ref(false)

function connectWS() {
  console.log('Подключение к серверу чата...')

  ws.value = new WebSocket('ws://prtf.localhost:8090/stream/chat/ws')

  ws.value.onopen = (e) => {
    console.log('Успешное подключение')

    chatState.value = true

    sendMessage('старт')
  }

  ws.value.onmessage = (e) => {
    const msg: ChatMessage = JSON.parse(e.data)
    // console.log('Получено сообщение', msg)

    chatMessages.value.push(msg)
  }

  ws.value.onclose = function (e) {
    console.log('Соединение разорвано: ', e)

    chatState.value = false
  }

  ws.value.onerror = function (err) {
    console.log('Ошибка соединения: ', err)

    chatState.value = false
  }
}

function sendMessage(text: string) {
  if (text) {
    const msg: ChatMessage = {
      stream_channel_id: props.channel.id,
      stream_id: props.channel.rf_active_stream_id,
      text: text,
      time: moment().format(),
      channel: channelStore.channel
    }
    // console.log('Отправка сообщения', msg)
    ws.value.send(JSON.stringify(msg))
  }
}

watch(
  () => props.channel,
  (newChannel: Channel) => {
    if (newChannel && newChannel.id) {
      connectWS()
    }
  },
  { immediate: true }
)

// onMounted(() => {
//   connectWS()
//   console.log('channel', props.channel)
// })
</script>

<template>
  <div class="my-10 bg-zinc-800 rounded-md" v-if="chatState">
    <ChatHistory :chat-messages="chatMessages" />
    <ChatInput class="my-5" :send="sendMessage" />
  </div>
</template>
