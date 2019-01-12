<template>
  <v-container fluid fill-height>
    <v-layout justify-center align-center>
      <v-flex xs12 sm8 md4>
        <v-card>
          <v-toolbar dark color="secondary">
            <v-toolbar-title>Login</v-toolbar-title>
          </v-toolbar>
          <v-card-text>
            <div v-if="errors.length" color="error">
              <v-alert v-for="(err, i) in errors" :key="i" :value="err" type="error">{{err}}</v-alert>
            </div>
            <div v-if="welcome" color="info">
              <v-alert type="info" :value="welcomeMessage">{{ welcomeMessage }}</v-alert>
            </div>

            <v-form @keypress.enter.native="submit" ref="form" v-model="valid">
              <v-text-field
                label="Email"
                type="text"
                v-model="user.email"
                name="email"
                :rules="emailRules"
                autocomplete="email"
              ></v-text-field>
              <v-text-field
                label="Password"
                :type="'password'"
                v-model="user.password"
                name="password"
                counter="8"
                :rules="passwordRules"
                autocomplete="current-password"
              ></v-text-field>
            </v-form>
          </v-card-text>
          <v-card-actions>
            <v-spacer></v-spacer>
            <v-btn @click="submit" :disabled="!valid">Login</v-btn>
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

import { AUTH_REQUEST } from "@/store/actions/auth";

export default {
  name: "SignIn",
  mixins: [validationMixin],
  props: { welcome: { type: String, default: null } },
  validations: {
    password: { required, minLength: maxLength(8) },
    email: { required, email }
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
        v => !!v || "E-mail is required",
        v =>
          /^\w+([\.-]?\w+)*@\w+([\.-]?\w+)*(\.\w{2,3})+$/.test(v) ||
          "E-mail must be valid"
      ],
      passwordRules: [
        v => !!v || "Password is required",
        v =>
          /^(?=.*\d)(?=.*[a-z])(?=.*[A-Z])[0-9a-zA-Z]{8,}$/.test(v) ||
          "Password must be valid"
      ]
    };
  },
  methods: {
    async submit(e) {
      if (!this.valid) return;
      this.$store
        .dispatch(AUTH_REQUEST, this.user)
        .then(() => {
          this.$router.push("/profile");
        })
        .catch((data) => {
          this.errors = [];
          this.errors.push(data);
        });
    }
  },
  computed: {
    welcomeMessage() {
      return this.welcome;
    },
    emailErrors() {
      const errors = [];
      if (!this.$v.user.email.$dirty) return errors;
      !this.$v.user.email.email && errors.push("Must be valid e-mail");
      !this.$v.user.email.required && errors.push("E-mail is required");
      return errors;
    },
    passwordErrors() {
      const errors = [];
      if (!this.$v.user.password.$dirty) return errors;
      !this.$v.user.passwor.maxLength &&
        errors.push("Password must be at most 8 characters long");
      !this.$v.user.password.required && errors.push("Password is required.");
      return errors;
    }
  }
};
</script>

