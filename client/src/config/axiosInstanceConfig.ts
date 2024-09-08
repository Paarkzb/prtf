import axios from 'axios'
import router from '@/router'

const axiosInstance = axios.create()
axiosInstance.defaults.timeout = 2500
axiosInstance.interceptors.request.use(
  function (config) {
    config.withCredentials = true
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
    if (error.response.status == 401) {
      router.push({ name: 'login' })
    }
    return Promise.reject(error)
  }
)

export default axiosInstance
