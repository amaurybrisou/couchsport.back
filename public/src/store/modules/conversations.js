import {
  GET_CONVERSATIONS,
  GOT_CONVERSATIONS,
  CONVERSATION_ADD_MESSAGE,
  CONVERSATION_MESSAGE_SENT,
  CONVERSATION_SEND_MESSAGE,
  NEW_CONVERSATION,
  REMOVE_CONVERSATION,
  CONVERSATION_REMOVED,
  MESSAGES_READ,
  CONVERSATION_ERROR
} from '../actions/conversations'

import conversationRepo from '../../repositories/conversation'
import Vue from 'vue'
import { AUTH_LOGOUT } from '../actions/auth'

const state = { status: '', conversations: [], unread: 0 }

const getters = {}

const actions = {
  [GET_CONVERSATIONS]: ({ commit, dispatch }) => {
    commit(GET_CONVERSATIONS)
    return conversationRepo
      .mines()
      .then(({ data }) => {
        commit(GOT_CONVERSATIONS, data)
      })
      .catch((resp) => {
        if (resp.response.statusCode === 401) {
          commit(CONVERSATION_ERROR)
          // if resp is unauthorized, logout, to
          dispatch(AUTH_LOGOUT)
        }
        throw resp
      })
  },
  [CONVERSATION_SEND_MESSAGE]: ({ commit, dispatch }, message) => {
    commit(CONVERSATION_SEND_MESSAGE)
    return conversationRepo
      .sendMessage(message)
      .then(({ data }) => {
        commit(CONVERSATION_MESSAGE_SENT, data)
        return data
      })
      .catch((resp) => {
        if (resp.response.statusCode === 401) {
          commit(CONVERSATION_ERROR)
          // if resp is unauthorized, logout, to
          dispatch(AUTH_LOGOUT)
        }
        throw resp
      })
  },
  [REMOVE_CONVERSATION]: ({ commit, dispatch }, id) => {
    commit(REMOVE_CONVERSATION)
    return conversationRepo
      .delete(id)
      .then(() => {
        commit(CONVERSATION_REMOVED, id)
      })
      .catch((resp) => {
        if (resp.response.statusCode === 401) {
          commit(CONVERSATION_ERROR)
          // if resp is unauthorized, logout, to
          dispatch(AUTH_LOGOUT)
        }
        throw resp
      })
  }
}

const mutations = {
  [MESSAGES_READ]: (state, conversationIDX) => {
    state.unread = 0
    if (conversationIDX > -1) {
      Vue.set(state.conversations[conversationIDX], 'unread', false)
    }
  },
  [GET_CONVERSATIONS]: (state) => {
    state.status = 'loading'
  },
  [GOT_CONVERSATIONS]: (state, conversations) => {
    state.status = 'get_success'
    state.conversations = conversations
  },
  [CONVERSATION_ADD_MESSAGE]: (state, message, rootState) => {
    state.unread++
    message.data = JSON.parse(message.data)
    for (var i = 0; i < state.conversations.length; i++) {
      var c = state.conversations[i]
      if (c.ID === message.data.OwnerID) {
        Vue.set(state.conversations[i], 'unread', true)
        state.conversations[i].Messages.push(message.data)
        break
      }
    }
  },
  [NEW_CONVERSATION]: (state, message) => {
    state.unread++
    let m = JSON.parse(message.data)
    m.unread = true
    state.conversations.push(m)
  },
  [CONVERSATION_SEND_MESSAGE]: (state) => {
    state.status = 'sending'
  },
  [CONVERSATION_MESSAGE_SENT]: (state, message) => {
    for (var i = 0; i < state.conversations.length; i++) {
      var c = state.conversations[i]
      if (c.ID === message.OwnerID) {
        state.conversations[i].Messages.push(message)
        break
      }
    }
    state.status = 'send_success'
  },
  [REMOVE_CONVERSATION]: (state) => {
    state.status = 'removing'
  },
  [CONVERSATION_REMOVED]: (state, id) => {
    state.conversations = state.conversations.filter((c) => id !== c.ID)
    state.status = 'remove_success'
  }
}

export default {
  namespaced: true,
  state,
  getters,
  actions,
  mutations
}
