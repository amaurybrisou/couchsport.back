import {
  GET_PAGES,
  GOT_PAGES,
  GET_PAGE,
  GOT_PAGE,
  SAVE_PAGE,
  PAGE_SAVED,
  PAGE_ADD_IMAGE,
  MODIFY_IMAGE_ALT,
  PAGE_DELETE_IMAGE,
  PAGE_IMAGE_DELETED,
  NEW_PAGE,
  NEW_PAGE_SAVED,
  PUBLISH_PAGE,
  PAGE_PUBLISHED,
  DELETE_PAGE,
  PAGE_DELETED,
  EDIT_PAGE,
  CANCEL_EDIT_PAGE,
  MODIFY_PAGE,
  REMOVE_ACTIVITY,
  PAGE_ERROR
} from '../actions/pages'

import { AUTH_LOGOUT } from '../actions/auth'

import pagesRepo from '../../repositories/pages'
import imagesRepo from '../../repositories/images'

import Vue from 'vue'

const state = {
  status: '',
  pages: {},
  edited_page: { Activities: [], Images: [] }
}

const getters = {}

function catchError ({ commit, dispatch }, resp) {
  if (resp.response.statusCode === 401) {
    commit(PAGE_ERROR)
    // if resp is unauthorized, logout, to
    dispatch(AUTH_LOGOUT)
  }
  throw resp
}
const actions = {
  [GET_PAGES]: ({ commit, dispatch }) => {
    commit(GET_PAGES)
    return pagesRepo
      .mines()
      .then(({ data }) => {
        commit(GOT_PAGES, data)
      })
      .catch(resp => {
        catchError({ commit, dispatch }, resp)
      })
  },
  [GET_PAGE]: ({ commit, dispatch }, params) => {
    commit(GET_PAGE)
    return pagesRepo
      .get(params)
      .then(({ data }) => {
        commit(GOT_PAGE)
        return data
      })
      .catch(resp => {
        catchError({ commit, dispatch }, resp)
      })
  },
  [SAVE_PAGE]: ({ commit, dispatch }, what) => {
    if (['edit', 'new'].indexOf(what) < 0) throw new Error('unknow method')

    what === 'edit' && commit(SAVE_PAGE)
    what === 'new' && commit(NEW_PAGE)

    return pagesRepo[what](state.edited_page)
      .then(({ data }) => {
        what === 'edit' && commit(PAGE_SAVED, data)
        what === 'new' && commit(NEW_PAGE_SAVED, data)
        commit(CANCEL_EDIT_PAGE)
      })
      .catch(resp => {
        catchError({ commit, dispatch }, resp)
      })
  },
  [DELETE_PAGE]: ({ commit, dispatch }, page) => {
    commit(DELETE_PAGE)
    return pagesRepo
      .delete(page)
      .then(() => {
        commit(PAGE_DELETED, page.ID)
      })
      .catch(resp => {
        catchError({ commit, dispatch }, resp)
      })
  },
  [PAGE_DELETE_IMAGE]: ({ commit, dispatch }, imageIDX) => {
    commit(PAGE_DELETE_IMAGE)
    if (!state.edited_page.Images[imageIDX].ID) { return commit(PAGE_IMAGE_DELETED, imageIDX) }
    return imagesRepo
      .delete(state.edited_page.Images[imageIDX])
      .then(() => {
        commit(PAGE_IMAGE_DELETED, imageIDX)
      })
      .catch(resp => {
        catchError({ commit, dispatch }, resp)
      })
  },
  [PUBLISH_PAGE]: ({ commit, dispatch }, page) => {
    commit(PUBLISH_PAGE, page.Public)
    return pagesRepo
      .publish(page)
      .then(() => {
        commit(PAGE_PUBLISHED, page)
      })
      .catch(resp => {
        catchError({ commit, dispatch }, resp)
      })
  }
}

const mutations = {
  [GET_PAGES]: state => {
    state.status = 'getting_pages'
  },
  [GOT_PAGES]: (state, pages) => {
    state.status = 'got_pages'
    state.pages = pages
  },

  [GET_PAGE]: state => {
    state.status = 'getting_one_page'
  },
  [GOT_PAGE]: state => {
    state.status = 'got_one_page'
  },

  [SAVE_PAGE]: state => {
    state.status = 'saving_page'
  },
  [PAGE_SAVED]: (state, { Images, ID }) => {
    for (var i = 0; i < state.pages.length; i++) {
      let p = state.pages[i]
      if (p.ID === ID) {
        Vue.set(state.pages[i], 'Images', Images)
        break
      }
    }
    state.status = 'page_saved'
  },

  [NEW_PAGE]: state => {
    state.status = 'saving_new_page'
  },
  [NEW_PAGE_SAVED]: (state, page) => {
    state.pages.push(page)
    state.status = 'new_page_saved'
  },

  [EDIT_PAGE]: (state, pageID) => {
    state.status = 'editing_page'
    for (var i = 0; i < state.pages.length; i++) {
      let p = state.pages[i]
      if (p.ID === pageID) {
        Vue.set(state, 'edited_page', state.pages[i])
        break
      }
    }
  },
  [CANCEL_EDIT_PAGE]: state => {
    Vue.set(state, 'edited_page', {
      Activities: [],
      Images: [],
      CouchNumber: 0
    })
    state.status = 'edit_page_canceled'
  },

  [MODIFY_PAGE]: (state, { key, value }) => {
    state.status = 'modifying_page'
    Vue.set(state.edited_page, key, value)
    state.status = 'page_modified'
  },
  [REMOVE_ACTIVITY]: (state, activity) => {
    state.status = 'removing_activity'
    state.edited_page.Activities = state.edited_page.Activities.filter(a => activity.ID !== a.ID)
    state.status = 'activity_removed'
  },
  [PAGE_ADD_IMAGE]: (state, image) => {
    state.edited_page.Images.push(image)
    // state.status = "page_adding_photo";
  },
  [PAGE_DELETE_IMAGE]: state => {
    state.status = 'page_deleting_image'
  },
  [MODIFY_IMAGE_ALT]: (state, { idx, value }) => {
    Vue.set(state.edited_page.Images[idx], 'Alt', value)
  },
  [PAGE_IMAGE_DELETED]: (state, imageIDX) => {
    state.edited_page.Images = state.edited_page.Images.filter((i, j) =>
      j !== imageIDX ? i : null
    )
    state.status = 'page_photo_deleteed'
  },

  [DELETE_PAGE]: state => {
    state.status = 'removing_page'
  },
  [PAGE_DELETED]: (state, pageID) => {
    state.pages = state.pages.filter(p => pageID !== p.ID)
    state.status = 'page_removed'
  },

  [PUBLISH_PAGE]: (state, Public) => {
    state.status = (Public ? '' : 'un') + 'publishing_page'
  },
  [PAGE_PUBLISHED]: (state, { ID, Public }) => {
    for (var i = 0; i < state.pages.length; i++) {
      let p = state.pages[i]
      if (p.ID === ID) {
        Vue.set(state.pages[i], 'Public', Public)
        break
      }
    }
    state.status = 'page_removed'
  }
}

export default {
  namespaced: true,
  state,
  getters,
  actions,
  mutations
}
