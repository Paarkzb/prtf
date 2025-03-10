<script setup lang="ts">
import { onMounted, ref } from 'vue'
import router from '@/router'
import Swal from 'sweetalert2'
import { useChannelStore } from '@/stores/store'
import { Channel } from './types'
import { FwbButton, FwbHeading } from 'flowbite-vue'
import ChannelAvatar from './Channel/ChannelAvatar.vue'
import { RouterLink } from 'vue-router'

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
      let channel: Channel = resp.data
      channelStore.login(channel)
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
      <ChannelAvatar :channel="channel" imgSize="md" />
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
</template>
