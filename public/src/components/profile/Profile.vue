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
          <v-tab to="#informations" href="#informations" class="subheader">Personal Information</v-tab>
          <v-tab to="#activities" href="#activities" class="subheader">Activities</v-tab>
          <v-tab to="#conversations" href="#conversations" class="subheader" @click="$store.commit('conversations/MESSAGES_READ')">Conversations</v-tab>
          <v-tab to="#pages" href="#pages" class="subheader">My Spots</v-tab>
          <v-tabs-items v-model="activeTab">
            <v-tab-item value="informations">
              <v-card flat>
                <informations :allLanguages="allLanguages"></informations>
              </v-card>
            </v-tab-item>
            <v-tab-item lazy value="activities">
              <v-card flat>
                <activities
                  :activities="getProfile.Activities || []"
                  :allActivities="allActivities"
                ></activities>
              </v-card>
            </v-tab-item>
            <v-tab-item lazy value="conversations">
              <v-card flat>
                <conversations></conversations>
              </v-card>
            </v-tab-item>
            <v-tab-item lazy value="pages">
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

import { GET_CONVERSATIONS } from "@/store/actions/conversations";
import { mapGetters, mapActions } from "vuex";

import activityRepo from "@/repositories/activity.js";
import languageRepo from "../../repositories/language.js";

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
  asyncComputed: {
    async allActivities() {
      return await activityRepo.all().then(({ data }) => data);
    },
    async allLanguages() {
      return await languageRepo.all().then(({ data }) => data);
    }
  }
};
</script>
