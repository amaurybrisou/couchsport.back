<template>
  <v-container fill-height>
    <v-layout row justify-center align-center>
      <v-flex xs6>
        <v-card>
          <v-toolbar color="secondary">
            <h1 class="display-2 font-weight-thin">CoSport</h1>
          </v-toolbar>
          <v-card-text>
            <p>Lorem ipsum dolor sit amet consectetur adipisicing elit. Enim non eveniet labore! Rem cum mollitia odit ullam animi, deleniti nemo praesentium perferendis doloribus officia ratione molestiae, vel quaerat architecto unde?</p>
            <p>Lorem ipsum dolor sit amet consectetur adipisicing elit. Enim non eveniet labore! Rem cum mollitia odit ullam animi, deleniti nemo praesentium perferendis doloribus officia ratione molestiae, vel quaerat architecto unde?</p>
            <p>Lorem ipsum dolor sit amet consectetur adipisicing elit. Enim non eveniet labore! Rem cum mollitia odit ullam animi, deleniti nemo praesentium perferendis doloribus officia ratione molestiae, vel quaerat architecto unde?</p>
            <v-layout row wrap>
                <v-spacer></v-spacer>
                <v-btn flat to="/signin" >Sign In</v-btn>
            </v-layout>
          </v-card-text>
        </v-card>
      </v-flex>
    </v-layout>
  </v-container>
</template>

<script>
export default {
  name: "About"
};
</script>
