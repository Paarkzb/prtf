<template>
  <div>
    <video ref="videoPlayer" class="video-js"></video>
  </div>
</template>

<script>
import videojs from 'video.js'
import 'video.js/dist/video-js.css'
import 'videojs-hls-quality-selector/src/plugin'

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
      this.player.hlsQualitySelector()
    }
  },
  beforeUnmount() {
    if (this.player) {
      this.player.dispose()
    }
  }
}
</script>
