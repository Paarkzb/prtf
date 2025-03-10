import axios from 'axios'
import router from '@/router'
import { useUserStore } from '@/stores/store'

const axiosInstance = axios.create({
  headers: {
    'Content-Type': 'application/json'
  }
})

axiosInstance.interceptors.request.use(
  function (config) {
    const store = useUserStore()
    if (store.tokens) {
      config.headers.Authorization = 'Bearer ' + store.tokens.access_token
    }
    return config
  },
  function (error) {
    return Promise.reject(error)
  }
)
axiosInstance.interceptors.response.use(
  function (response) {
    return response
  },
  function (error) {
    if (error.response?.status == 401) {
      const store = useUserStore()
      store.logout()
      router.push({ name: 'login' })
    }
    return Promise.reject(error)
  }
)

export default axiosInstance
