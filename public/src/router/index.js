import Vue from 'vue'
import Router from 'vue-router'

import SignIn from '@/components/auth/SignIn'
import Login from '@/components/auth/Login'
import About from '@/components/About'
import Home from '@/components/Home'
import Explore from '@/components/explore/Explore'
import Profile from '@/components/profile/Profile'

import store from '../store'


Vue.use(Router)


const ifNotAuthenticated = (to, from, next) => {
  if (!store.getters.isAuthenticated) {
    next()
    return
  }
  next('/')
}

const ifAuthenticated = (to, from, next) => {
  if (store.getters.isAuthenticated) {
    next()
    return
  }
  next('/login')
}

export default new Router({
  mode: 'history',
  routes: [
    {
      path: '/',
      name: 'home',
      component: Home
    },
    {
      path: '/explore',
      name: 'explore',
      component: Explore
    },
    {
      path: '/profile',
      name: 'profile',
      component: Profile,
      beforeEnter: ifAuthenticated,
    },
    {
      path: '/signin',
      name: 'signin',
      component: SignIn
    },
    {
      path: '/login',
      name: 'login',
      component: Login,
      beforeEnter: ifNotAuthenticated,
    },
    {
      path: '/about',
      name: 'about',
      component: About
    }
  ]
})
