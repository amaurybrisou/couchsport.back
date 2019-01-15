
<template>
  <v-container fluild grid-list-xs>
    <v-layout align-center pb-3>
      <div>Select the activities you're looking for.</div>
    </v-layout>
    <v-divider></v-divider>
    <v-layout v-if="allActivities" wrap align-center justify-center>
      <v-flex xs6 md2 v-for="(item, i) in allActivities" :key="i">
        <v-checkbox
          height="0"
          :label="item.Name"
          :value="item"
          multiple
          v-model="selected_activities"
        >{{ item.Name }}</v-checkbox>
      </v-flex>
    </v-layout>
    <v-layout row wrap>
      <v-btn
        color="success"
        :disabled="activities && activities.length == 0"
        @click="submit"
        block
        flat
      >Save</v-btn>
    </v-layout>
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
import AppSnackBar from "@/components/utils/AppSnackBar";
import { MODIFY_PROFILE, SAVE_PROFILE } from "@/store/actions/user.js";
import { mapMutations, mapActions } from "vuex";
export default {
  name: "Activities",
  props: ["activities", "allActivities"],
  components: { AppSnackBar },
  data() {
    return {
      snackbar: false,
      snackbarTimeout: 3000,
      snackbarText: "your profile has been successfully saved",
      showSavingProfileDialog: false
    };
  },
  computed: {
    selected_activities: {
      set(val) {
        this.MODIFY_PROFILE({ Activities: val });
      },
      get() {
        return this.$store.state.user.profile.Activities;
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
          that.showSavingProfileDialog = false;
          that.snackbar = true;
        })
        .catch(e => {
          that.snackbarText = "there was and error saving your profile";
          that.snackbar = true;
        });
    }
  }
};
</script>
