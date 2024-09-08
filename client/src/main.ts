import './assets/main.css'
import './index.css'
import '../node_modules/flowbite-vue/dist/index.css'

import { createApp } from 'vue'
import { createPinia } from 'pinia'
import piniaPluginPersistedstate from 'pinia-plugin-persistedstate'

import App from './App.vue'
import router from './router'
import { type AxiosInstance } from 'axios'
import axiosInstanceConfig from '@/config/axiosInstanceConfig'

import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'
import { library } from '@fortawesome/fontawesome-svg-core'
import { fas } from '@fortawesome/free-solid-svg-icons'

library.add(fas)

const app = createApp(App)
const pinia = createPinia()
pinia.use(piniaPluginPersistedstate)

app.use(pinia)
app.use(router)
app.component('font-awesome-icon', FontAwesomeIcon)

declare global {
  interface Window {
    axios: AxiosInstance
    quizApiURL: string
    chatApiURL: string
  }
}

window.quizApiURL = 'http://localhost:8086/'
window.chatApiURL = 'http://localhost:8071/'

window.axios = axiosInstanceConfig

app.mount('#app')
