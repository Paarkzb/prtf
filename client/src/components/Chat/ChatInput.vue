<script setup lang="ts">
import { FwbInput } from 'flowbite-vue'

const props = defineProps({
  send: {
    type: Function
  }
})
</script>

<template>
  <div>
    <fwb-input placeholder="Написать" label="Сообщение" @keydown="props.send!" />
    <!-- <input placeholder="Написать" label="Сообщение" @keydown="props.send" /> -->
  </div>
</template>
