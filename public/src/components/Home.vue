<template>
  <v-container fluid>
    <v-layout>
      <v-flex xs12 sm4>
        <v-card>
          <v-toolbar color="secondary" dense>
            <v-toolbar-title>{{ $t('cosport') }}</v-toolbar-title>
          </v-toolbar>
          <v-card-text class="text-content">
            <h3
              class="subheading"
              :data-intro="$t('p.home.help.first_step')"
              data-step="1"
            >
              {{ $t('home.speech_title') }}
            </h3>
            <v-divider />
            <div
              class="body-1 mt-2"
              :data-intro="$t('p.home.help.second_step')"
              data-step="2"
            >
              {{ $t('home.speech_content') }}
            </div>
          </v-card-text>
          <!-- <v-flex class="text-xs-center"> -->
          <!-- <h1 class="display-3 font-weight-thin mb-3">CoSport</h1> -->
          <!-- </v-flex> -->
        </v-card>
      </v-flex>
    </v-layout>
  </v-container>
</template>

<script>
  import IntroJS from 'mixins/intro'

  export default {
    name: 'Home',
    mixins: [IntroJS],
    mounted() {
      var that = this
      this.help.setOption('doneLabel', this.$t('help.next_page'))
      this.help.oncomplete(function () {
        that.$router.push({ name: 'explore' })
      })
      this.help.start()
    }
  }
</script>

<style lang="stylus">
  .frontend-content {
    position: absolute;
    top: 10vh;
    background-color: rgba(255, 255, 255, 0.8) !important;
  }
</style>
