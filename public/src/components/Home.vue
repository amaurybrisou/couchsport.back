<template>
  <v-container fluid>
    <v-layout>
      <v-flex xs12 sm4>
        <v-card>
          <v-toolbar color="secondary" dense>
            <v-toolbar-title>{{ $t("cosport") }}</v-toolbar-title>
          </v-toolbar>
          <v-card-text class="text-content">
            <h3 class="subheading">{{ $t("home.speech_title") }}</h3>
            <v-divider></v-divider>
            <div class="body-1 mt-2">{{ $t("home.speech_content") }}</div>
          </v-card-text>
          <!-- <v-flex class="text-xs-center"> -->
          <!-- <h1 class="display-3 font-weight-thin mb-3">CoSport</h1> -->
          <!-- </v-flex> -->
        </v-card>
      </v-flex>
    </v-layout>
  </v-container>
</template>


<script>
export default {
  name: "Home"
};
</script>

<style lang="scss">
.frontend-content {
  position: absolute;
  top: 10vh;
  background-color: rgba(255, 255, 255, 0.8) !important;
}
</style>
