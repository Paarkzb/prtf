<script setup lang="ts">
import { defineProps, ref } from 'vue'
import { Recording } from './types'
import router from '@/router'
import { FwbCard } from 'flowbite-vue'
import moment from 'moment'

const gateway = ref(window.gatewayURL)

const props = defineProps({
  recording: {
    type: Recording,
    required: true
  }
})

const duration = moment.duration(props.recording.duration / 1000000, 'milliseconds')
const hours = duration.hours()
</script>

<template>
  <!-- :href="'/stream/video/' + recording.id" -->
  <fwb-card
    @click="router.push({ name: 'videoById', params: { id: recording.id } })"
    img-alt="no poster"
    :img-src="gateway + '/stream/rec/' + recording.poster"
    variant="image"
  >
    <div class="text-gray-100">
      <div>Название стрима</div>
      <div>{{ recording.channel_name }}</div>
      <div>
        {{ duration.humanize() + ' назад' }}
      </div>
    </div>
  </fwb-card>
</template>
