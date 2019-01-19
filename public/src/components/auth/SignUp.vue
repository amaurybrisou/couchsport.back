<template>
  <v-container fluid fill-height>
    <v-layout justify-center align-center>
      <v-flex xs12 sm8 md4>
        <v-card>
          <v-toolbar dark color="secondary">
            <v-toolbar-title>{{ $t('signup') }}</v-toolbar-title>
          </v-toolbar>
          <v-card-text>
            <div v-if="errors.length" color="error">
              <v-alert v-for="(err, i) in errors" :key="i" :value="err" type="error">{{err}}</v-alert>
            </div>
            <v-form @keypress.enter.native="submit" ref="form" v-model="valid">
              <v-text-field
                :label="$t('email')"
                type="text"
                v-model="user.email"
                name="email"
                :rules="emailRules"
                autocomplete="email"
              ></v-text-field>
              <v-text-field
                :label="$t('password')"
                :counter="8"
                :type="'password'"
                v-model="user.password"
                name="password"
                :rules="passwordRules"
                autocomplete="current-password"
              ></v-text-field>
            </v-form>
          </v-card-text>
          <v-card-actions>
            <v-spacer></v-spacer>
            <v-btn @click="submit" :disabled="!valid">{{ $t('signup') }}</v-btn>
          </v-card-actions>
        </v-card>
      </v-flex>
    </v-layout>
  </v-container>
</template>

<script>
import { repositoryFactory } from "@/repositories/repositoryFactory";
const userRepository = repositoryFactory.get("user");

import { validationMixin } from "vuelidate";
import { required, maxLength, email } from "vuelidate/lib/validators";

export default {
  name: "SignUp",
  mixins: [validationMixin],

  validations: {
    email: { required, email },
    password: { required }
  },
  data() {
    return {
      valid: false,
      errors: [],
      user: {
        email: "",
        password: ""
      },
      emailRules: [
        v => !!v || this.$t("message.required", ["", this.$t("email")]),
        v =>
          /^\w+([\.-]?\w+)*@\w+([\.-]?\w+)*(\.\w{2,3})+$/.test(v) ||
          this.$t("message.invalid", [this.$t("email")])
      ],
      passwordRules: [
        v => !!v || this.$t("message.required", ["", this.$t("password")]),
        v =>
          /^(?=.*\d)(?=.*[_!?,]?)(?=.*[a-z])(?=.*[A-Z])[0-9a-zA-Z_?,!]{8,}$/.test(v) ||
          this.$t("message.invalid", [this.$t("password")], [8])
      ]
    };
  },
  methods: {
    async submit(e) {
      if (!this.valid) return;
      var that = this;
      await userRepository
        .create(this.user)
        .then(({ data }) => {
          that.$router.push({
            name: "login",
            params: { welcome: this.$t("message.signup_success_welcome") }
          });
        })
        .catch(({ response: { data } }) => {
          that.errors = [];
          that.errors.push(data);
        });
    }
  },
  computed: {
    emailErrors() {
      const errors = [];
      if (!this.$v.user.email.$dirty) return errors;
      !this.$v.user.email.email &&
        errors.push(this.$t("message.invalid", [this.$t("email")]));
      !this.$v.user.email.required &&
        errors.push(this.$t("message.required", ["", this.$t('email')]));
      return errors;
    },
    passwordErrors() {
      const errors = [];
      if (!this.$v.user.password.$dirty) return errors;
      !this.$v.user.passwor.maxLength &&
        errors.push(this.$t("message.password_hint", [8]));
      !this.$v.user.password.required &&
        errors.push(this.$t("message.required", ["", this.$t('password')]));
      return errors;
    }
  }
};
</script>

