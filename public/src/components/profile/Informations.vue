<template>
  <v-container fluild grid-list-xs>
    <v-form v-if="isProfileLoaded" @keypress.enter.native="submit" ref="form" v-model="rules.valid">
      <v-layout row wrap align-center justify-center>
        <v-flex
          :class="{ 'sm4 pr-5': $vuetify.breakpoint.smAndUp, 'xs12 pa-2': $vuetify.breakpoint.xsOnly }"
        >
          <upload-button @formData="handleImage" :accept="rules.imageFormatsAllowed">
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
          <v-input :rules="rules['Avatar']" v-model="Avatar" hidden class="right"></v-input>
        </v-flex>
        <v-flex :class="{ 'sm4': $vuetify.breakpoint.smAndUp, 'xs12': $vuetify.breakpoint.xsOnly }">
          <v-text-field
            flat
            disabled
            readonly
            :label="$t('fields.email') | capitalize"
            v-model="email"
          ></v-text-field>
          <v-text-field
            flat
            autofocus
            :rules="rules['username']"
            :label="$t('fields.username') | capitalize"
            @keypress.enter.native="submit"
            v-model="Username"
          ></v-text-field>
          <v-text-field
            flat
            :rules="rules['firstname']"
            :label="$t('fields.firstname') | capitalize"
            @keypress.enter.native="submit"
            v-model="Firstname"
          ></v-text-field>
          <v-text-field
            flat
            :rules="rules['lastname']"
            :label="$t('fields.lastname') | capitalize"
            @keypress.enter.native="submit"
            v-model="Lastname"
          ></v-text-field>
          <v-select
            :items="[`` ,`Male`, `Female`]"
            v-model="Gender"
            :label="$t('fields.gender') | capitalize"
            @keypress.enter.native="submit"
          ></v-select>
        </v-flex>
        <v-flex :class="{ 'sm4': $vuetify.breakpoint.smAndUp, 'xs12': $vuetify.breakpoint.xsOnly }">
          <v-text-field
            flat
            :label="$t('fields.streetname') | capitalize"
            @keypress.enter.native="submit"
            v-model="StreetName"
          ></v-text-field>
          <v-text-field
            flat
            :label="$t('fields.city') | capitalize"
            @keypress.enter.native="submit"
            v-model="City"
          ></v-text-field>
          <v-text-field
            flat
            :label="$t('fields.zipcode') | capitalize"
            @keypress.enter.native="submit"
            v-model="ZipCode"
          ></v-text-field>
          <v-text-field
            flat
            :label="$t('fields.country') | capitalize"
            @keypress.enter.native="submit"
            v-model="Country"
          ></v-text-field>
          <v-text-field
            flat
            :rules="rules['phone']"
            :label="$t('fields.phone') | capitalize"
            @keypress.enter.native="submit"
            v-model="Phone"
          ></v-text-field>
        </v-flex>
      </v-layout>
      <v-layout row wrap>
        <v-flex xs8 offset-xs1>
          <v-autocomplete
            v-model="Languages"
            :items="allLanguages"
            :label="$t('languages') | capitalize"
            return-object
            item-text="Name"
            multiple
          ></v-autocomplete>
        </v-flex>
        <v-flex xs3>
          <v-btn color="warning" @click="showChangePasswordDialog = true">{{ $t('change_password') }}</v-btn>
        </v-flex>
      </v-layout>
      <v-layout row wrap>
        <v-btn color="success" :disabled="!rules.valid" @click="submit" block flat>{{ $t('save') }}</v-btn>
      </v-layout>
    </v-form>

    <!-- Warn section  used to display application state (saving and success) -->
    <app-snack-bar :state="snackbar" @snackClose="snackbar = false" :text="snackbarText"></app-snack-bar>
    <v-dialog v-model="showChangePasswordDialog" persistent width="400">
            <auth-form
              :title="$t('change_password') | capitalize"
              @submit="changePassword"
              @hideChangePasswordDialog="showChangePasswordDialog = false"
              :buttonMessage="$t('change_password') | capitalize"
            ></auth-form>
       </v-dialog>
    <v-dialog v-model="showSavingProfileDialog" hide-overlay persistent width="300">
      <v-card color="primary" dark>
        <v-card-text>
          {{ $t('message.stand_by') }}
          <v-progress-linear indeterminate color="white" class="mb-0"></v-progress-linear>
        </v-card-text>
      </v-card>
    </v-dialog>
  </v-container>
</template>


<script>
import { mapGetters, mapMutations, mapActions } from "vuex";
import {
  MODIFY_PROFILE,
  SAVE_PROFILE,
  GET_LANGUAGES,
} from "@/store/actions/profile.js";

import {  AUTH_CHANGE_PASSWORD} from "@/store/actions/auth.js";

import UploadButton from "@/components/utils/UploadButton";
import AppSnackBar from "@/components/utils/AppSnackBar";
import AuthForm from "@/components/auth/AuthForm";

export default {
  name: "PersonalInformation",
  components: { UploadButton, AppSnackBar, AuthForm },
  data() {
    return {
      snackbar: false,
      snackbarTimeout: 4000,
      snackbarText: "your profile has been successfully saved",
      
      showSavingProfileDialog: false,
      showChangePasswordDialog: false,

      rules: {
        valid: false,
        imageFormatsAllowed: "image/jpeg, image:jpg, image/png, image/gif",
        username: [
          v =>
            /^[àéèïîôoa-zA-Z0-9]{6,15}$/.test(v) ||
            this.$t("message.length_between", [
              this.$t("fields.username"),
              6,
              15
            ])
        ],
        firstname: [
          v =>
            /^[àéèïîôo a-zA-Z]{0,35}$/.test(v) ||
            this.$t("message.length_below", [this.$t("fields.firstname"), 35])
        ],
        lastname: [
          v =>
            /^[àéèïîaôö a-zA-Z]{0,35}$/.test(v) ||
            this.$t("message.length_below", [this.$t("fields.lastname"), 35])
        ],
        Avatar: [
          v => !!v || this.$t("message.required", ["e", this.$t("image")]),
          v =>
            /(?:png|jpg|jpeg|gif)$/i.test(v) ||
            /^\s*data:([a-z]+\/[a-z]+(;[a-z\-]+\=[a-z\-]+)?)?(;base64)?,[a-z0-9\!\$\&\'\,\(\)\*\+\,\;\=\-\.\_\~\:\@\/\?\%\s]*\s*$/i.test(
              v
            ) ||
            this.$t("message.invalid", [this.$t("image_link")])
        ]
      }
    };
  },
  computed: {
    ...mapGetters({
      email: "email",
      isProfileLoaded: "isProfileLoaded",
      allLanguages: "languages"
    }),
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
    ...mapActions([SAVE_PROFILE, GET_LANGUAGES, AUTH_CHANGE_PASSWORD]),
    ...mapMutations([MODIFY_PROFILE]),
    submit() {
      if (!this.$refs.form.validate()) {
        return;
      }
      this.showSavingProfileDialog = true;
      this.SAVE_PROFILE()
        .then(() => {
          this.showSavingProfileDialog = false;
          this.snackbarText = this.$t("message.success_saving", ["profile"]);
          this.snackbar = true;
        })
        .catch(e => {
          this.showSavingProfileDialog = false;
          this.snackbarText = this.$t("message.error_saving", ["profile"]);
          this.snackbar = true;
        });
    },
    changePassword(user){
      this.showSavingProfileDialog = true;
      this[AUTH_CHANGE_PASSWORD](user).then(() => {
        this.snackbarText = this.$t("message.success_updating", [this.$t("password")])
        this.snackbar = true;
        this.showChangePasswordDialog = false;
        this.showSavingProfileDialog = false;
      })
      .catch(() => {
        this.showSavingProfileDialog = false;
        this.snackbarText = this.$t("message.error_updating", [this.$t("password")])
        this.snackbar = true;
      })
    },
    handleImage(formData) {
      var file = formData.get("file");
      if (file instanceof File && file.size) {
        if (file.size > 100000) {
          this.snackbarText = this.$t("message.too_big", ["image", "100ko"]);
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
