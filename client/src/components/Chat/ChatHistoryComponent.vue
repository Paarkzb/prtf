<script setup lang="ts">
import { ChatMessage } from './types'
import ChatMessageComponent from './ChatMessageComponent.vue'
import { type PropType } from 'vue'
import { FwbButton } from 'flowbite-vue'
import { useChatStore } from '@/stores/store'

const store = useChatStore()

const props = defineProps({
  chatHistory: Array as PropType<ChatMessage[]>
})
</script>

<template>
  <div class="bg-gray-200 m-0 p-[20px] rounded-md">
    <h2 class="text-center text-lg bg-gray-300 m-0 py-2 rounded-md">История чата</h2>
    <div>
      <p v-for="(msg, index) in props.chatHistory" :key="index">
        <ChatMessageComponent :message="msg" />
      </p>
    </div>
    <fwb-button class="mt-[10px]" @click="store.chatHistory = []">Clear chat history</fwb-button>
  </div>
</template>
