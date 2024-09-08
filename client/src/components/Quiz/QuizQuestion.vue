<script setup>
import { ref, watch } from 'vue'
import { FwbInput } from 'flowbite-vue'

const props = defineProps({
  Index: Number,
  Title: String,
  Answer: String
})

const emit = defineEmits(['title', 'answer'])

const titleValue = ref(props.Title)
const answerValue = ref(props.Answer)

function titleChange() {
  emit('title', titleValue)
}

function answerChange() {
  emit('answer', answerValue)
}

watch(titleValue, titleChange())
watch(answerValue, answerChange())
</script>

<template>
  <div class="flex">
    <div class="p-5 text-white self-center">{{ Index }}</div>
    <div class="flex-1">
      <div class="flex flex-col gap-2">
        <fwb-input :id="'question' + Index" v-model="titleValue" placeholder="Вопрос" />

        <fwb-input :id="'answer' + Index" v-model="answerValue" placeholder="Ответ" />
      </div>
    </div>
    <div class="p-5 text-white self-center">
      <div>+</div>
      <div>-</div>
    </div>
  </div>
</template>
