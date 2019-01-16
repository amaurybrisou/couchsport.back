import {
  SOCKET_ONOPEN,
  SOCKET_ONCLOSE,
  SOCKET_ONERROR,
  SOCKET_ONMESSAGE,
  SOCKET_RECONNECT,
  SOCKET_RECONNECT_ERROR,
  EMIT,
  SOCKET_CONNECT,
  MESSAGES_READ
} from "../actions/ws.js";

import Vue from "vue";

const state = {
  socket: {
    isConnected: false,
    reconnectError: false
  },
  messages: []
};

const getters = {};

const actions = {
  [EMIT]: ({ commit, dispatch }, message) => {
    if (state.socket.isConnected) {
      Vue.prototype.$socket.sendObj(message);
    }
  },
  [SOCKET_CONNECT]: ({ commit, dispatch }, profileID) => {
    Vue.prototype.$connect(
      `ws://${window.location.hostname}:${process.env.PORT ||
        8080}/api/ws?id=${profileID}`
    );
  }
};

const mutations = {
  [SOCKET_ONOPEN]: state => {
    state.socket.isConnected = true;
  },
  [SOCKET_ONCLOSE]: state => {
    state.socket.isConnected = false;
  },
  [SOCKET_ONERROR]: state => {
    console.error(state, event);
  },
  // default handler called for all methods
  [SOCKET_ONMESSAGE]: (state, message) => {
    state.messages.push(message);
  },
  [MESSAGES_READ]: (state, message) => {
    state.messages = [];
  },
  // mutations for reconnect methods
  [SOCKET_RECONNECT]: (state, count) => {
    console.info(state, count);
  },
  [SOCKET_RECONNECT_ERROR]: state => {
    state.socket.reconnectError = true;
  }
};

export default {
  state,
  getters,
  actions,
  mutations
};
