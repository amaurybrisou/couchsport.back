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
  status: "",
  socket: {
    isConnected: false,
    reconnectError: false
  }
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
    state.status = "ws_error";
  },
  // default handler called for all methods
  // [SOCKET_ONMESSAGE]: (state, rootState, message) => {
  //   console.log(rootState)
  //   rootState.profile.conversations.conversations = MessageHub(state.messenger.conversations, message);
  //   state.messenger.unread_message = true;
  // },
  [SOCKET_RECONNECT]: (state, count) => {
    state.status = `ws_reconnecting ${count}`;
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
