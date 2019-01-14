<template>
  <v-container class="first-profile-step" fluild grid-list-md>
    <v-layout align-center justify-center>
      <v-flex xs4 pa-5>
        <upload-button @formData="handleImage">
          <v-card slot="appearance" flat tile class="d-flex profile-avatar">
            <v-img
              :src="Avatar"
              :alt="local_profile.Username"
              :value="Avatar"
              aspect-ratio="1"
              class="grey lighten-2"
            >
              <v-layout slot="placeholder" fill-height align-center justify-center ma-0>
                <!-- <v-progress-circular indeterminate color="grey lighten-5"></v-progress-circular> -->
                <v-icon large>person</v-icon>
              </v-layout>
            </v-img>
          </v-card>
        </upload-button>
      </v-flex>
      <v-flex xs12 sm6 md4>
        <v-text-field flat disabled readonly label="Email" v-model="email"></v-text-field>
        <v-text-field flat label="Username" v-model="local_profile.Username"></v-text-field>
        <v-text-field flat label="Firstname" v-model="local_profile.Firstname"></v-text-field>
        <v-text-field flat label="Lastname" v-model="local_profile.Lastname"></v-text-field>
        <v-select :items="[`` ,`Male`, `Female`]" v-model="local_profile.Gender" label="Gender"></v-select>
      </v-flex>
      <v-flex xs12 sm8 md4>
        <v-text-field flat label="StreetName" v-model="local_profile.StreetName"></v-text-field>
        <v-text-field flat label="City" v-model="local_profile.City"></v-text-field>
        <v-text-field flat label="ZipCode" v-model="local_profile.ZipCode"></v-text-field>
        <v-text-field flat label="Country" v-model="local_profile.Country"></v-text-field>
        <v-text-field flat label="Phone" v-model="local_profile.Phone"></v-text-field>
      </v-flex>
    </v-layout>
    <v-layout row wrap>
      <v-flex xs8 offset-xs2>
        <v-autocomplete
          v-model="local_profile.Languages"
          :items="allLanguages"
          label="Languages"
          return-object
          item-text="Name"
          multiple
        ></v-autocomplete>
      </v-flex>
    </v-layout>
    <v-layout row wrap>
      <v-btn color="success" @click="submit" block flat>Save</v-btn>
    </v-layout>

    <!-- Warn section  used to display application state (saving and success) -->
    <app-snack-bar :state="snackbar" @snackClose="snackbar = false" :text="snackbarText"></app-snack-bar>
    <v-dialog v-model="showSavingProfileDialog" hide-overlay persistent width="300">
      <v-card color="primary" dark>
        <v-card-text>Please stand by
          <v-progress-linear indeterminate color="white" class="mb-0"></v-progress-linear>
        </v-card-text>
      </v-card>
    </v-dialog>
  </v-container>
</template>


<script>
import { mapGetters, mapMutations, mapActions } from "vuex";
import { MODIFY_PROFILE, SAVE_PROFILE } from "@/store/actions/user.js";

import UploadButton from "@/components/utils/UploadButton";
import AppSnackBar from "@/components/utils/AppSnackBar";

export default {
  name: "PersonalInformation",
  props: ["profile", "allLanguages"],
  components: { UploadButton, AppSnackBar },
  data() {
    return {
      local_profile: { ...{}, ...this.profile },
      snackbar: false,
      snackbarTimeout: 4000,
      snackbarText: "your profile has been successfully saved",
      showSavingProfileDialog: false
    };
  },
  computed: {
    ...mapGetters(["email"]),
    Avatar: {
      get() {
        return this.$store.state.user.profile.Avatar;
      }
    }
  },
  watch: {
    local_profile: {
      handler: function(v) {
        this.MODIFY_PROFILE(v);
      },
      deep: true
    },
    snackbar(v) {
      if (!v) return;
      var that = this;
      setTimeout(function() {
        that.snackbar = false;
      }, that.snackbarTimeout);
    }
  },
  methods: {
    ...mapActions([SAVE_PROFILE]),
    ...mapMutations([MODIFY_PROFILE]),
    submit() {
      this.showSavingProfileDialog = true;
      var that = this;
      this.SAVE_PROFILE()
        .then(() => {
          this.showSavingProfileDialog = false;
          this.snackbarText = "your profile has been successfully saved";
          that.snackbar = true;
        })
        .catch(e => {
          this.showSavingProfileDialog = false;
          that.snackbarText = "there was and error saving your profile";
          that.snackbar = true;
        });
    },
    handleImage(formData) {
      var file = formData.get("file");
      if (file instanceof File && file.size) {
        if (file.size > 100000) {
          this.snackbarText = "This image is too big";
          return (this.snackbar = true);
        }
        var that = this;
        var reader = new FileReader();
        reader.onload = function(e) {
          that.MODIFY_PROFILE({
            Avatar: e.target.result,
            AvatarFile: file.name
          });
        };

        reader.readAsDataURL(file);
      }
    }
  }
};
</script>


<style lang="scss">
.profile-avatar {
  border-radius: 50%;
}
</style>
