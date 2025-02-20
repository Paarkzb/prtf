<script setup lang="ts">
import { useRoute } from 'vue-router'
import { onMounted, ref } from 'vue'
import Swal from 'sweetalert2'

const route = useRoute()

const id = route.params.id

declare interface Channel {
  id: string
  rf_user_id: string
  live: boolean
  rf_active_stream_id: string
  created_at: Date
  updated_ad: Date
}

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

onMounted(() => {
  getChannelData()
})
</script>

<template>
  <div>
    <h1 class="text-6xl">{{ channelData.rf_user_id }}</h1>
    <h2 class="text-6xl">{{ channelData.live ? 'Онлайн' : 'Оффлайн' }}</h2>
  </div>
  <div>
    <button @click="startStream()">Запустить стрим</button>
    <div class="p-2 bg-white">{{ streamToken }}</div>
  </div>
</template>
