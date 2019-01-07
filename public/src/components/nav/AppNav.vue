<template>
  <nav class="app-nav">
    <v-system-bar color="primary">
      <!-- <v-toolbar-side-icon></v-toolbar-side-icon> -->
      <v-toolbar-items>
        <v-btn to="/" flat>
          <v-icon>home</v-icon>
        </v-btn>
        <v-btn to="/explore" flat>Explore</v-btn>
      </v-toolbar-items>
      <v-spacer></v-spacer>
      <v-toolbar-items class="hidden-sm-and-down">
        <v-btn v-if="isProfileLoaded" to="/profile" flat>Profile</v-btn>
        <v-btn to="/about" flat>About</v-btn>
        <v-btn v-if="!isAuthenticated" to="/signin" flat>Sign In</v-btn>
        <v-btn v-if="!isAuthenticated && !authLoading" to="/login" flat>Login</v-btn>
        <v-btn v-if="isAuthenticated" @click="logout" to="/logout" flat>Logout</v-btn>
      </v-toolbar-items>
    </v-system-bar>
  </nav>
</template>


<script>
import { mapGetters, mapState } from "vuex";
import { AUTH_LOGOUT } from "@/store/actions/auth";

export default {
  name: "AppNav",
  computed: {
    ...mapGetters(["getProfile", "isAuthenticated", "isProfileLoaded"]),
    ...mapState({
      authLoading: state => state.auth.status === "loading",
      name: state =>
        `${state.user.profile.Firstname} ${state.user.profile.Lastname}`
    })
  },
  mounted() {
    this.$root.$emit("navBarLoaded");
  },
  methods: {
    logout: function() {
      this.$store.dispatch(AUTH_LOGOUT).then(() => this.$router.push("/"));
    }
  }
};
</script>
