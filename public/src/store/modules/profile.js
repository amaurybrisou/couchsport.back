import {
  MODIFY_PROFILE,
  SAVE_PROFILE,
  PROFILE_ERROR,
  PROFILE_REQUEST,
  PROFILE_SUCCESS,
  SET_ACTIVITIES,
  GET_ACTIVITIES,
  GET_LANGUAGES,
  SET_LANGUAGES
} from '../actions/profile'

import { SOCKET_CONNECT } from '../actions/ws'

import pages from './pages'
import conversations from './conversations'

import profileRepo from '../../repositories/profile'
import activityRepo from 'repositories/activity.js'
import languageRepo from '../../repositories/language.js'

import axios from 'repositories/repository.js'

import Vue from 'vue'
import { AUTH_LOGOUT } from '../actions/auth'

const state = {
  status: '',
  profile: {
    locale: 'fr'
  }
}

const getters = {
  getProfile: (state) => state.profile,
  isProfileLoaded: (state) => !!state.profile.ID,
  getLocale: (state) => state.profile.locale,
  activities: (state) => state.activities,
  languages: (state) => state.languages
}

const actions = {
  [PROFILE_REQUEST]: ({ commit, dispatch }) => {
    commit(PROFILE_REQUEST)
    profileRepo
      .get()
      .then(({ data }) => {
        dispatch(SOCKET_CONNECT, data.ID)
        commit(PROFILE_SUCCESS, data)
      })
      .catch((resp) => {
        if (resp.response && resp.response.status === 401) {
          commit(PROFILE_ERROR)
          // if resp is unauthorized, logout, to
          dispatch(AUTH_LOGOUT)
        }
      })
  },
  [SAVE_PROFILE]: ({ commit, dispatch }) => {
    commit(SAVE_PROFILE)
    return profileRepo
      .update(state.profile)
      .then(({ data }) => {
        commit(PROFILE_SUCCESS, data)
      })
      .catch((resp) => {
        if (resp.response.statusCode === 401) {
          commit(PROFILE_ERROR)
          // if resp is unauthorized, logout, to
          dispatch(AUTH_LOGOUT)
        }
        throw resp
      })
  },
  [GET_ACTIVITIES]: ({ dispatch }) => {
    if (sessionStorage.activities) {
      state.activities = JSON.parse(sessionStorage.activities)
      return state.activities
    }
    return dispatch(SET_ACTIVITIES)
  },
  [SET_ACTIVITIES]: ({ commit }) => {
    return activityRepo.all().then(({ data }) => {
      commit(SET_ACTIVITIES, data)
      return data
    })
  },
  [GET_LANGUAGES]: ({ dispatch }) => {
    if (sessionStorage.languages) {
      state.languages = JSON.parse(sessionStorage.languages)
      return state.languages
    }
    return dispatch(SET_LANGUAGES)
  },
  [SET_LANGUAGES]: ({ commit }) => {
    return languageRepo.all().then(({ data }) => {
      commit(SET_LANGUAGES, data)
      return data
    })
  }
}

const mutations = {
  [GET_LANGUAGES]: (state) => {
    state.status = 'loading_languages'
  },
  [SET_LANGUAGES]: (state, languages) => {
    Vue.set(state, 'languages', languages)
    sessionStorage.removeItem('languages')
    sessionStorage.setItem('languages', JSON.stringify(state.languages))
    state.status = 'languages_loaded'
  },
  [GET_ACTIVITIES]: (state) => {
    state.status = 'loading_activities'
  },
  [SET_ACTIVITIES]: (state, activities) => {
    Vue.set(state, 'activities', activities)
    sessionStorage.removeItem('activities')
    sessionStorage.setItem('activities', JSON.stringify(state.activities))
    state.status = 'activities_loaded'
  },
  [PROFILE_REQUEST]: (state) => {
    state.status = 'loading'
  },
  [PROFILE_SUCCESS]: (state, profile) => {
    state.status = 'success'
    Vue.set(state, 'profile', profile)
  },
  [PROFILE_ERROR]: (state) => {
    state.status = 'error'
  },
  [MODIFY_PROFILE]: (state, { key, value }) => {
    if (key === 'locale') {
      axios.defaults.headers.common['Accept-Language'] = value
    }
    Vue.set(state.profile, key, value)
  },
  [SAVE_PROFILE]: (state) => {
    state.status = 'loading'
  },
  [AUTH_LOGOUT]: (state) => {
    state.profile = {}
  }
}

export default {
  state,
  getters,
  actions,
  mutations,
  modules: {
    pages,
    conversations
  }
}
