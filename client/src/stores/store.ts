import { defineStore } from 'pinia'
import { ChatMessage } from '@/components/Chat/types'
import { parse, stringify } from 'zipson'
import { Channel, ChannelData } from '@/components/Stream/types'
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

export const useChannelStore = defineStore('channel', {
  state: () => ({
    channel: {} as ChannelData
  }),
  actions: {
    login(channel: ChannelData) {
      this.channel.id = channel.id || ''
      this.channel.channel_name = channel.channel_name || ''
      this.channel.icon = channel.icon || ''
    },
    logout() {
      this.channel = {} as ChannelData
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
