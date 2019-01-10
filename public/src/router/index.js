import Vue from 'vue'
import Router from 'vue-router'

import SignUp from '@/components/auth/SignUp'
import Login from '@/components/auth/Login'
import About from '@/components/About'
import Home from '@/components/Home'
import Explore from '@/components/explore/Explore'
import Profile from '@/components/profile/Profile'
import PageDetails from '@/components/page/PageDetails'

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
      path: '/signup',
      name: 'signup',
      component: SignUp
    },
    {
      path: '/pages/:page_id',
      name: 'page-details',
      component: PageDetails
    },
    {
      path: '/login',
      name: 'login',
      component: Login,
      props: true,
      beforeEnter: ifNotAuthenticated,
    },
    {
      path: '/about',
      name: 'about',
      component: About
    }
  ]
})
