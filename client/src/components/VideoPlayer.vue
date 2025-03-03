<template>
  <div>
    <video ref="videoPlayer" class="video-js"></video>
  </div>
</template>

<script>
import videojs from 'video.js'
import 'video.js/dist/video-js.css'

export default {
  name: 'VideoPlayer',
  props: {
    options: {
      type: Object,
      default() {
        return {}
      }
    }
  },
  data() {
    return {
      player: null
    }
  },
  mounted() {
    this.initializePlayer()
  },
  methods: {
    initializePlayer() {
      this.player = videojs(this.$refs.videoPlayer, this.options, () => {
        this.player.log('onPlayerReady', this)
      })

      // this.player.ready(() => {
      //   this.player.hlsQualitySelector({
      //     displayCurrentQuality: true
      //   })
      // })
    }
  },
  beforeUnmount() {
    if (this.player) {
      this.player.dispose()
    }
  }
}
</script>
