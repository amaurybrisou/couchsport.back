// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.
import Vue from 'vue'
import Vuetify from 'vuetify'

import 'assets/css/app.styl'

import '~/regenerator-runtime/runtime'
import { i18n } from './trans'
import VueNativeSock from 'vue-native-websocket'
import Filters from './plugins/filter'
import AppMessenger from './plugins/messenger'
import App from './App'
import router from './router'
import store from './store'

import {
  SOCKET_ONOPEN,
  SOCKET_ONCLOSE,
  SOCKET_ONERROR,
  SOCKET_ONMESSAGE,
  SOCKET_RECONNECT,
  SOCKET_RECONNECT_ERROR
} from './store/actions/ws.js'

import {
  MESSAGES_READ,
  CONVERSATION_SEND_MESSAGE
} from './store/actions/conversations'

Vue.use(Vuetify, {
  theme: {
    primary: '#06B998',
    secondary: '#06B998',
    accent: '#00bcd4',
    error: '#ff5722',
    warning: '#ff9800',
    info: '#ffc107',
    success: '#607d8b'
  }
})

const mutations = {
  SOCKET_ONOPEN,
  SOCKET_ONCLOSE,
  SOCKET_ONERROR,
  SOCKET_ONMESSAGE,
  SOCKET_RECONNECT,
  SOCKET_RECONNECT_ERROR
}
Vue.use(
  VueNativeSock,
  `ws://${window.location.hostname}:${process.env.PORT || 8080}/api/ws`,
  {
    connectManually: true,
    store: store,
    mutations: mutations,
    format: 'json',
    reconnection: true,
    reconnectionAttempts: 5
  }
)
Vue.use(Filters)

Vue.use(AppMessenger, {
  namespace: 'conversations',
  mutations: {
    MESSAGES_READ
  },
  actions: {
    CONVERSATION_SEND_MESSAGE
  },
  store
})

/* eslint-disable-next-line no-new */
new Vue({
  el: '#app',
  router,
  i18n,
  store,
  AppMessenger,
  render: h => h(App)
})
