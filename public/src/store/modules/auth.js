/* eslint-disable promise/param-names */
import {
  AUTH_REQUEST,
  AUTH_ERROR,
  AUTH_SUCCESS,
  AUTH_LOGOUT
} from "../actions/auth";
import { PROFILE_REQUEST } from "../actions/profile";
import userRepository from "../../repositories/user";
import axios from "../../repositories/repository";

const state = {
  email: localStorage.getItem("user-email") || "",
  status: "",
  hasLoadedOnce: false
};

const getters = {
  isAuthenticated: state => !!state.email,
  authStatus: state => state.status,
  email: state => state.email
};

const actions = {
  [AUTH_REQUEST]: ({ commit, dispatch }, user) => {
    return new Promise((resolve, reject) => {
      commit(AUTH_REQUEST);
      userRepository
        .login(user)
        .then(resp => {
          localStorage.setItem("user-email", resp.data.Email);
          // Here set the header of your ajax library to the token value.
          // example with axios
          // axios.defaults.headers.common = {'Authorization': `bearer ${resp.data.Token}`}

          commit(AUTH_SUCCESS, resp);
          dispatch(PROFILE_REQUEST);
          resolve(resp);
        })
        .catch(({ response: { data } }) => {
          commit(AUTH_ERROR);
          reject(data);
        });
    });
  },
  [AUTH_LOGOUT]: ({ commit }) => {
    return new Promise((resolve, reject) => {
      commit(AUTH_LOGOUT);
      userRepository.logout();
      resolve();
    });
  },
  [AUTH_ERROR]: ({ commit }) => {
    return new Promise((resolve, reject) => {
      commit(AUTH_LOGOUT);
      resolve();
    });
  }
};

const mutations = {
  [AUTH_REQUEST]: state => {
    state.status = "loading";
  },
  [AUTH_SUCCESS]: (state, resp) => {
    state.status = "success";
    axios.defaults.headers.common.Authorization = resp.data.Token;
    state.email = resp.data.Email;
    state.hasLoadedOnce = true;
  },
  [AUTH_ERROR]: state => {
    state.status = "error";
    state.hasLoadedOnce = true;
  },
  [AUTH_LOGOUT]: state => {
    localStorage.removeItem("user-email");
    delete axios.defaults.headers.common.Authorization;
    state.email = null;
  }
};

export default {
  state,
  getters,
  actions,
  mutations
};
