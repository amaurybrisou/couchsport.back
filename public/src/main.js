// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.
import Vue from "vue";
import Vuetify from "vuetify";

Vue.use(Vuetify, {
  theme: {
    primary: "#06b9a8",
    secondary: "#009688",
    accent: "#00bcd4",
    error: "#ff5722",
    warning: "#ff9800",
    info: "#ffc107",
    success: "#607d8b"
  }
});

import "material-design-icons-iconfont/dist/material-design-icons.css";
import "vuetify/dist/vuetify.min.css";
import AsyncComputed from "vue-async-computed";
import moment from "moment";

Vue.filter("formatDate", function(value, format) {
  if (value) {
    return moment(String(value)).format(format || "MM/DD/YYYY hh:mm");
  }
});

Vue.use(AsyncComputed);

import App from "./App";
import router from "./router";

import store from "./store";

Vue.config.productionTip = false;

/* eslint-disable no-new */
new Vue({
  el: "#app",
  router,
  store,
  components: { App },
  template: "<App/>"
});
