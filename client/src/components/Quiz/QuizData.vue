<script setup lang="ts">
import { useRoute } from 'vue-router'
import { onMounted, ref, type Ref } from 'vue'
import Swal from 'sweetalert2'
import { v4 as uuidv4 } from 'uuid'
import QuizQuestion from './QuizQuestion.vue'
import router from '@/router'
import { FwbButton } from 'flowbite-vue'

const route = useRoute()

declare interface question {
  id: string
  index: number
  title: string
  answer: string
}

declare interface quiz {
  id: string
  user: {
    id: string
    name: string
    username: string
  }
  name: string
  description: string
  questions: question[]
}

const id = route.params.id
const quizData: Ref<quiz> = ref({
  id: '',
  user: {
    id: '',
    name: '',
    username: ''
  },
  name: '',
  description: '',
  questions: []
})

function getQuizData() {
  window.axios
    .get(window.quizApiURL + 'api/quiz/' + id)
    .then((response) => response.data)
    .then((data) => {
      console.log(data)
      quizData.value = data
    })
    .catch((error) => {
      console.log(error)
    })
}

function addQuestion() {
  quizData.value.questions.push({
    id: uuidv4(),
    index: quizData.value.questions.length + 1,
    title: '',
    answer: ''
  })
}

function deleteQuestion() {
  quizData.value.questions.pop()
}

function saveQuiz() {
  window.axios
    .put(window.quizApiURL + 'api/quiz/' + id, quizData.value)
    .then(function (response) {
      if (response.status === 200) {
        Swal.fire({
          title: 'Успех',
          text: 'Квиз сохранен',
          icon: 'success'
        })
      }
    })
    .catch(function (error) {
      console.log(error)
      Swal.fire({
        title: 'Ошибка',
        text: 'Квиз не сохранен',
        icon: 'error'
      })
    })
}

onMounted(() => {
  getQuizData()
})
</script>

<template>
  <!-- <div>{{ quizData }}</div> -->
  <div>
    <fwb-button color="yellow" @click="router.push('/admin/quiz')">Назад</fwb-button>
  </div>
  <div class="text-center">
    <div>
      <h1 class="text-xl font-bold">{{ quizData.name }}</h1>
    </div>
    <div>
      <h3>Создатель: {{ quizData?.user.username }}</h3>
    </div>
  </div>
  <div class="my-5">Описание: {{ quizData.description }}</div>
  <div>
    <form class="flex flex-col">
      <div class="flex flex-col gap-4">
        <QuizQuestion
          v-for="q in quizData.questions.sort((a, b) => a.index - b.index)"
          :key="q.id"
          :Index="q.index"
          :Title="q.title"
          :Answer="q.answer"
          @title="(msg) => (quizData.questions[q.index - 1].title = msg)"
          @answer="(msg) => (quizData.questions[q.index - 1].answer = msg)"
        />
      </div>
      <div class="my-4 flex justify-between">
        <fwb-button @click="saveQuiz">Сохранить</fwb-button>
        <fwb-button color="green" @click.prevent="addQuestion">Добавить вопрос</fwb-button>
        <fwb-button color="red" @click.prevent="deleteQuestion">Удалить вопрос</fwb-button>
      </div>
    </form>
  </div>
</template>
