<template>
  <v-container fluild grid-list-xs>
    <v-layout v-if="isProfileLoaded" row wrap align-center justify-center>
      <v-flex
        :class="{ 'sm4 pr-5': $vuetify.breakpoint.smAndUp, 'xs12 pa-2': $vuetify.breakpoint.xsOnly }"
      >
        <upload-button @formData="handleImage">
          <v-card slot="appearance" flat tile class="d-flex profile-avatar">
            <v-img
              :src="Avatar"
              :alt="Username"
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
      <v-flex :class="{ 'sm4': $vuetify.breakpoint.smAndUp, 'xs12': $vuetify.breakpoint.xsOnly }">
        <v-text-field flat disabled readonly label="Email" v-model="email"></v-text-field>
        <v-text-field flat label="Username" v-model="Username"></v-text-field>
        <v-text-field flat label="Firstname" v-model="Firstname"></v-text-field>
        <v-text-field flat label="Lastname" v-model="Lastname"></v-text-field>
        <v-select :items="[`` ,`Male`, `Female`]" v-model="Gender" label="Gender"></v-select>
      </v-flex>
      <v-flex :class="{ 'sm4': $vuetify.breakpoint.smAndUp, 'xs12': $vuetify.breakpoint.xsOnly }">
        <v-text-field flat label="StreetName" v-model="StreetName"></v-text-field>
        <v-text-field flat label="City" v-model="City"></v-text-field>
        <v-text-field flat label="ZipCode" v-model="ZipCode"></v-text-field>
        <v-text-field flat label="Country" v-model="Country"></v-text-field>
        <v-text-field flat label="Phone" v-model="Phone"></v-text-field>
      </v-flex>
    </v-layout>
    <v-layout row wrap>
      <v-flex xs8 offset-xs2>
        <v-autocomplete
          v-model="Languages"
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
import { MODIFY_PROFILE, SAVE_PROFILE } from "@/store/actions/profile.js";

import UploadButton from "@/components/utils/UploadButton";
import AppSnackBar from "@/components/utils/AppSnackBar";

export default {
  name: "PersonalInformation",
  props: ["allLanguages"],
  components: { UploadButton, AppSnackBar },
  data() {
    return {
      snackbar: false,
      snackbarTimeout: 4000,
      snackbarText: "your profile has been successfully saved",
      showSavingProfileDialog: false
    };
  },
  computed: {
    ...mapGetters(["email", "isProfileLoaded"]),
    Avatar: {
      get() {
        return this.$store.state.profile.profile.Avatar;
      }
    },
    Username: {
      get() {
        return this.$store.state.profile.profile.Username;
      },
      set(v) {
        this.MODIFY_PROFILE({ key: "Username", value: v });
      }
    },
    Firstname: {
      get() {
        return this.$store.state.profile.profile.Firstname;
      },
      set(v) {
        this.MODIFY_PROFILE({ key: "Firstname", value: v });
      }
    },
    Lastname: {
      get() {
        return this.$store.state.profile.profile.Lastname;
      },
      set(v) {
        this.MODIFY_PROFILE({ key: "Lastname", value: v });
      }
    },
    Gender: {
      get() {
        return this.$store.state.profile.profile.Gender;
      },
      set(v) {
        this.MODIFY_PROFILE({ key: "Gender", value: v });
      }
    },
    StreetName: {
      get() {
        return this.$store.state.profile.profile.StreetName;
      },
      set(v) {
        this.MODIFY_PROFILE({ key: "StreetName", value: v });
      }
    },
    City: {
      get() {
        return this.$store.state.profile.profile.City;
      },
      set(v) {
        this.MODIFY_PROFILE({ key: "City", value: v });
      }
    },
    ZipCode: {
      get() {
        return this.$store.state.profile.profile.ZipCode;
      },
      set(v) {
        this.MODIFY_PROFILE({ key: "ZipCode", value: v });
      }
    },
    Country: {
      get() {
        return this.$store.state.profile.profile.Country;
      },
      set(v) {
        this.MODIFY_PROFILE({ key: "Country", value: v });
      }
    },
    Phone: {
      get() {
        return this.$store.state.profile.profile.Phone;
      },
      set(v) {
        this.MODIFY_PROFILE({ key: "Phone", value: v });
      }
    },
    Languages: {
      get() {
        return this.$store.state.profile.profile.Languages;
      },
      set(v) {
        this.MODIFY_PROFILE({ key: "Languages", value: v });
      }
    }
  },
  watch: {
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
          that.MODIFY_PROFILE({ key: "Avatar", value: e.target.result });
          that.MODIFY_PROFILE({ key: "AvatarFile", value: file.name });
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
