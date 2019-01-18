// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.
import Vue from "vue";
import Vuetify from "vuetify";
import "material-design-icons-iconfont/dist/material-design-icons.css";
import "vuetify/dist/vuetify.min.css";

Vue.use(Vuetify, {
  theme: {
    primary: "#06b9a8",
    secondary: "#009688",
    accent: "#00bcd4",
    error: "#ff5722",
    warning: "#ff9800",
    info: "#ffc107",
    success: "#607d8b"
  }
});

import AsyncComputed from "vue-async-computed";
Vue.use(AsyncComputed);

import App from "./App";
import router from "./router";
import store from "./store";

import {
  SOCKET_ONOPEN,
  SOCKET_ONCLOSE,
  SOCKET_ONERROR,
  SOCKET_ONMESSAGE,
  SOCKET_RECONNECT,
  SOCKET_RECONNECT_ERROR
} from "./store/actions/ws.js";

const mutations = {
  SOCKET_ONOPEN,
  SOCKET_ONCLOSE,
  SOCKET_ONERROR,
  SOCKET_ONMESSAGE,
  SOCKET_RECONNECT,
  SOCKET_RECONNECT_ERROR
};

import { i18n } from "./trans";
import FlagIcon from "vue-flag-icon";

Vue.use(FlagIcon);

import VueNativeSock from "vue-native-websocket";
Vue.use(
  VueNativeSock,
  `ws://${window.location.hostname}:${process.env.PORT || 8080}/api/ws`,
  {
    connectManually: true,
    store: store,
    mutations: mutations,
    format: "json",
    reconnection: true,
    reconnectionAttempts: 5
  }
);

import Filters from "./plugins/filter";
Vue.use(Filters);
import AppMessenger from "./plugins/messenger";

import {
  MESSAGES_READ,
  CONVERSATION_SEND_MESSAGE
} from "./store/actions/conversations";

Vue.use(AppMessenger, {
  namespace: "conversations",
  mutations: {
    MESSAGES_READ
  },
  actions: {
    CONVERSATION_SEND_MESSAGE
  },
  store
});

Vue.config.productionTip = false;
/* eslint-disable no-new */
new Vue({
  el: "#app",
  router,
  i18n,
  store,
  AppMessenger,
  components: { App },
  template: "<App/>"
});
