<script setup lang="ts">
import { onMounted, ref } from 'vue'
import Swal from 'sweetalert2'
import CreateQuizModal from './QuizModal.vue'
import router from '@/router'
import {
  FwbButton,
  FwbTable,
  FwbTableBody,
  FwbTableCell,
  FwbTableHead,
  FwbTableHeadCell,
  FwbTableRow
} from 'flowbite-vue'

let quizes = ref([] as any)

function getAllQuiz() {
  window.axios
    .get(window.quizApiURL + 'api/quiz')
    .then((data) => data.data)
    .then(function (response) {
      console.log(response)
      quizes.value = response.data
    })
    .catch(function (error) {
      console.log(error)
      Swal.fire({
        title: 'Ошибка',
        text: 'Квизы не получены',
        icon: 'error'
      })
    })
}

function deleteQuiz(id: string) {
  window.axios
    .delete(window.quizApiURL + 'api/quiz/' + id)
    .then((response) => {
      if (response.status == 200) {
        Swal.fire({
          title: 'Удачно',
          text: 'Квиз удален',
          icon: 'success'
        })
        getAllQuiz()
      } else {
        Swal.fire({
          title: 'Ошибка',
          text: 'Квиз не удален',
          icon: 'error'
        })
      }
    })
    .catch((error) => {
      console.log(error)
      Swal.fire({
        title: 'Ошибка',
        text: 'Квиз не удален',
        icon: 'error'
      })
    })
}

onMounted(() => {
  getAllQuiz()
})

const modalTypes = ref(['create', 'edit'])
const showCreateQuizModal = ref(false)
const editId = ref('')
const editName = ref('')
const editDescription = ref('')
const type = ref('')

function showQuizModal(t: string, quiz = { id: '', name: '', description: '', questions: [] }) {
  showCreateQuizModal.value = true
  type.value = t
  editId.value = quiz?.id
  editName.value = quiz?.name
  editDescription.value = quiz?.description
}

function closeQuizModal() {
  showCreateQuizModal.value = false
  editName.value = ''
  editDescription.value = ''
}
</script>

<template>
  <div>
    <CreateQuizModal
      :show="showCreateQuizModal"
      :type="type"
      :editId="editId"
      :editName="editName"
      :editDescription="editDescription"
      @close="closeQuizModal()"
      @save="getAllQuiz()"
    >
    </CreateQuizModal>
  </div>

  <div class="flex justify-end"></div>

  <div>
    <fwb-table hoverable>
      <fwb-table-head>
        <fwb-table-head-cell>#</fwb-table-head-cell>
        <fwb-table-head-cell>Создатель</fwb-table-head-cell>
        <fwb-table-head-cell>Название</fwb-table-head-cell>
        <fwb-table-head-cell>Описание</fwb-table-head-cell>
        <fwb-table-head-cell>Количество вопросов</fwb-table-head-cell>
        <fwb-table-head-cell><span class="sr-only">Действие</span></fwb-table-head-cell>
      </fwb-table-head>
      <fwb-table-body>
        <fwb-table-row v-for="(q, index) in quizes" :key="q.id">
          <fwb-table-cell @click="router.push({ name: 'quizById', params: { id: q.id } })">{{
            index + 1
          }}</fwb-table-cell>
          <fwb-table-cell @click="router.push({ name: 'quizById', params: { id: q.id } })">{{
            q?.user.username
          }}</fwb-table-cell>
          <fwb-table-cell @click="router.push({ name: 'quizById', params: { id: q.id } })">{{
            q?.name
          }}</fwb-table-cell>
          <fwb-table-cell @click="router.push({ name: 'quizById', params: { id: q.id } })">{{
            q?.description
          }}</fwb-table-cell>
          <fwb-table-cell @click="router.push({ name: 'quizById', params: { id: q.id } })">{{
            q?.questions.length
          }}</fwb-table-cell>
          <fwb-table-cell>
            <fwb-button color="alternative" @click="showQuizModal(modalTypes[1], q)"
              ><font-awesome-icon icon="pencil"
            /></fwb-button>
            <fwb-button color="alternative" @click="deleteQuiz(q.id)"
              ><font-awesome-icon icon="trash"
            /></fwb-button>
          </fwb-table-cell>
        </fwb-table-row>
      </fwb-table-body>
    </fwb-table>
  </div>
</template>
