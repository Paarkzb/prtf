<script setup lang="ts">
import router from '@/router'
import Swal from 'sweetalert2'
import { ref } from 'vue'
import { useUserStore } from '@/stores/store'
import { FwbInput, FwbButton } from 'flowbite-vue'

const email = ref('')
const username = ref('')
const password = ref('')
const passwordConfirm = ref('')

const store = useUserStore()

function signUp() {
  let signUpData = {
    email: email.value,
    username: username.value,
    password: password.value
  }
  let signInData = {
    username: username.value,
    password: password.value
  }
  window.axios
    .post(window.gatewayURL + '/auth/v1/sign-up', signUpData)
    .then((response) => response.data)
    .then(function (response) {
      window.axios
        .post(window.gatewayURL + '/auth/v1/sign-in', signInData)
        .then((resp) => resp.data)
        .then((resp) => {
          store.login(resp)
        })
        .catch(function (error) {
          console.log(error)
          Swal.fire({
            title: 'Ошибка',
            text: 'Неправильный логин или пароль',
            icon: 'error'
          })
        })
    })
    .then(() => {
      router.push({ name: 'home' })
    })
    .catch(function (error) {
      console.log(error)
      Swal.fire({
        title: 'Ошибка',
        text: 'Логин или почта занята',
        icon: 'error'
      })
    })
}
</script>

<template>
  <form class="flex flex-col mx-auto pt-20 max-w-sm content-center">
    <fwb-input v-model="email" placeholder="Почта" label="Почта" />
    <fwb-input v-model="username" placeholder="Логин" label="Логин" />
    <fwb-input v-model="password" placeholder="Пароль" label="Пароль" type="password" />
    <fwb-input
      v-model="passwordConfirm"
      placeholder="Подтверждение пароля"
      label="Подтверждение пароля"
      type="password"
    />

    <fwb-button class="mt-5" color="default" @click.prevent="signUp()"
      >Зарегистрироваться</fwb-button
    >
  </form>
</template>
