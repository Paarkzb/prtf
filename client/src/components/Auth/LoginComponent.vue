<script setup lang="ts">
import router from '@/router'
import Swal from 'sweetalert2'
import { ref } from 'vue'
import { useUserStore } from '@/stores/store'
import { FwbInput, FwbButton } from 'flowbite-vue'

const username = ref('')
const password = ref('')

const store = useUserStore()

function login() {
  let loginData = {
    username: username.value,
    password: password.value
  }
  window.axios
    .post(window.quizApiURL + 'auth/sign-in', loginData)
    .then((response) => response.data)
    .then(function (response) {
      store.login(response)
    })
    .then(() => {
      router.push({ name: 'home' })
    })
    .catch(function (error) {
      console.log(error)
      Swal.fire({
        title: 'Ошибка',
        text: 'Неправильный логин или пароль',
        icon: 'error'
      })
    })
}
</script>

<template>
  <form class="flex flex-col mx-auto pt-20 max-w-sm content-center">
    <fwb-input v-model="username" placeholder="Логин" label="Логин" />
    <fwb-input v-model="password" placeholder="Пароль" label="Пароль" type="password" />

    <fwb-button class="mt-5" color="default" @click.prevent="login">Войти</fwb-button>
  </form>
</template>
