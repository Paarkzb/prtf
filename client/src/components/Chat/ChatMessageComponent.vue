<script setup lang="ts">
import type { PropType } from 'vue'
import { ChatMessage } from './types'

const props = defineProps({
  message: Object as PropType<ChatMessage>
})

const msg: ChatMessage = props.message!
</script>

<template>
  <div class="block bg-white mt-[10px] px-[20px] py-[10px] rounded-md clear-both shadow-inherit">
    {{ msg.body }}
  </div>
</template>
