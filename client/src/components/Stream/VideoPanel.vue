<script setup lang="ts">
import { defineProps, ref, type PropType } from 'vue'
import { Recording } from './types'
import router from '@/router'
import { FwbCard } from 'flowbite-vue'
import moment from 'moment'

const gateway = ref(window.gatewayURL)

const props = defineProps({
  recording: {
    type: Object as PropType<Recording>,
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
        {{
          hours > 0
            ? moment.utc(duration.asMilliseconds()).format('HH:mm:ss')
            : moment.utc(duration.asMilliseconds()).format('mm:ss')
        }}
      </div>
      <div>
        {{ moment(recording.date).format('DD.MM.YYYY HH:mm') }}
      </div>
    </div>
  </fwb-card>
</template>
