<script setup lang="ts">
import { useRoute } from 'vue-router'
import { onMounted, ref } from 'vue'
import Swal from 'sweetalert2'
import { Channel, Recording } from '../types'
import VideoPlayer from '@/components/VideoPlayer.vue'
import router from '@/router'
import VideoPanel from '../VideoPanel.vue'
import { FwbButton, FwbTab, FwbTabs, FwbP } from 'flowbite-vue'
import ChannelAvatar from './ChannelAvatar.vue'
import moment from 'moment'
import ChatComponent from '../Chat/ChatComponent.vue'
import ChannelStream from './ChannelStream.vue'

const route = useRoute()
const id = route.params.id

const channelData = ref({} as Channel)

function getChannelData() {
  window.axios
    .get(window.gatewayURL + '/stream/api/channels/' + id)
    .then((rec) => {
      channelData.value = rec.data
    })
    .catch((error) => {
      console.log(error)
      Swal.fire({
        title: 'Ошибка',
        text: 'Неудалось получить записи',
        icon: 'error'
      })
    })
}

const recordings = ref([] as Recording[])

function getRecordings() {
  window.axios
    .get(window.gatewayURL + '/stream/api/channels/' + id + '/recordings')
    .then((rec) => {
      recordings.value = rec.data
    })
    .catch((error) => {
      console.log(error)
      Swal.fire({
        title: 'Ошибка',
        text: 'Неудалось получить записи',
        icon: 'error'
      })
    })
}

const activeTab = ref('first')

onMounted(() => {
  getChannelData()
  getRecordings()
})
</script>

<template>
  <div class="flex justify-between">
    <ChannelAvatar :channel="channelData" />
    <div class="text-white flex items-center">
      <fwb-button
        color="light"
        @click="router.push({ name: 'channelByIdSettings', params: { id: channelData.id } })"
      >
        Настройки
      </fwb-button>
    </div>
  </div>
  <div class="mt-10">
    <ChannelStream :channel="channelData" />
  </div>

  <div class="my-10" v-if="channelData">
    <ChatComponent :channel="channelData" />
  </div>

  <div class="mt-10">
    <fwb-tabs v-model="activeTab" class="p-5">
      <fwb-tab name="first" title="Описание">
        <fwb-p>
          Lorem ipsum dolor, sit amet consectetur adipisicing elit. Temporibus nemo quidem vero
          blanditiis dolor ut a dignissimos sint quis harum praesentium distinctio, aliquam dicta
          nulla ducimus alias non magnam! Eaque fugit doloribus nesciunt dolor accusantium incidunt,
          in odio ipsa libero nam qui rerum sint iure ratione repellendus voluptas animi dicta.
        </fwb-p>
      </fwb-tab>
      <fwb-tab name="second" title="Все видео">
        <div>
          <div class="grid grid-cols-3 gap-4">
            <div
              v-for="(rec, idx) in recordings?.sort((a, b) => moment(b.date).diff(a.date))"
              :key="idx"
            >
              <VideoPanel :recording="rec" />
            </div>
          </div>
        </div>
      </fwb-tab>
    </fwb-tabs>
  </div>
</template>
