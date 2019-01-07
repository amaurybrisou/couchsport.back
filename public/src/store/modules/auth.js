/* eslint-disable promise/param-names */
import { AUTH_REQUEST, AUTH_ERROR, AUTH_SUCCESS, AUTH_LOGOUT } from '../actions/auth'
import { USER_REQUEST } from '../actions/user'
import userRepository from '../../repositories/user'
import axios from '../../repositories/repository';

const state = { email: localStorage.getItem('user-email') || '', status: '', hasLoadedOnce: false }

const getters = {
  isAuthenticated: state => !!state.email,
  authStatus: state => state.status,
  email: state => state.email,
}

const actions = {
  [AUTH_REQUEST]: ({commit, dispatch}, user) => {
    return new Promise((resolve, reject) => {
      commit(AUTH_REQUEST)
      userRepository.login(user)
      .then(resp => {
        localStorage.setItem('user-email', resp.data.Email)
        // Here set the header of your ajax library to the token value.
        // example with axios
        // axios.defaults.headers.common = {'Authorization': `bearer ${resp.data.Token}`}
        axios.defaults.headers.common['Authorization'] = resp.data.Token
        commit(AUTH_SUCCESS, resp)
        dispatch(USER_REQUEST)
        resolve(resp)
      })
      .catch(err => {
        commit(AUTH_ERROR, err)
        localStorage.removeItem('user-email')
        reject(err)
      })
    })
  },
  [AUTH_LOGOUT]: ({commit, dispatch}) => {
    return new Promise((resolve, reject) => {
      commit(AUTH_LOGOUT)
      userRepository.logout()
      localStorage.removeItem('user-email')
      delete axios.defaults.headers.common['Authorization']
      resolve()
    })
  },
  [AUTH_ERROR]: ({commit, dispatch}) => {
    return new Promise((resolve, reject) => {
      commit(AUTH_LOGOUT)
      localStorage.removeItem('user-email')
      delete axios.defaults.headers.common['Authorization']
      resolve()
    })
  }
}

const mutations = {
  [AUTH_REQUEST]: (state) => {
    state.status = 'loading'
  },
  [AUTH_SUCCESS]: (state, resp) => {
    state.status = 'success'
    // state.token = resp.data.Token
    state.email = resp.data.Email,
    state.hasLoadedOnce = true
  },
  [AUTH_ERROR]: (state) => {
    state.status = 'error'
    state.hasLoadedOnce = true
  },
  [AUTH_LOGOUT]: (state) => {
    state.token = ''
  }
}

export default {
  state,
  getters,
  actions,
  mutations,
}
