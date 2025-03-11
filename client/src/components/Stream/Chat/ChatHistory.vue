<script setup lang="ts">
import { defineProps } from 'vue'
import ChatMessage from './ChatMessage.vue'
import { FwbHeading } from 'flowbite-vue'
import { ChatMessage as chatMsg } from '../types'

const props = defineProps({
  chatMessages: {
    type: Array<chatMsg>,
    required: true
  }
})
</script>

<template>
  <div>
    <fwb-heading tag="h4" class="text-center py-2">Чат трансляции</fwb-heading>
    <hr class="border-b border-white" />
    <div v-for="msg in chatMessages" :key="msg.text">
      <ChatMessage :message="msg" />
    </div>
  </div>
</template>
