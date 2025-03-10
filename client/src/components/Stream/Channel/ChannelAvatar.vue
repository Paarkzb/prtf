<script setup lang="ts">
import { defineProps } from 'vue'
import { Channel } from '../types'
import { FwbAvatar } from 'flowbite-vue'
import type { AvatarSize } from 'node_modules/flowbite-vue/dist/components/FwbAvatar/types'
import type { PropType } from 'vue'

const props = defineProps({
  channel: {
    type: Object as PropType<Channel>,
    required: true
  },
  imgSize: {
    type: Object as PropType<AvatarSize>,
    default: () => {
      return 'lg' as AvatarSize
    }
  }
})
</script>

<template>
  <div>
    <router-link :to="{ name: 'channelById', params: { id: channel.id } }">
      <div class="flex gap-x-2 text-gray-100">
        <div>
          <fwb-avatar
            :img="channel.icon || ''"
            status-position="bottom-right"
            rounded
            :size="imgSize"
            :status="channel.live ? 'online' : 'offline'"
          />
        </div>
        <div class="flex items-center text-xl">
          {{ channel.channel_name }}
        </div>
      </div>
    </router-link>
  </div>
</template>
