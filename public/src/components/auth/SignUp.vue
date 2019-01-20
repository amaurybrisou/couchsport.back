<template>
  <v-container fluid fill-height>
    <v-layout justify-center align-center>
      <v-flex xs12 sm8 md4>
        <auth-form :title="$t('signup')" @submit="submit"  :buttonMessage="$t('signup')" :errors="errors"></auth-form>
      </v-flex>
    </v-layout>
  </v-container>
</template>

<script>
import { repositoryFactory } from "@/repositories/repositoryFactory";
const userRepository = repositoryFactory.get("user");

import AuthForm from './AuthForm';

export default {
  name: "SignUp",
  components: {AuthForm},
  data() {
    return { errors: [] }
  },
  methods: {
    async submit(user) {
      await userRepository
        .create(user)
        .then(({ data }) => {
          this.$router.push({
            name: "login",
            params: { welcome: this.$t("message.signup_success_welcome") }
          });
        })
        .catch(({ response: { data } }) => {
          this.errors = [];
          this.errors.push(data);
        });
    }
  },
};
</script>

