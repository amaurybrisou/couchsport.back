<template>
  <v-container fill-height>
    <!-- <v-layout class="frontend-content" row wrap align-center> -->
    <v-card class="frontend-content">
      <v-toolbar color="secondary" dense  >
        <v-toolbar-title>CoSport</v-toolbar-title>
      </v-toolbar>
      <v-card-text class="text-content">
        <h3 class="subheading">Travel, make friends, practicing your favorite sports together !</h3>
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
  background-color: rgba(255, 255, 255, .8) !important;
}
</style>
