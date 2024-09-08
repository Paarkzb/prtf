<script setup lang="ts">
import { RouterView } from 'vue-router'
import { useUserStore } from '@/stores/store'
import {
  FwbNavbar,
  FwbNavbarCollapse,
  FwbNavbarLink,
  FwbDropdown,
  FwbListGroup,
  FwbListGroupItem
} from 'flowbite-vue'
const store = useUserStore()

function logout() {
  store.$reset()
}
</script>

<template>
  <fwb-navbar>
    <template #logo>
      <fwb-navbar-collapse>
        <fwb-navbar-link link="/"> <span class="text-2xl">Сервер</span> </fwb-navbar-link>
      </fwb-navbar-collapse>
    </template>
    <template #default="{ isShowMenu }">
      <fwb-navbar-collapse :is-show-menu="isShowMenu">
        <fwb-navbar-link is-active link="/"> Главная </fwb-navbar-link>
        <fwb-navbar-link link="/about"> Контакты </fwb-navbar-link>
        <fwb-dropdown v-if="store.isLogged">
          <template #trigger>
            <span
              class="cursor-pointer text-gray-700 hover:bg-gray-100 md:hover:bg-transparent md:border-0 md:hover:text-blue-700 dark:text-gray-400 md:dark:hover:text-white dark:hover:bg-gray-700 dark:hover:text-white md:dark:hover:bg-transparent"
              >Админ
            </span>
          </template>
          <fwb-list-group>
            <fwb-list-group-item>
              <fwb-navbar-link link="/admin/quiz"> Квизы </fwb-navbar-link>
            </fwb-list-group-item>
          </fwb-list-group>
        </fwb-dropdown>
        <fwb-navbar-link link="/chat"> Чат </fwb-navbar-link>
      </fwb-navbar-collapse>
    </template>
    <template #right-side>
      <fwb-navbar-collapse>
        <fwb-navbar-link link="/login" v-if="!store.isLogged"> Войти </fwb-navbar-link>
        <fwb-navbar-link @click="logout()" v-else>
          Выйти ({{ store.user?.username }})</fwb-navbar-link
        >
      </fwb-navbar-collapse></template
    >
  </fwb-navbar>

  <div class="bg-slate-800">
    <div class="max-w-screen-lg items-center mx-auto h-full min-h-[100vh] pt-5">
      <router-view></router-view>
    </div>
  </div>
</template>
