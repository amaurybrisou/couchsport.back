
<template>
  <v-container fluild grid-list-xs>
    <v-layout align-center pb-3>
      <div>{{ $t('p.activities.hint') | capitalize }}</div>
    </v-layout>
    <v-divider></v-divider>
    <v-layout v-if="activities" wrap align-center justify-center>
      <v-flex xs6 md2 v-for="(item, i) in activities" :key="i">
        <v-checkbox
          height="0"
          :label="$t(`allActivities.${item.Name}`) | capitalize"
          :value="item"
          multiple
          v-model="selected_activities"
        >{{ item.Name | capitalize }}</v-checkbox>
      </v-flex>
    </v-layout>
    <v-layout row wrap>
      <v-btn
        color="success"
        :disabled="!selected_activities || selected_activities.length == 0"
        @click="submit"
        block
        flat
      >Save</v-btn>
    </v-layout>
    <app-snack-bar :state="snackbar" @snackClose="snackbar = false" :text="snackbarText"></app-snack-bar>
    <v-dialog v-model="showSavingProfileDialog" hide-overlay persistent width="300">
      <v-card color="primary" dark>
        <v-card-text>{{ $t('message.stand_by') | capitalize }}
          <v-progress-linear indeterminate color="white" class="mb-0"></v-progress-linear>
        </v-card-text>
      </v-card>
    </v-dialog>
  </v-container>
</template>


<script>
import AppSnackBar from "@/components/utils/AppSnackBar";
import { MODIFY_PROFILE, SAVE_PROFILE } from "@/store/actions/profile.js";
import { mapMutations, mapActions, mapGetters, mapState } from "vuex";

export default {
  name: "Activities",
  components: { AppSnackBar },
  data() {
    return {
      snackbar: false,
      snackbarTimeout: 3000,
      snackbarText: this.$t("message.success_saving", [this.$t("profile")]),
      showSavingProfileDialog: false
    };
  },
  computed: {
    ...mapGetters(["activities"]),
    selected_activities: {
      set(val) {
        this.MODIFY_PROFILE({ key: "Activities", value: val });
      },
      get() {
        return this.$store.state.profile.profile.Activities;
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
          that.snackbarText = this.$t("message.success_saving", [this.$t("profile")]);
          that.snackbar = true;
        })
        .catch(e => {
          that.snackbarText = this.$t("message.error_saving", [this.$t("profile")]);
          that.showSavingProfileDialog = false;
          that.snackbar = true;
        });
    }
  }
};
</script>
