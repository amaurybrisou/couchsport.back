<template>
  <v-container fluid fill-height>
    <v-layout justify-center align-center>
      <v-flex xs12 sm8 md4>
        <auth-form :title="$t('login') | capitalize" @submit="submit" :buttonMessage="$t('login') | capitalize" :welcome="welcome" :errors="errors"></auth-form>
      </v-flex>
    </v-layout>
  </v-container>
</template>

<script>
import { AUTH_REQUEST } from "@/store/actions/auth";
import { mapActions } from "vuex";
import AuthForm from "./AuthForm";

export default {
  name: "Login",
  components: { AuthForm },
  props: { welcome: null },
  data() {
    return { errors: [] };
  },
  methods: {
    ...mapActions([AUTH_REQUEST]),
    async submit(user) {
      this.AUTH_REQUEST(user)
        .then(() => {
          this.$router.push({ name: "profile" });
        })
        .catch(data => {
          this.errors = [];
          this.errors.push(data);
        });
    }
  }
};
</script>

