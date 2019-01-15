<template>
  <div fluild grid-list-xs class="v-list__container pa-0">
    <v-layout v-if="local_pages" row wrap>
      <v-flex>
        <v-list>
          <template v-for="(p) in local_pages">
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
                    v-model="p.Public"
                    @change="publishPage(p.ID, $event)"
                  ></v-checkbox>

                  <page-edition-dialog
                    @NewPageCreated="onNewPageCreated"
                    :state="'edit'"
                    :page="p"
                    :allActivities="allActivities"
                  >
                    <template slot="open-btn">
                      <v-btn
                        small
                        :color="$vuetify.breakpoint.xsOnly ? '' : 'primary'"
                        :class="{'v-btn--flat v-btn--floating': $vuetify.breakpoint.xsOnly}"
                        :to="`/pages/${p.ID}`"
                      >
                        <v-icon>visibility</v-icon>
                      </v-btn>
                      <v-btn
                        small
                        :color="$vuetify.breakpoint.xsOnly ? '' : 'primary'"
                        :class="{'v-btn--flat v-btn--floating': $vuetify.breakpoint.xsOnly}"
                        @click.prevent
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

                    <span slot="submitText">Save Modifications</span>
                    <span slot="pageTitle">Edit page : {{p.title}}</span>
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
          @NewPageCreated="onNewPageCreated"
          :state="'new'"
          :page="new_page"
          :allActivities="allActivities"
        >
          <template slot="open-btn">
            <v-btn block color="success" flat>New Page</v-btn>
          </template>
          <span slot="pageTitle">New Page</span>
        </page-edition-dialog>
      </v-flex>
    </v-layout>
    <app-snack-bar :state="snackbar" @snackClose="snackbar = false" :text="snackbarText"></app-snack-bar>
  </div>
</template>


<script>
import { LMap, LTileLayer, LMarker } from "vue2-leaflet";
import PageEditionDialog from "@/components/profile/PageEditionDialog";
import pageRepo from "@/repositories/page.js";
import AppSnackBar from "@/components/utils/AppSnackBar";
import { MODIFY_PROFILE } from "@/store/actions/user.js";
import { mapMutations } from "vuex";

export default {
  name: "Pages",
  props: ["pages", "allActivities"],
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
      },
      local_pages: JSON.parse(JSON.stringify(this.pages))
    };
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
    ...mapMutations([MODIFY_PROFILE]),
    deletePage(id) {
      if (id != null) {
        var that = this;
        pageRepo
          .delete({ ID: id })
          .then(function({ data }) {
            that.local_pages = that.local_pages.filter(p => p.ID != id);
            that.snackbarText = "the page has been successfully deleted";
            that.snackbar = true;
          })
          .catch(() => {
            that.snackbarText = "there was an error deleting your page";
            that.snackbar = true;
          });
      }
    },
    onNewPageCreated(page, state) {
      if (state === "error") {
        this.snackbarText = "there was an error saving your page";
        this.snackbar = true;
        return;
      } else if (state === "edit") {
        this.snackbarText = "your page has been successfully edited";
        var idx = this.local_pages.map(p => p.ID).indexOf(page.ID);
        this.local_pages[idx] = page;
        this.MODIFY_PROFILE({ OwnedPages: this.local_pages });
        this.snackbar = true;
        return;
      } else if (state === "new") {
        this.local_pages.push(page);
        this.MODIFY_PROFILE({ OwnedPages: this.local_pages });
        this.new_page = {
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
        };
        this.snackbar = true;
      }
    },
    publishPage(id, state) {
      var that = this;
      if (id != null && (state == false || state == true)) {
        pageRepo
          .publish({ ID: id, Public: state })
          .then(({ data }) => {
            that.snackbarText = state
              ? "your page is now public"
              : "your page is now private";
            that.snackbar = true;
          })
          .catch(() => {
            that.local_pages[id].Public = true;
            that.snackbarText = "there was an error publishing your page";
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

