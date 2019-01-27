<template>
  <div id="app">
    <v-app>
      <v-carousel
        v-if="[`/${$i18n.locale}`, `/${$i18n.locale}/`, `/${$i18n.locale}/about`, `/${$i18n.locale}/signup`, `/${$i18n.locale}/login`].indexOf($route.path) >= 0"
        class="background-carousel"
        height="100vh"
        hide-delimiters
        hide-controls
      >
        <v-carousel-item
          v-for="(bg,i) in backgrounds"
          :key="i"
          class="carousel-item"
          lazy
          :src="bg.src"
        />
      </v-carousel>
      <app-nav />
      <v-content>
        <router-view />
      </v-content>
      <!-- <v-footer app></v-footer> -->
      <!-- <v-layout row justify-center>
        <v-dialog v-model="isNewOnSite" persistent max-width="290">
          <v-card>
            <v-card-title class="headline">{{ $t('p.app.new_on_site.modal_title') }}</v-card-title>
            <v-card-text>text</v-card-text>
            <v-card-actions>
              <v-spacer></v-spacer>
              <v-btn color="primary" flat @click.native="isNewOnSite = false">Disagree</v-btn>
              <v-btn color="primary" flat @click.native="isNewOnSite = false">Agree</v-btn>
            </v-card-actions>
          </v-card>
        </v-dialog>
      </v-layout>-->
    </v-app>
  </div>
</template>

<script>
import WebFontLoader from 'webfontloader'
import AppNav from 'components/nav/AppNav'

import { PROFILE_REQUEST } from 'store/actions/profile'
import { AUTH_ERROR } from 'store/actions/auth'
import axios from 'repositories/repository'

import L from '~/leaflet'
delete L.Icon.Default.prototype._getIconUrl

L.Icon.Default.mergeOptions({
  iconRetinaUrl: require('~/leaflet/dist/images/marker-icon-2x.png'),
  iconUrl: require('static/img/marker-icon.png'),
  shadowUrl: require('~/leaflet/dist/images/marker-shadow.png')
})

export default {
  name: 'App',
  components: {
    'app-nav': AppNav
  },
  data () {
    return {
      backgrounds: [
        { src: '/static/img/bg.jpg' },
        { src: '/static/img/bg1.jpg' },
        { src: '/static/img/bg2.jpg' },
        { src: '/static/img/bg3.jpg' }
      ]
    }
  },
  created: function () {
    if (this.$store.getters.isAuthenticated) {
      this.$store.dispatch(PROFILE_REQUEST)
    }

    var that = this
    axios.interceptors.response.use(
      function (response) {
        // Do something with response data
        return response
      },
      function (err) {
        if (err.response.status === 401) {
          // if you ever get an unauthorized, logout the user
          that.$store.dispatch(AUTH_ERROR)
          that.$router.push({ name: 'login' })

          // you can also redirect to /login if needed !
        }
        return Promise.reject(err)
      }
    )
  },
  mounted () {
    WebFontLoader.load({
      google: {
        families: ['Roboto:100,300,400,500,700,900']
      },
      active: this.setFontLoaded
    })
  },
  methods: {
    setFontLoaded () {
      this.$emit('font-loaded')
    }
  }
}
</script>

<style lang="stylus">

#app {
  background: none;

  .background-carousel {
    position: absolute;
    z-index: -1;
  }
}

.highlight-help {
  border-radius: 5px;
}

.help {

  tooltip-color = rgba(
    #fff,
    0.9
  );

  tooltip-text-color = rgba(
    #16191b,
    1
  );

  tooltip-border-color = rgba(
    #16191b,
    0.8
  );

  border-radius: 10px;
  background-color: tooltip-color;
  color: tooltip-text-color;

  .introjs-arrow.top {
    border-bottom-color: tooltip-color;
  }

  .introjs-arrow.bottom {
    border-top-color: tooltip-color;
  }

  .introjs-button {
    color: tooltip-text-color;
    border-color: tooltip-border-color;
  }

  .introjs-disabled {
    display: none;
  }

  .introjs-prevbutton {
    margin-right: 5px;
  }
}
</style>
