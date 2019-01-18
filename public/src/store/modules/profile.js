import {
  MODIFY_PROFILE,
  SAVE_PROFILE,
  PROFILE_ERROR,
  PROFILE_REQUEST,
  PROFILE_SUCCESS,
  SET_ACTIVITIES,
  GET_ACTIVITIES,
  GET_LANGUAGES,
  SET_LANGUAGES,
  SET_LOCALE
} from "../actions/profile";

import { SOCKET_CONNECT } from "../actions/ws";

import pages from "./pages";
import conversations from "./conversations";

import profileRepo from "../../repositories/profile";
import activityRepo from "@/repositories/activity.js";
import languageRepo from "../../repositories/language.js";

import { i18n } from "@/trans";

import Vue from "vue";
import { AUTH_LOGOUT } from "../actions/auth";
import { lang } from "moment";

const state = { status: "", profile: { locale: "fr" }, activities: [] };

const getters = {
  getProfile: state => state.profile,
  isProfileLoaded: state => !!state.profile.ID,
  getLocale: state => state.profile.locale,
  activities: state => state.activities,
  languages: state => state.languages
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
  },
  [SET_LOCALE]: ({ dispatch }, locale) => {
    dispatch(SET_ACTIVITIES, locale);
  },
  [GET_ACTIVITIES]: ({ dispatch }, locale) => {
    return dispatch(SET_ACTIVITIES, locale);
  },
  [SET_ACTIVITIES]: ({ commit }, locale) => {
    return new Promise((resolve, reject) => {
      if (!localStorage.saved_activities) {
        if (locale !== state.profile.locale) {
          commit(MODIFY_PROFILE, { key: "locale", value: locale });
        }
        return activityRepo
          .all()
          .then(({ data }) => {
            commit(SET_ACTIVITIES, data);
            return resolve(state.activities);
          })
          .catch(res => {
            return reject(res);
          });
      } else {
        if (
          locale === state.profile.locale &&
          localStorage.activities &&
          state.activities.length > 0
        ) {
          return resolve(state.activities);
        }

        commit(MODIFY_PROFILE, { key: "locale", value: locale });
        var data = JSON.parse(localStorage.getItem("saved_activities"));
        commit(SET_ACTIVITIES, data);
        return resolve(state.activities);
      }
      reject();
    });
  },
  [GET_LANGUAGES]: ({ commit }, locale) => {
    return new Promise((resolve, reject) => {
      if (!localStorage.saved_languages) {
        if (locale !== state.profile.locale) {
          commit(MODIFY_PROFILE, { key: "locale", value: locale });
        }
        return languageRepo
          .all()
          .then(({ data }) => {
            commit(SET_LANGUAGES, data);
            return resolve(state.languages);
          })
          .catch(res => {
            return reject(res);
          });
      } else {
        if (
          locale === state.profile.locale &&
          localStorage.languages &&
          state.languages.length > 0
        ) {
          return resolve(state.languages);
        }

        commit(MODIFY_PROFILE, { key: "locale", value: locale });
        let data = JSON.parse(localStorage.getItem("saved_languages"));
        commit(SET_LANGUAGES, data);
        return resolve(state.languages);
      }
      reject();
    });
  }
};

const mutations = {
  [GET_LANGUAGES]: state => {
    state.status = "loading_languages";
  },
  [SET_LANGUAGES]: (state, languages) => {
    //KEEP THIS SNIPPET IF FUTURE TRANSLATIONS ARE NEEDED
    localStorage.setItem("saved_languages", JSON.stringify(languages));
    // state.languages = languages.map(function(e){
    //   e.Name = i18n.t(`allActivities.${e["Name"]}`)
    //   return e;
    // })
    state.languages = languages;
    localStorage.removeItem("languages");
    localStorage.setItem("languages", JSON.stringify(state.languages));
    state.status = "languages_loaded";
  },
  [GET_ACTIVITIES]: state => {
    state.status = "loading_activities";
  },
  [SET_ACTIVITIES]: (state, activities) => {
    localStorage.setItem("saved_activities", JSON.stringify(activities));
    state.activities = activities.map(function(e) {
      e.Name = i18n.t(`allActivities.${e["Name"]}`);
      return e;
    });
    localStorage.removeItem("activities");
    localStorage.setItem("activities", JSON.stringify(state.activities));
    state.status = "activities_loaded";
  },
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
    Vue.set(state.profile, key, value);
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
