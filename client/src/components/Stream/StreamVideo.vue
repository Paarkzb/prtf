<script setup lang="ts">
import VideoPlayer from '@/components/VideoPlayer.vue'
import { useRoute } from 'vue-router'
import { onMounted, ref } from 'vue'
import { Recording } from './types'
import Swal from 'sweetalert2'

const route = useRoute()

const id = route.params.id as string

const videoOptions = ref()

const recording = ref<Recording>()

function getRecording(id: string) {
  window.axios
    .get(window.gatewayURL + '/stream/api/streams/' + id)
    .then((rec) => {
      recording.value = rec.data
      console.log('test', recording.value)
      setVideoOptions(recording.value!.path)
    })
    .catch((error) => {
      console.log(error)
      Swal.fire({
        title: 'Ошибка',
        text: 'Неудалось получить запись',
        icon: 'error'
      })
    })
}

const recordingVideo = ref({} as HTMLMediaElement)

function setVideoOptions(path: string) {
  videoOptions.value = {
    controls: true,
    preload: 'auto',
    width: 1000,
    aspectRatio: '16:9',
    playbackRates: [0.5, 0.75, 1, 1.25, 1.5, 1.75, 2],
    controlBar: {
      skipButtons: {
        forward: 5,
        backward: 5
      }
    },
    sources: [
      {
        src: `${window.gatewayURL}/stream/rec/${path}`,
        type: 'application/x-mpegURL'
      }
    ]
  }
}

onMounted(() => {
  getRecording(id)
})
</script>
<template>
  <div v-if="recording"><VideoPlayer ref="recordingVideo" :options="videoOptions" /></div>
</template>
