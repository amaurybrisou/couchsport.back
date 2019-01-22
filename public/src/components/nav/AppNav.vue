<template>
  <nav>
    <v-system-bar color="primary" height="30vh">
      <!-- <v-toolbar-side-icon></v-toolbar-side-icon> -->
      <v-toolbar-items>
        <v-btn :to="{name: 'home'}" flat>
          <v-icon>home</v-icon>
        </v-btn>
        <v-btn :to="{ name : 'explore' }" flat>{{ $t("explore") }}</v-btn>
        <v-menu offset-y transition="slide-y-transition" style="z-index: 500;">
          <v-btn flat slot="activator" class="body-1 font-weight-regular">
            <v-icon class="mr-1">language</v-icon>
            {{ $i18n.locale.toUpperCase() }}
          </v-btn>
          <v-list>
            <v-list-tile v-for="(item, i) in languages" :key="i" @click="changeLocale(i)">
              <v-list-tile-title>{{item.title}}</v-list-tile-title>
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
        <span v-if="$vuetify.breakpoint.mdAndUp" style="height: 100%;">
          <v-btn
            v-if="!isAuthenticated && !isProfileLoaded"
            :to="{ name: 'about' }"
            flat
          >{{ $t("about") }}</v-btn>
          <v-btn
            v-if="!isAuthenticated"
            :to="{ name: 'signup' }"
            flat
          >{{ $t("signup") | capitalize }}</v-btn>
          <v-btn
            v-if="!isAuthenticated && !authLoading"
            :to="{ name: 'login' }"
            flat
          >{{ $t("login") | capitalize }}</v-btn>
        </span>
        <v-menu v-else open-on-hover offset-y transition="slide-y-transition" style="z-index: 500;">
          <v-btn icon slot="activator">
            <v-icon>account_box</v-icon>
          </v-btn>
          <v-list>
            <v-list-tile class="v-list__tile--link" :to="{ name : 'about'}">
              <v-list-tile-title>{{ $t('about') | capitalize }}</v-list-tile-title>
            </v-list-tile>
            <v-list-tile v-if="!isAuthenticated" :to="{ name: 'login'}" class="v-list__tile--link">
              <v-list-tile-title>{{ $t('login') | capitalize }}</v-list-tile-title>
            </v-list-tile>
            <v-list-tile v-if="!isAuthenticated" :to="{ name: 'signup'}" class="v-list__tile--link">
              <v-list-tile-title>{{ $t('signup') | capitalize }}</v-list-tile-title>
            </v-list-tile>
          </v-list>
        </v-menu>
      </v-toolbar-items>
    </v-system-bar>
  </nav>
</template>


<script>
import { mapGetters, mapState, mapMutations } from "vuex";
import { AUTH_LOGOUT } from "@/store/actions/auth";
import { MODIFY_PROFILE } from "@/store/actions/profile";

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
      languages: {
        en: { title: "English" },
        fr: { title: "Français" }
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
    this.MODIFY_PROFILE({ key: "locale", value: this.$i18n.locale });
  },
  methods: {
    ...mapMutations([MODIFY_PROFILE]),
    goToConversations() {
      this.$messenger.setMessagesRead();
      this.$router.push({ name: "profile", hash: "#conversations" });
    },
    changeLocale(locale) {
      this.$i18n.locale = locale;
      this.MODIFY_PROFILE({ key: "locale", value: locale });

      this.$router.push({ name: this.$route.name, params: { locale: locale } });
    },
    logout: function() {
      this.$store
        .dispatch(AUTH_LOGOUT)
        .then(() => this.$router.push({ name: "home" }));
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
