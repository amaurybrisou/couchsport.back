import Vue from "vue";
import Router from "vue-router";

import SignUp from "@/components/auth/SignUp";
import Login from "@/components/auth/Login";
import About from "@/components/About";
import Home from "@/components/Home";
import Explore from "@/components/explore/Explore";
import Profile from "@/components/profile/Profile";
import PageDetails from "@/components/page/PageDetails";

import { defaultLocale, i18n } from "../trans";
import store from "../store";

Vue.use(Router);

const ifNotAuthenticated = (to, from, next) => {
  if (!store.getters.isAuthenticated) {
    next();
    return;
  }
  next("/");
};

const ifAuthenticated = (to, from, next) => {
  if (store.getters.isAuthenticated) {
    next();
    return;
  }
  next("/login");
};

let router = new Router({
  mode: "history",
  routes: [
    {
      path: "/",
      redirect: `/${defaultLocale}`
    },
    {
      path: "/:locale",
      component: {
        template: "<router-view></router-view>"
      },
      children: [
        {
          path: "/",
          name: "home",
          component: Home
        },
        {
          path: "explore",
          name: "explore",
          component: Explore
        },
        {
          path: "signup",
          name: "signup",
          component: SignUp
        },
        {
          path: "pages/:page_id",
          name: "page-details",
          component: PageDetails
        },
        {
          path: "login",
          name: "login",
          component: Login,
          props: true,
          beforeEnter: ifNotAuthenticated
        },
        {
          path: "about",
          name: "about",
          component: About
        },
        {
          path: "profile",
          name: "profile",
          component: Profile,
          beforeEnter: ifAuthenticated
          // children: [
          //   {
          //     path: "informations",
          //     name: "informations",
          //     beforeEnter: ifAuthenticated
          //   },
          //   {
          //     path: "activities",
          //     name: "activities",
          //     beforeEnter: ifAuthenticated
          //   },
          //   {
          //     path: "conversations",
          //     name: "conversations",
          //     beforeEnter: ifAuthenticated
          //   },
          //   {
          //     path: "pages",
          //     name: "pages",
          //     beforeEnter: ifAuthenticated
          //   }
          // ]
        }
      ]
    }
  ]
});

router.beforeEach((to, from, next) => {
  let language = to.params.locale;
  if (!language) {
    language = defaultLocale;
  }

  i18n.locale = language;
  next();
});

export default router;
