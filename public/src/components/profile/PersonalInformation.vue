<template>
  <v-container class="first-profile-step" fluild grid-list-md>
    <v-layout align-center justify-center>
      <v-flex xs4 pa-5>
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
      <v-flex xs12 sm6 md4>
        <v-text-field flat disabled readonly label="Email" v-model="email"></v-text-field>
        <v-text-field flat label="Username" v-model="Username"></v-text-field>
        <v-text-field flat label="Firstname" v-model="Firstname"></v-text-field>
        <v-text-field flat label="Lastname" v-model="Lastname"></v-text-field>
        <v-select :items="[`` ,`Male`, `Female`]" v-model="Gender" label="Gender"></v-select>
      </v-flex>
      <v-flex xs12 sm8 md4>
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
          :items="languages"
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
    <app-snack-bar
      :state="snackbar"
      @snackClose="snackbar = false"
      @snackOpen="setTimeout"
      :text="snackbarText"
    ></app-snack-bar>
  </v-container>
</template>


<script>
import { mapState, mapGetters } from "vuex";
import { MODIFY_PROFILE, SAVE_PROFILE } from "@/store/actions/user.js";

import UploadButton from "@/components/utils/UploadButton";
import AppSnackBar from "@/components/utils/AppSnackBar";

import profileRepo from "@/repositories/profile.js";
import languageRepo from "../../repositories/language.js";

export default {
  name: "PersonalInformation",
  components: { UploadButton, AppSnackBar },
  data() {
    return {
      snackbar: false,
      snackbarTimeout: 4000,
      snackbarText: "your profile has been successfully saved"
    };
  },
  methods: {
    setTimeout() {
      var that = this;
      setTimeout(function() {
        that.snackbar = false;
      }, that.snackbarTimeout);
    },
    submit() {
      var that = this;
      this.$store
        .dispatch(SAVE_PROFILE)
        .then(() => {
          that.snackbar = true;
        })
        .catch(e => {
          that.snackbarText = "there was and error saving your profile";
          that.snackbar = true;
        });
    },
    handleImage(formData) {
      var file = formData.get("file");
      if (file instanceof File && file.size) {
        if(file.size > 100000) {
          this.snackbarText = "This image is too big";
          return this.snackbar = true;
        }
        var that = this;
        var reader = new FileReader();
        reader.onload = function(e) {
          that.$store.dispatch(MODIFY_PROFILE, {
            Avatar: e.target.result,
            AvatarFile: file.name
          });
        };

        reader.readAsDataURL(file);
      }
    }
  },
  asyncComputed: {
    languages: async function() {
      return await languageRepo.all().then(({ data }) => data);
    }
  },
  computed: {
    ...mapState({
      profile: state => state.user.profile,
      Avatar: state => state.user.profile.Avatar
    }),
    ...mapGetters(["email"]),
    Username: {
      get() {
        return this.$store.state.user.profile.Username;
      },
      set(val) {
        this.$store.dispatch(MODIFY_PROFILE, {
          Username: val
        });
      }
    },
    Firstname: {
      get() {
        return this.$store.state.user.profile.Firstname;
      },
      set(val) {
        this.$store.dispatch(MODIFY_PROFILE, {
          Firstname: val
        });
      }
    },
    Lastname: {
      get() {
        return this.$store.state.user.profile.Lastname;
      },
      set(val) {
        this.$store.dispatch(MODIFY_PROFILE, {
          Lastname: val
        });
      }
    },
    Country: {
      get() {
        return this.$store.state.user.profile.Country;
      },
      set(val) {
        this.$store.dispatch(MODIFY_PROFILE, {
          Country: val
        });
      }
    },
    Languages: {
      get() {
        return this.$store.state.user.profile.Languages;
      },
      set(val) {
        this.$store.dispatch(MODIFY_PROFILE, {
          Languages: val
        });
      }
    },
    Phone: {
      get() {
        return this.$store.state.user.profile.Phone;
      },
      set(val) {
        this.$store.dispatch(MODIFY_PROFILE, {
          Phone: val
        });
      }
    },
    ZipCode: {
      get() {
        return this.$store.state.user.profile.ZipCode;
      },
      set(val) {
        this.$store.dispatch(MODIFY_PROFILE, {
          ZipCode: val
        });
      }
    },
    City: {
      get() {
        return this.$store.state.user.profile.City;
      },
      set(val) {
        this.$store.dispatch(MODIFY_PROFILE, {
          City: val
        });
      }
    },
    Gender: {
      get() {
        return this.$store.state.user.profile.Gender;
      },
      set(val) {
        this.$store.dispatch(MODIFY_PROFILE, {
          Gender: val
        });
      }
    },
    StreetName: {
      get() {
        return this.$store.state.user.profile.StreetName;
      },
      set(val) {
        this.$store.dispatch(MODIFY_PROFILE, {
          StreetName: val
        });
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
