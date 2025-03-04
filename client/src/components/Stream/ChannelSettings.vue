<script setup lang="ts">
import { useRoute } from 'vue-router'
import { onMounted, ref, watchEffect } from 'vue'
import { Channel } from './types'
import Swal from 'sweetalert2'
const route = useRoute()

const id = route.params.id

const channelData = ref({} as Channel)

function getChannelData() {
  window.axios
    .get(window.gatewayURL + '/stream/api/channels/user')
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

onMounted(() => {
  getChannelData()
})
</script>

<template>
  <div>{{ channelData }}</div>
</template>
