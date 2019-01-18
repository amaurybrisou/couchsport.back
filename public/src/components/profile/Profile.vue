<template>
  <v-container v-if="isProfileLoaded" id="profile" fill-height>
    <v-layout justify-center fill-height>
      <v-flex xs12 xm6>
        <v-tabs
          slot="extension"
          v-model="activeTab"
          slider-color="secondary"
          color="primary"
          centered
          show-arrows
        >
          <v-tab
            to="#informations"
            href="#informations"
            class="subheader"
          >{{ $t('personal_informations') }}</v-tab>
          <v-tab to="#activities" href="#activities" class="subheader">{{ $t('activities') }}</v-tab>
          <v-tab
            to="#conversations"
            href="#conversations"
            class="subheader"
            @click="$messenger.setMessagesRead()"
          >{{ $t('conversations') }}</v-tab>
          <v-tab to="#pages" href="#pages" class="subheader">{{ $t('pages') }}</v-tab>
          <v-tabs-items v-model="activeTab">
            <v-tab-item value="informations">
              <v-card flat>
                <informations></informations>
              </v-card>
            </v-tab-item>
            <v-tab-item lazy value="activities">
              <v-card flat>
                <activities></activities>
              </v-card>
            </v-tab-item>
            <v-tab-item lazy value="conversations">
              <v-card flat>
                <conversations></conversations>
              </v-card>
            </v-tab-item>
            <v-tab-item lazy value="pages">
              <v-card flat>
                <pages></pages>
              </v-card>
            </v-tab-item>
          </v-tabs-items>
        </v-tabs>
      </v-flex>
    </v-layout>
  </v-container>
</template>

<script>
import Informations from "./Informations";
import Activities from "./Activities";
import Pages from "./Pages";
import Conversations from "./Conversations";
import { mapGetters, mapActions } from "vuex";
import { GET_ACTIVITIES, GET_LANGUAGES } from "@/store/actions/profile";

export default {
  name: "Profile",
  components: {
    Informations,
    Activities,
    Pages,
    Conversations
  },
  data() {
    return {
      activeTab: "informations"
    };
  },
  computed: {
    ...mapGetters(["isProfileLoaded", "getProfile"])
  },
  methods: {
    ...mapActions([GET_ACTIVITIES, GET_LANGUAGES])
  },
  mounted() {
    this.GET_ACTIVITIES(this.$route.params.locale, this.$route.params.locale);
    this.GET_LANGUAGES(this.$route.params.locale, this.$route.params.locale);
  }
};
</script>
