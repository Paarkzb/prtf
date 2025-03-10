<script setup lang="ts">
import { defineProps } from 'vue'
import { ChatMessage } from '../types'
import moment from 'moment'
import { FwbP } from 'flowbite-vue'
import router from '@/router'

const props = defineProps({
  message: {
    type: ChatMessage,
    required: true
  }
})
</script>

<template>
  <fwb-p class="px-4 py-2 hover:bg-zinc-500 hover:rounded-md">
    <span>{{ moment(message.time).format('HH:mm') }} &nbsp;</span>
    <router-link :to="{ name: 'channelById', params: { id: message.channel.id } }"
      ><b class="hover:bg-zinc-600 hover:rounded-md"
        >{{ message.channel.channel_name }}
      </b></router-link
    >
    <span>:&nbsp;{{ message.text }}</span>
  </fwb-p>
</template>
