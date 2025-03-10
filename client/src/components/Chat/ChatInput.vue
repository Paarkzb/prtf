<script setup lang="ts">
import { FwbInput } from 'flowbite-vue'
import { ref } from 'vue'

const props = defineProps({
  send: {
    type: Function
  }
})

const text = ref('')
</script>

<template>
  <div>
    <fwb-input v-model="text" placeholder="Написать" label="Сообщение" @keydown="props.send!" />
  </div>
</template>
