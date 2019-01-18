<template>
  <nav>
    <v-system-bar color="primary" height="30vh">
      <!-- <v-toolbar-side-icon></v-toolbar-side-icon> -->
      <v-toolbar-items>
        <v-btn to="/" flat>
          <v-icon>home</v-icon>
        </v-btn>
        <v-btn :to="{ name : 'explore' }" flat>{{ $t("explore") }}</v-btn>
        <v-menu offset-y transition="slide-y-transition" style="z-index: 500;">
          <v-btn flat slot="activator">
            <flag :iso="languagesFlags[$i18n.locale].flag" v-bind:squared="false"/>
            <!-- {{ $i18n.locale }} -->
          </v-btn>
          <v-list>
            <v-list-tile v-for="(item, i) in languagesFlags" :key="i" @click="changeLocale(i)">
              <v-list-tile-title>
                <flag :iso="item.flag" v-bind:squared="false"/>
                {{item.title}}
              </v-list-tile-title>
            </v-list-tile>
          </v-list>
        </v-menu>
      </v-toolbar-items>
      <v-spacer></v-spacer>
      <v-toolbar-items>
        <a
          class="new-message-link"
          v-if="isProfileLoaded && unread_message && $route.hash !== '#conversations'"
          @click="goToConversations"
        >
          <v-chip small color="info">
            {{unread_message}}
            <v-icon>mail_outline</v-icon>
          </v-chip>
        </a>
        <v-menu
          v-if="isProfileLoaded"
          open-on-hover
          offset-y
          transition="slide-y-transition"
          style="z-index: 500;"
        >
          <v-btn icon slot="activator">
            <v-icon>account_box</v-icon>
          </v-btn>
          <v-list>
            <v-list-tile
              v-for="link in links"
              :key="link.name"
              :to="{ name : 'profile', hash : link.to }"
            >
              <v-list-tile-title>{{ $t(`${link.name}`) | capitalize }}</v-list-tile-title>
            </v-list-tile>
            <v-list-tile class="v-list__tile--link" :to="{ name : 'about'}">
              <v-list-tile-title>{{ $t('about') | capitalize }}</v-list-tile-title>
            </v-list-tile>
            <v-list-tile v-if="isAuthenticated" class="v-list__tile--link">
              <v-list-tile-title @click="logout">{{ $t('logout') | capitalize }}</v-list-tile-title>
            </v-list-tile>
          </v-list>
        </v-menu>
        <v-btn
          v-if="!isAuthenticated && !isProfileLoaded"
          :to="{ name: 'about' }"
          flat
        >{{ $t("about") }}</v-btn>
        <v-btn v-if="!isAuthenticated" :to="{ name: 'signup' }" flat>{{ $t("signup") | capitalize }}</v-btn>
        <v-btn
          v-if="!isAuthenticated && !authLoading"
          :to="{ name: 'login' }"
          flat
        >{{ $t("login") | capitalize }}</v-btn>
      </v-toolbar-items>
    </v-system-bar>
  </nav>
</template>


<script>
import { mapGetters, mapState, mapMutations } from "vuex";
import { AUTH_LOGOUT } from "@/store/actions/auth";

const NAMESPACE = "conversations/";

export default {
  name: "AppNav",
  data() {
    return {
      links: [
        { auth: true, to: "#informations", name: "profile" },
        { auth: true, to: "#activities", name: "activities" },
        { auth: true, to: "#conversations", name: "conversations" },
        { auth: true, to: "#pages", name: "pages" }
      ],
      languagesFlags: {
        en: { flag: "gb", title: "English" },
        fr: { flag: "fr", title: "Français" }
      }
    };
  },
  computed: {
    ...mapGetters(["getProfile", "isAuthenticated", "isProfileLoaded"]),
    ...mapState({
      unread_message: state => state.profile.conversations.unread,
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
      this.$messenger.setMessagesRead();
      this.$router.push({ name: "profile", hash: "#conversations" });
    },
    changeLocale(locale) {
      this.$i18n.locale = locale;
      this.$router.push({ name: this.$route.name, params: { locale: locale } });
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
