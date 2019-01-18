<template>
  <div fluild grid-list-xs class="v-list__container pa-0">
    <v-layout v-if="pages" row wrap>
      <v-flex>
        <v-list>
          <template v-for="(p, idx) in pages">
            <v-divider :key="p.ID"></v-divider>
            <div class="v-list__tile page-line px-0 mx-0" :key="`preview-image-${p.ID}`">
              <v-layout row wrap d-flex class="v-list__tile px-0 mx-0">
                <v-flex class="v-list__tile__avatar xs2 sm2 md2 px-0 mx-0">
                  <v-list-tile-avatar v-if="p.Images && p.Images.length > 0">
                    <img :src="p.Images[0].URL" :alt="p.Images[0].Alt">
                  </v-list-tile-avatar>
                </v-flex>
                <v-flex class="v-list__tile__content xs2 sm2 md2 px-0 mx-0">
                  <v-list-tile-title>{{ p.Description| shorten(10) }}</v-list-tile-title>
                </v-flex>

                <v-flex class="xs8 sm8 md8 px-0 d-flex v-list__tile__action text-xs-right">
                  <v-checkbox
                    :class="{'v-btn v-btn--small v-btn--flat v-btn--floating': $vuetify.breakpoint.xsOnly}"
                    :color="$vuetify.breakpoint.xsOnly ? '' : 'primary'"
                    :label="$vuetify.breakpoint.xsOnly ? '' : (p.Public ? 'Public' : 'Private')"
                    v-model="isPublic[idx]"
                    @change="publishPage(p.ID, $event)"
                  ></v-checkbox>

                  <page-edition-dialog
                    @page_saved="onPageSaved"
                    :state="'edit'"
                    :allActivities="allActivities"
                  >
                    <template slot="open-btn">
                      <v-btn
                        small
                        :color="$vuetify.breakpoint.xsOnly ? '' : 'primary'"
                        :class="{'v-btn--flat v-btn--floating': $vuetify.breakpoint.xsOnly}"
                        :to="{ name: 'page-details', params: { page_id: p.ID }}"
                      >
                        <v-icon>visibility</v-icon>
                      </v-btn>
                      <v-btn
                        small
                        :color="$vuetify.breakpoint.xsOnly ? '' : 'primary'"
                        :class="{'v-btn--flat v-btn--floating': $vuetify.breakpoint.xsOnly}"
                        v-on:click="editPage(p.ID)"
                      >
                        <v-icon>edit</v-icon>
                      </v-btn>
                      <v-btn
                        small
                        :color="$vuetify.breakpoint.xsOnly ? '' : 'primary'"
                        :class="{'v-btn--flat v-btn--floating': $vuetify.breakpoint.xsOnly}"
                        v-on:click.stop="deletePage(p.ID)"
                      >
                        <v-icon>delete</v-icon>
                      </v-btn>
                    </template>

                    <span slot="submitText">{{ $t('save') }} {{ $t('modifications') }}</span>
                    <span slot="pageTitle">{{ $t('edit') }} {{ $t('page') }} : {{p.title}}</span>
                  </page-edition-dialog>
                </v-flex>
              </v-layout>
            </div>
          </template>
        </v-list>
      </v-flex>
    </v-layout>
    <v-layout>
      <v-flex d-flex>
        <page-edition-dialog
          @page_saved="onPageSaved"
          :state="'new'"
          :page="new_page"
          :allActivities="allActivities"
        >
          <template slot="open-btn">
            <v-btn block color="success" flat>{{ $t('new') }} {{ $t('page') }}</v-btn>
          </template>
          <span slot="pageTitle">{{ $t('new') }} {{ $t('page') }}</span>
        </page-edition-dialog>
      </v-flex>
    </v-layout>
    <app-snack-bar :state="snackbar" @snackClose="snackbar = false" :text="snackbarText"></app-snack-bar>
  </div>
</template>


<script>
import { LMap, LTileLayer, LMarker } from "vue2-leaflet";
import PageEditionDialog from "@/components/profile/PageEditionDialog";
import AppSnackBar from "@/components/utils/AppSnackBar";

import {
  MODIFY_PAGE,
  GET_PAGES,
  EDIT_PAGE,
  PUBLISH_PAGE,
  DELETE_PAGE,
  PAGE_ADD_PHOTO,
  NEW_PAGE
} from "@/store/actions/pages";
import { mapMutations, mapActions, mapState, mapGetters } from "vuex";

const NAMESPACE = "pages/";

export default {
  name: "Pages",
  components: { LMap, LTileLayer, LMarker, PageEditionDialog, AppSnackBar },
  data() {
    return {
      snackbar: false,
      snackbarTimeout: 3000,
      snackbarText: "your page has been successfully created",
      new_page: {
        ID: null,
        Name: "",
        Description: "",
        Public: true,
        LongDescription: "",
        Images: [],
        Lat: null,
        Lng: null,
        CouchNumber: 0,
        Activities: []
      }
    };
  },
  computed: {
    ...mapGetters({ allActivities: "activities" }),
    pages: {
      get() {
        return this.$store.state.profile.pages.pages;
      }
    },
    isPublic: {
      get() {
        return this.$store.state.profile.pages.pages.map((p, i) => {
          return p.Public ? "true" : "false";
        });
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
  mounted() {
    this[NAMESPACE + GET_PAGES]();
  },
  methods: {
    ...mapActions([
      NAMESPACE + GET_PAGES,
      NAMESPACE + NEW_PAGE,
      NAMESPACE + PUBLISH_PAGE,
      NAMESPACE + DELETE_PAGE
    ]),
    ...mapMutations([NAMESPACE + EDIT_PAGE]),
    onPageSaved(state) {
      if (state) {
        this.snackbarText = this.$t("message.success_saving", ["page"]);
        this.snackbar = true;
      } else {
        this.snackbarText = this.$t("message.error_saving", ["page"]);
        this.snackbar = true;
      }
    },
    editPage(id) {
      this[NAMESPACE + EDIT_PAGE](id);
    },
    deletePage(id) {
      if (id != null) {
        var that = this;
        this[NAMESPACE + DELETE_PAGE]({ ID: id })
          .then(function() {
            that.snackbarText = this.$t("message.success_deleting", ["page"]);
            that.snackbar = true;
          })
          .catch(() => {
            that.snackbarText = this.$t("message.error_deleting", ["page"]);
            that.snackbar = true;
          });
      }
    },
    publishPage(id, state) {
      var that = this;
      if (id != null && (state == false || state == true)) {
        this[NAMESPACE + PUBLISH_PAGE]({ ID: id, Public: state })
          .then(() => {
            that.snackbarText = state
              ? this.$t("message.state", ["page", "public"])
              : this.$t("message.state", ["page", "private"]);
            that.snackbar = true;
          })
          .catch(() => {
            that.pages[id].Public = true;
            that.snackbarText = this.$t("message.error_updating", ["page"]);
            that.snackbar = true;
          });
      }
    }
  }
};
</script>

<style lang="scss">
.page-map {
  height: 350px;
}

.page-line:hover {
  background: rgba($color: #607d8b, $alpha: 0.12);
}
</style>

