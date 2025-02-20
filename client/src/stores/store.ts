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
      this.tokens.access_token = ''
      this.tokens.refresh_token = ''
    }
  },
  persist: true
})

interface ChannelData {
  id: string
  rf_user_id: string
  live: boolean
  rf_active_stream_id: string
  created_at: string
  updated_at: string
}

export const useChannelStore = defineStore('channel', {
  state: () => ({
    channel: {
      id: '',
      rf_user_id: '',
      live: false,
      rf_active_stream_id: '',
      created_at: '',
      updated_at: ''
    }
  }),
  actions: {
    login(channelData: ChannelData) {
      this.channel = channelData
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
