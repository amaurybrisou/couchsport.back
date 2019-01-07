import {
  MODIFY_PROFILE,
  PROFILE_MODIFIED,
  SAVE_PROFILE,
  PROFILE_SAVED,
  PROFILE_ERROR,
  USER_REQUEST,
  USER_ERROR,
  USER_SUCCESS,
  MODIFY_PROFILE_ACTIVITY
} from "../actions/user";
import profileRepo from "../../repositories/profile";
import Vue from "vue";
import { AUTH_LOGOUT } from "../actions/auth";

const state = { status: "", profile: {} };

const getters = {
  getProfile: state => state.profile,
  isProfileLoaded: state => !!state.profile.UserID
};

const actions = {
  [USER_REQUEST]: ({ commit, dispatch }) => {
    commit(USER_REQUEST);
    profileRepo
      .get()
      .then(({ data }) => {
        commit(USER_SUCCESS, data);
      })
      .catch(resp => {
        commit(USER_ERROR);
        // if resp is unauthorized, logout, to
        dispatch(AUTH_LOGOUT);
      });
  },
  [MODIFY_PROFILE]: ({ commit, dispatch }, updates) => {
    commit(MODIFY_PROFILE);
    var profile = { ...state.profile, ...updates };
    commit(PROFILE_MODIFIED, profile);
  },
  [MODIFY_PROFILE_ACTIVITY]: ({ commit, dispatch }, activities) => {
    commit(MODIFY_PROFILE_ACTIVITY, activities)
  },
  [SAVE_PROFILE]: ({ commit, dispatch }) => {
    commit(SAVE_PROFILE);
    return profileRepo
      .update(state.profile)
      .then(({ data }) => {
        commit(PROFILE_SAVED, data);
      })
      .catch(resp => {
        commit(PROFILE_ERROR);
        // if resp is unauthorized, logout, to
        dispatch(AUTH_LOGOUT);
        throw resp
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
  [MODIFY_PROFILE]: (state, payload) => {
    state.status = "modifying";
  },
  [PROFILE_MODIFIED]: (state, profile) => {
    state.status = "success";
    Vue.set(state, "profile", profile);
    console.log(state.profile);
  },
  [SAVE_PROFILE]: state => {
    state.status = "loading";
  },
  [PROFILE_SAVED]: (state, profile) => {
    state.status = "success";
    Vue.set(state, "profile", profile);
  },
  [MODIFY_PROFILE_ACTIVITY]: (state, activities) => {
    state.profile.Activities = activities
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
