<template>
  <nav>
    <v-system-bar color="primary" height="30vh">
      <!-- <v-toolbar-side-icon></v-toolbar-side-icon> -->
      <v-toolbar-items>
        <v-btn to="/" flat>
          <v-icon>home</v-icon>
        </v-btn>
        <v-btn to="/explore" flat>Explore</v-btn>
      </v-toolbar-items>
      <v-spacer></v-spacer>
      <v-toolbar-items>
        <a class="new-message-link" v-if="messages.length" @click="goToConversations">
          <v-chip small color="info">
            {{messages.length}}
            <v-icon>mail_outline</v-icon>
          </v-chip>
        </a>
        <v-menu
          v-if="isProfileLoaded"
          open-on-hover
          offset-y
          origin="left top"
          transition="scale-transition"
        >
          <v-btn icon slot="activator">
            <v-icon>account_box</v-icon>
          </v-btn>
          <v-list>
            <v-list-tile v-for="link in links" :key="link.name" :to="link.to">
              <v-list-tile-title>{{ link.name }}</v-list-tile-title>
            </v-list-tile>
            <v-list-tile v-if="isAuthenticated" class="v-list__tile--link">
              <v-list-tile-title @click="logout">Logout</v-list-tile-title>
            </v-list-tile>
          </v-list>
        </v-menu>
        <v-btn v-if="!isAuthenticated" to="/signup" flat>Sign Up</v-btn>
        <v-btn v-if="!isAuthenticated && !authLoading" to="/login" flat>Login</v-btn>
      </v-toolbar-items>
    </v-system-bar>
  </nav>
</template>


<script>
import { mapGetters, mapState } from "vuex";
import { AUTH_LOGOUT } from "@/store/actions/auth";
import { MESSAGES_READ } from "@/store/actions/ws";

export default {
  name: "AppNav",
  data() {
    return {
      links: [
        { to: "/profile#informations", name: "Profile" },
        { to: "/profile#activities", name: "Activities" },
        { to: "/profile#conversations", name: "Conversations" },
        { to: "/profile#pages", name: "Pages" },
        { to: "/about", name: "About" }
      ]
    };
  },
  computed: {
    ...mapGetters(["getProfile", "isAuthenticated", "isProfileLoaded"]),
    ...mapState({
      messages: state => state.ws.messages,
      authLoading: state => state.auth.status === "loading",
      name: state =>
        `${state.user.profile.Firstname} ${state.user.profile.Lastname}`
    })
  },
  mounted() {
    this.$root.$emit("navBarLoaded");
  },
  methods: {
    goToConversations() {
      this.$store.commit(MESSAGES_READ);
      this.$router.push("/profile#conversations");
    },
    logout: function() {
      this.$store.dispatch(AUTH_LOGOUT).then(() => this.$router.push("/"));
    }
  }
};
</script>

<style lang="scss">
.new-message-link {
  text-decoration: none;
  color: none;
  outline: none;
  span:focus:after {
    background: none !important;
  }
}
</style>
