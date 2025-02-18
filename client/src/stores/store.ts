import { defineStore } from 'pinia'
import { ChatMessage } from '@/components/Chat/types'
import { parse, stringify } from 'zipson'

interface LoginResp {
  tokens: {
    access_token: string
    refresh_token: string
  }
  user: {
    id: string
    name: string
    username: string
  }
}

export const useUserStore = defineStore('user', {
  state: () => ({
    isLogged: false,
    tokens: {
      access_token: '',
      refresh_token: ''
    },
    user: {
      id: '',
      name: '',
      username: ''
    }
  }),
  actions: {
    login(loginResp: LoginResp) {
      this.isLogged = true
      this.tokens = loginResp.tokens
      this.user = loginResp.user
    },
    logout() {
      this.isLogged = false
    }
  },
  persist: true
})

export const useChatStore = defineStore('chat', {
  state: () => ({
    chatHistory: <ChatMessage[]>[]
  }),
  persist: {
    serializer: {
      deserialize: parse,
      serialize: stringify
    }
  }
})
