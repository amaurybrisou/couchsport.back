import {
  MODIFY_PROFILE,
  SAVE_PROFILE,
  PROFILE_ERROR,
  PROFILE_REQUEST,
  PROFILE_SUCCESS
} from "../actions/profile";

import { SOCKET_CONNECT } from "../actions/ws";

import pages from "./pages";
import conversations from "./conversations";

import profileRepo from "../../repositories/profile";

import Vue from "vue";
import { AUTH_LOGOUT } from "../actions/auth";

const state = { status: "", profile: {} };

const getters = {
  getProfile: state => state.profile,
  isProfileLoaded: state => !!state.profile.ID
};

const actions = {
  [PROFILE_REQUEST]: ({ commit, dispatch }) => {
    commit(PROFILE_REQUEST);
    profileRepo
      .get()
      .then(({ data }) => {
        dispatch(SOCKET_CONNECT, data.ID);
        commit(PROFILE_SUCCESS, data);
      })
      .catch(resp => {
        if (resp.response.status == 401) {
          commit(PROFILE_ERROR);
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
        commit(PROFILE_SUCCESS, data);
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
  [PROFILE_REQUEST]: state => {
    state.status = "loading";
  },
  [PROFILE_SUCCESS]: (state, profile) => {
    state.status = "success";
    Vue.set(state, "profile", profile);
  },
  [PROFILE_ERROR]: state => {
    state.status = "error";
  },
  [MODIFY_PROFILE]: (state, { key, value }) => {
    state.profile[key] = value;
  },
  [SAVE_PROFILE]: state => {
    state.status = "loading";
  },
  [AUTH_LOGOUT]: state => {
    state.profile = {};
  }
};

export default {
  state,
  getters,
  actions,
  mutations,
  modules: {
    pages,
    conversations
  }
};
