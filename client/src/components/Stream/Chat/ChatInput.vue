<script setup lang="ts">
import { ref, defineProps } from 'vue'
import { FwbInput, FwbButton } from 'flowbite-vue'

const text = ref('')

const props = defineProps({
  send: {
    type: Function,
    required: true
  }
})

function sendMessage() {
  props.send(text.value)
  text.value = ''
}
</script>

<template>
  <div class="flex gap-x-2">
    <div class="flex-1">
      <fwb-input
        size="sm"
        v-model="text"
        v-on:keypress="
          (e) => {
            if (e.key === 'Enter') sendMessage()
          }
        "
        placeholder="Отправить сообщение"
      />
    </div>
    <fwb-button color="default" pill square @click="sendMessage()"
      ><svg
        class="w-5 h-5"
        fill="currentColor"
        viewBox="0 0 20 20"
        xmlns="http://www.w3.org/2000/svg"
      >
        <path
          clip-rule="evenodd"
          d="M10.293 3.293a1 1 0 011.414 0l6 6a1 1 0 010 1.414l-6 6a1 1 0 01-1.414-1.414L14.586 11H3a1 1 0 110-2h11.586l-4.293-4.293a1 1 0 010-1.414z"
          fill-rule="evenodd"
        />
      </svg>
    </fwb-button>
  </div>
</template>
