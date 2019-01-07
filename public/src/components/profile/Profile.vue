<template>
  <v-container id="profile" fill-height>
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
          <v-tab href="#tab-1" class="subheader">Personal Information</v-tab>
          <v-tab href="#tab-2" class="subheader">Activities</v-tab>
          <v-tab href="#tab-3" class="subheader">My Spots</v-tab>
          <v-tabs-items v-model="activeTab">
            <v-tab-item value="tab-1">
              <v-card flat>
                <personal-information></personal-information>
              </v-card>
            </v-tab-item>
            <v-tab-item value="tab-2">
              <v-card flat>
                <personal-activities></personal-activities>
              </v-card>
            </v-tab-item>
            <v-tab-item value="tab-3">
              <v-card flat>
                <personal-pages></personal-pages>
              </v-card>
            </v-tab-item>
          </v-tabs-items>
        </v-tabs>
      </v-flex>
    </v-layout>
  </v-container>
</template>

<script>

import PersonalInformation from "./PersonalInformation";
import PersonalActivities from "./PersonalActivities";
import PersonalPages from "./PersonalPages";

export default {
  name: "Profile",
  components: {
    PersonalInformation,
    PersonalActivities,
    PersonalPages,
  },
  data() {
    return {
      activeTab: "tab-1",
    };
  }
}
</script>
