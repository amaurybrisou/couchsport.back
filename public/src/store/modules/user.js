import {
  MODIFY_PROFILE,
  SAVE_PROFILE,
  PROFILE_SAVED,
  PROFILE_ERROR,
  USER_REQUEST,
  USER_ERROR,
  USER_SUCCESS
} from "../actions/user";

import { SOCKET_CONNECT } from "../actions/ws";

import profileRepo from "../../repositories/profile";

import Vue from "vue";
import { AUTH_LOGOUT } from "../actions/auth";

const state = { status: "", profile: {} };

const getters = {
  getProfile: state => state.profile,
  isProfileLoaded: state => !!state.profile.ID
};

const actions = {
  [USER_REQUEST]: ({ commit, dispatch }) => {
    commit(USER_REQUEST);
    profileRepo
      .get()
      .then(({ data }) => {
        dispatch(SOCKET_CONNECT, data.ID);
        commit(USER_SUCCESS, data);
      })
      .catch(resp => {
        if (resp.response.status == 401) {
          commit(USER_ERROR);
          // if resp is unauthorized, logout, to
          dispatch(AUTH_LOGOUT);
        }
      });
  },
  [SAVE_PROFILE]: ({ commit, dispatch }) => {
    commit(SAVE_PROFILE);
    return profileRepo
      .update(state.profile)
      .then(({ data }) => {
        commit(PROFILE_SAVED, data);
      })
      .catch(resp => {
        if (resp.response.statusCode == 401) {
          commit(PROFILE_ERROR);
          // if resp is unauthorized, logout, to
          dispatch(AUTH_LOGOUT);
        }
        throw resp;
      });
  }
};

const mutations = {
  [USER_REQUEST]: state => {
    state.status = "loading";
  },
  [USER_SUCCESS]: (state, profile) => {
    state.status = "success";
    Vue.set(state, "profile", profile);
  },
  [USER_ERROR]: state => {
    state.status = "error";
  },

  [MODIFY_PROFILE]: (state, profile) => {
    state.profile = { ...state.profile, ...profile };
  },
  [SAVE_PROFILE]: state => {
    state.status = "loading";
  },
  [PROFILE_SAVED]: (state, profile) => {
    state.status = "success";
    Vue.set(state, "profile", profile);
  },
  [PROFILE_ERROR]: state => {
    state.status = "error";
  },
  [AUTH_LOGOUT]: state => {
    state.profile = {};
  }
};

export default {
  state,
  getters,
  actions,
  mutations
};
