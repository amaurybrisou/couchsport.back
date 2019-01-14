<template>
  <v-container fill-height>
    <!-- <v-layout class="frontend-content" row wrap align-center> -->
    <v-card class="frontend-content">
      <v-toolbar color="secondary" dense>
        <v-toolbar-title>CoSport</v-toolbar-title>
      </v-toolbar>
      <v-card-text class="text-content">
        <h3 class="subheading">Travel, make friends, practicing your favorite sports together !</h3>
        <v-divider></v-divider>
        <div class="body-1 mt-2">Begin the tour browsing available spots in the
          <router-link to="/explore">
            <span color="accent">explore</span>
          </router-link>page or create your
          <router-link to="/signup">
            <span color="accent">account</span>
          </router-link>and promote your own spot!
        </div>
      </v-card-text>
      <!-- <v-flex class="text-xs-center"> -->
      <!-- <h1 class="display-3 font-weight-thin mb-3">CoSport</h1> -->
      <!-- </v-flex> -->
    </v-card>
    <!-- </v-layout> -->
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
  left: 10vh;
  background-color: rgba(255, 255, 255, 0.8) !important;
}
</style>
