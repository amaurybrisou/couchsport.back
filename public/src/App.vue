<template>
  <div id="app">
    <v-app>
      <v-carousel
        v-if="[`/${$i18n.locale}`, `/${$i18n.locale}/`, `/${$i18n.locale}/about`, `/${$i18n.locale}/signup`, `/${$i18n.locale}/login`].indexOf($route.path) >= 0"
        class="background-carousel"
        height="100vh"
        hide-delimiters
        hide-controls
      >
        <v-carousel-item
          class="carousel-item"
          lazy
          v-for="(bg,i) in backgrounds"
          :key="i"
          :src="bg.src"
        ></v-carousel-item>
      </v-carousel>
      <app-nav></app-nav>
      <v-content>
        <router-view></router-view>
      </v-content>
      <!-- <v-footer app></v-footer> -->
    </v-app>
  </div>
</template>

<script>
import WebFontLoader from "webfontloader";
import AppNav from "@/components/nav/AppNav";

import { PROFILE_REQUEST } from "@/store/actions/profile";
import { AUTH_ERROR } from "@/store/actions/auth";
import axios from "@/repositories/repository";

import L from "leaflet";
delete L.Icon.Default.prototype._getIconUrl;

L.Icon.Default.mergeOptions({
  iconRetinaUrl: require("leaflet/dist/images/marker-icon-2x.png"),
  iconUrl: require("./../static/img/marker-icon.png"),
  shadowUrl: require("leaflet/dist/images/marker-shadow.png")
});

export default {
  name: "App",
  components: {
    "app-nav": AppNav
  },
  data() {
    return {
      backgrounds: [
        { src: "/static/img/bg.jpg" },
        { src: "/static/img/bg1.jpg" },
        { src: "/static/img/bg2.jpg" },
        { src: "/static/img/bg3.jpg" }
      ],
      navBarLoaded: false
    };
  },
  created: function() {
    if (this.$store.getters.isAuthenticated) {
      this.$store.dispatch(PROFILE_REQUEST);
    }

    var that = this;
    axios.interceptors.response.use(
      function(response) {
        // Do something with response data
        return response;
      },
      function(err) {
        if (err.response.status === 401) {
          // if you ever get an unauthorized, logout the user
          that.$store.dispatch(AUTH_ERROR);
          that.$router.push({ name: "login" });

          // you can also redirect to /login if needed !
        }
        return Promise.reject(err);
      }
    );
  },
  mounted() {
    WebFontLoader.load({
      google: {
        families: ["Roboto:100,300,400,500,700,900"]
      },
      active: this.setFontLoaded
    });
  },
  methods: {
    setFontLoaded() {
      this.$emit("font-loaded");
    }
  }
};
</script>

<style lang="scss">
#app {
  background: none;
  .background-carousel {
    position: absolute;
    z-index: -1;
  }
}
</style>
