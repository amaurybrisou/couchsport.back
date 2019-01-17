import Vue from "vue";
import Vuex from "vuex";
import profile from "./modules/profile";
import auth from "./modules/auth";
import ws from "./modules/ws";

Vue.use(Vuex);

const debug = process.env.NODE_ENV !== "production";

export default new Vuex.Store({
  modules: {
    profile,
    auth,
    ws
  },
  strict: debug,
  devTools: debug
});
