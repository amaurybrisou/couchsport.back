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
          <v-tab to="#tab-1" href="#tab-1" class="subheader">Personal Information</v-tab>
          <v-tab to="#tab-2" href="#tab-2" class="subheader">Activities</v-tab>
          <v-tab to="#tab-3" href="#tab-3" class="subheader">Conversations</v-tab>
          <v-tab to="#tab-4" href="#tab-4" class="subheader">My Spots</v-tab>
          <v-tabs-items v-model="activeTab">
            <v-tab-item value="tab-1">
              <v-card flat>
                <informations :profile="getProfile" :allLanguages="allLanguages"></informations>
              </v-card>
            </v-tab-item>
            <v-tab-item value="tab-2">
              <v-card flat>
                <activities :activities="getProfile.Activities || []" :allActivities="allActivities"></activities>
              </v-card>
            </v-tab-item>
            <v-tab-item value="tab-3">
              <v-card flat>
                <conversations v-if="conversations" :conversations="conversations || []"></conversations>
              </v-card>
            </v-tab-item>
            <v-tab-item value="tab-4">
              <v-card flat>
                <pages :pages="getProfile.OwnedPages || []" :allActivities="allActivities"></pages>
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
import { mapGetters } from "vuex";

import activityRepo from "@/repositories/activity.js";
import languageRepo from "../../repositories/language.js";
import conversationRepo from "../../repositories/conversation.js";

export default {
  name: "Profile",
  components: {
    Informations,
    Activities,
    Pages,
    Conversations,
  },
  data() {
    return {
      activeTab: "tab-1",
    };
  },
  computed: {
    ...mapGetters(['isProfileLoaded', 'getProfile'])
  },
  asyncComputed: {
    async allActivities() {
      return await activityRepo.all().then(({ data }) => data);
    },
    async allLanguages() {
      return await languageRepo.all().then(({ data }) => data);
    },
    async conversations() {
      return await conversationRepo.mines().then(({data}) => data)
    }
  },
}
</script>
