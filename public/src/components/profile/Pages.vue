<template>
  <v-container class="third-profile-step" fluild grid-list-md>
    <v-layout v-if="local_pages" row wrap>
      <v-flex>
        <v-list>
          <template v-for="(p) in local_pages">
            <v-divider :key="p.ID"></v-divider>
            <v-list-tile class="page-line" :key="`preview-image-${p.ID}`" avatar>
              <v-list-tile-avatar v-if="p.Images && p.Images.length > 0">
                <img :src="p.Images[0].URL" :alt="p.Images[0].Alt">
              </v-list-tile-avatar>

              <v-list-tile-content>
                <v-list-tile-title>{{ p.Description }}</v-list-tile-title>
              </v-list-tile-content>
              <v-list-tile-action>
                <v-layout row>
                  <v-checkbox
                    class="align-center"
                    color="primary"
                    label="Public"
                    v-model="p.Public"
                    @change="publishPage(p.ID, $event)"
                  ></v-checkbox>

                  <v-flex>
                    <page-edition-dialog @NewPageCreated="NewPageCreated" :state="'edit'" :page="p" :allActivities="allActivities">
                      <template slot="open-btn">
                        <v-btn :to="`/pages/${p.ID}`" class="align-center" color="primary">
                          <v-icon>visibility</v-icon>
                        </v-btn>
                        <v-btn color="primary" @click.prevent>
                          <v-icon>edit</v-icon>
                        </v-btn>
                      </template>

                      <span slot="submitText">Save Modifications</span>
                      
                      <span slot="pageTitle">Edit page : {{p.title}}</span>
                    </page-edition-dialog>

                    <v-btn color="primary" @click.prevent="deletePage(p.ID)">
                      <v-icon>delete</v-icon>
                    </v-btn>
                  </v-flex>
                </v-layout>
              </v-list-tile-action>
            </v-list-tile>
          </template>
        </v-list>
      </v-flex>
    </v-layout>
    <v-layout row wrap justify-center align-center>
      <v-flex xs1>
        <page-edition-dialog
          @NewPageCreated="NewPageCreated"
          :state="'new'"
          :page="new_page"
          :allActivities="allActivities"
        >
          <template slot="open-btn">
            <v-btn color="success" flat>New Page</v-btn>
          </template>
          <span slot="pageTitle">New Page</span>
        </page-edition-dialog>
      </v-flex>
    </v-layout>
    <app-snack-bar :state="snackbar" @snackClose="snackbar = false" :text="snackbarText"></app-snack-bar>
  </v-container>
</template>


<script>
import { LMap, LTileLayer, LMarker } from "vue2-leaflet";
import PageEditionDialog from "@/components/profile/PageEditionDialog";
import pageRepo from "@/repositories/page.js";
import { AUTH_ERROR } from "@/store/actions/auth";
import AppSnackBar from "@/components/utils/AppSnackBar";
import { mapGetters } from "vuex";

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
    NewPageCreated(page, state) {
      if (state === "error") {
        this.snackbarText = "there was an error saving your page";
        this.snackbar = true;
        return;
      } else if (state === "edit") {
        this.snackbarText = "your page has been successfully edited";
        var idx = this.local_pages.map(p => p.ID).indexOf(page.ID);

        this.local_pages[idx] = page;
        this.snackbar = true;
        return;
      } else if (state === "new") {
        this.local_pages.push(page);
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

