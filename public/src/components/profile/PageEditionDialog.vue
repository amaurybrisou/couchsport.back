<template>
  <v-dialog
    v-model="showEditPageDialog"
    @keydown.esc="cancelEdit()"
    fullscreen
    transition="dialog-bottom-transition"
    :overlay="false"
  >
    <template slot="activator">
      <slot name="open-btn"></slot>
    </template>

    <v-card>
      <v-toolbar dark color="primary">
        <v-toolbar-title>
          <slot name="pageTitle">Edit page : {{Name}}</slot>
        </v-toolbar-title>
        <v-spacer></v-spacer>
        <v-toolbar-items></v-toolbar-items>
        <v-btn dark flat @click.native="submit">
          <slot name="submitText">Save</slot>
        </v-btn>
        <v-btn icon @click.native.prevent="cancelEdit()" dark>
          <v-icon>close</v-icon>
        </v-btn>
      </v-toolbar>
      <template>
        <v-container
          fluid
          :class="{ 'sm4 px-5 pb-0': $vuetify.breakpoint.smAndUp, 'xs12 pa-1': $vuetify.breakpoint.xsOnly }"
        >
          <v-layout row wrap>
            <v-flex
              :class="{ 'sm6 pr-1 pb-0': $vuetify.breakpoint.smAndUp, 'xs12 pa-1': $vuetify.breakpoint.xsOnly }"
            >
              <v-form v-model="valid" ref="form" lazy-validation>
                <v-text-field label="Name" v-model="Name" :rules="nameRules" required hide-details></v-text-field>
                <v-text-field label="Description" v-model="Description" required hide-details></v-text-field>
                <v-textarea
                  name="LongDescription"
                  v-model="LongDescription"
                  maxlength="512"
                  placeholder="Describe the spot with more details"
                  row="1"
                  hide-details
                  no-resize
                ></v-textarea>
                <v-autocomplete
                  v-model="Activities"
                  :items="allActivities"
                  label="Activity"
                  item-text="Name"
                  return-object
                  multiple
                ></v-autocomplete>
                <v-slider
                  v-model="CouchNumber"
                  :rules="couchNumberRules"
                  color="primary"
                  label="Number of couch available"
                  hint="Number of couch available"
                  min="0"
                  max="15"
                  thumb-label
                ></v-slider>
              </v-form>
            </v-flex>

            <v-flex
              :class="{ 'sm6 pl-1 pb-0': $vuetify.breakpoint.smAndUp, 'xs12 pa-1': $vuetify.breakpoint.xsOnly }"
            >
              <v-card pa-5>
                <div>
                  <v-subheader color="warning">
                    <v-icon color="warning">help</v-icon>Left click to the determine the spot location, right click to remove it!
                  </v-subheader>
                </div>
                <l-map
                  :zoom="mapConfig.zoom"
                  :center="mapConfig.center"
                  :maxBounds="mapConfig.maxBounds"
                  :nowWrap="mapConfig.nowWrap"
                  ref="map"
                  style="height:40vh;width:100%;"
                >
                  <l-tile-layer :url="mapConfig.url" :attribution="mapConfig.attribution"></l-tile-layer>
                </l-map>
              </v-card>
            </v-flex>
            <v-flex v-if="Images" xs12 mt-1>
              <upload-button
                label="Add photos of the spot!"
                :multiple="false"
                title="Browser"
                :disabled="Images.length > 5"
                @formData="addImage"
              ></upload-button>
              <v-layout v-if="Images.length > 0" row wrap>
                <v-flex
                  :class="{ 'sm2 px-2': $vuetify.breakpoint.smAndUp, 'xs6 px-1 py-2': $vuetify.breakpoint.xsOnly }"
                  v-for="(i, idx) in Images"
                  :key="idx"
                >
                  <v-card class="rounded">
                    <v-img
                      :src="i.URL"
                      :lazy-src="i.URL"
                      :alt="i.Alt"
                      aspect-ratio="1"
                      class="grey lighten-2"
                    >
                      <v-layout slot="placeholder" fill-height align-center justify-center ma-0>
                        <v-progress-circular indeterminate color="grey lighten-5"></v-progress-circular>
                      </v-layout>
                      <input
                        placeholder="Image title"
                        class="image-alt-in"
                        label="Title"
                        :value="i.Alt"
                        @input="setImageAlt(idx, $event)"
                      >
                      <v-btn
                        @click="deleteImage(idx)"
                        type="button"
                        class="right pa-0 ma-1"
                        icon
                        light
                        small
                        color="success"
                      >
                        <v-icon>clear</v-icon>
                        <span slot="loader" class="custom-loader">
                          <v-icon light>cached</v-icon>
                        </span>
                      </v-btn>
                    </v-img>
                  </v-card>
                </v-flex>
              </v-layout>
            </v-flex>
          </v-layout>
          <app-snack-bar :state="snackbar" :text="snackbarText" @snackClose="snackbar = false"></app-snack-bar>
        </v-container>
        <v-dialog lazy v-model="showSavingPageDialog" hide-overlay persistent width="300">
          <v-card color="primary" dark>
            <v-card-text>Please stand by
              <v-progress-linear indeterminate color="white" class="mb-0"></v-progress-linear>
            </v-card-text>
          </v-card>
        </v-dialog>
      </template>
    </v-card>
  </v-dialog>
</template>


<script>
import AppSnackBar from "@/components/utils/AppSnackBar";
import UploadButton from "@/components/utils/UploadButton";
import { LMap, LMarker, LTileLayer } from "vue2-leaflet";

import {
  MODIFY_PAGE,
  PAGE_ADD_IMAGE,
  MODIFY_IMAGE_ALT,
  PAGE_DELETE_IMAGE,
  SAVE_PAGE,
  CANCEL_EDIT_PAGE
} from "@/store/actions/pages";
import { mapMutations, mapActions, mapGetters, mapState } from "vuex";

const NAMESPACE = "pages/";

export default {
  name: "profile-page-edition-dialog",
  props: ["pageID", "state", "allActivities"],
  components: { UploadButton, LMap, LMarker, LTileLayer, AppSnackBar },
  data() {
    return {
      snackbar: false,
      snackbarTimeout: 3000,
      snackbarText: "an error occured",

      isEditing: this.state === "edit",

      valid: true,
      nameRules: [
        v => !!v || "Name is required",
        v => (v && v.length <= 50) || "Name must be less than 10 characters"
      ],

      couchNumberRules: [val => val < 15 || `Really ?!`],

      showEditPageDialog: false,
      showSavingPageDialog: false,

      map: null,
      mapConfig: {
        zoom: 1,
        center: [this.Lat || 46, this.Lng || -1],
        maxBounds: [[-90, -180], [90, 180]],
        noWrap: true,
        url: "http://{s}.tile.osm.org/{z}/{x}/{y}.png",
        attribution:
          '&copy; <a href="http://openstreetmap.org/copyright">OpenStreetMap</a> contributors',
        showMarkers: true,
        hasSpotMarker: false,
        spotMarker: this.Lat && this.Lng ? L.marker([this.Lat, this.Lng]) : null
      }
    };
  },
  computed: {
    Name: {
      get() {
        return this.$store.state.profile.pages.edited_page.Name;
      },
      set(v) {
        this[NAMESPACE + MODIFY_PAGE]({ key: "Name", value: v });
      }
    },
    Description: {
      get() {
        return this.$store.state.profile.pages.edited_page.Description;
      },
      set(v) {
        this[NAMESPACE + MODIFY_PAGE]({ key: "Description", value: v });
      }
    },
    LongDescription: {
      get() {
        return this.$store.state.profile.pages.edited_page.LongDescription;
      },
      set(v) {
        this[NAMESPACE + MODIFY_PAGE]({ key: "LongDescription", value: v });
      }
    },
    Lat: {
      get() {
        return this.$store.state.profile.pages.edited_page.Lat;
      },
      set(v) {
        this[NAMESPACE + MODIFY_PAGE]({ key: "Lat", value: v });
      }
    },
    Lng: {
      get() {
        return this.$store.state.profile.pages.edited_page.Lng;
      },
      set(v) {
        this[NAMESPACE + MODIFY_PAGE]({ key: "Lng", value: v });
      }
    },
    CouchNumber: {
      get() {
        return this.$store.state.profile.pages.edited_page.CouchNumber;
      },
      set(v) {
        this[NAMESPACE + MODIFY_PAGE]({ key: "CouchNumber", value: v });
      }
    },
    Public: {
      get() {
        return this.$store.state.profile.pages.edited_page.Public;
      },
      set(v) {
        this[NAMESPACE + MODIFY_PAGE]({ key: "Public", value: v });
      }
    },
    Activities: {
      get() {
        return this.$store.state.profile.pages.edited_page.Activities;
      },
      set(v) {
        this[NAMESPACE + MODIFY_PAGE]({ key: "Activities", value: v });
      }
    },
    Images: {
      get() {
        return this.$store.state.profile.pages.edited_page.Images;
      },
      set(v) {
        this[NAMESPACE + MODIFY_PAGE]({ key: "Images", value: v });
      }
    }
  },
  mounted() {
    var that = this;
    this.$nextTick(function() {
      that.map = this.$refs.map.mapObject;
      if (that.mapConfig.spotMarker) {
        that.mapConfig.hasSpotMarker = true;
        that.mapConfig.zoom = 5;
        that.mapConfig.spotMarker.addTo(that.map);
      }

      that.map.on("click", that.hasClickOnMap);
      that.map.on("contextmenu", () => {
        if (that.mapConfig.spotMarker) {
          that.mapConfig.spotMarker.removeFrom(that.map);
          that.mapConfig.hasSpotMarker = false;
        }
      });
    });
  },
  watch: {
    showEditPageDialog(v) {
      if (!v) return;
      var that = this;
      setTimeout(function() {
        that.map.invalidateSize();
      }, 200);
    },
    snackbar(v) {
      if (!v) return;
      var that = this;
      setTimeout(function() {
        that.snackbar = false;
      }, that.snackbarTimeout);
    }
  },
  methods: {
    ...mapMutations([
      NAMESPACE + MODIFY_IMAGE_ALT,
      NAMESPACE + MODIFY_PAGE,
      NAMESPACE + CANCEL_EDIT_PAGE,
      NAMESPACE + PAGE_ADD_IMAGE
    ]),
    ...mapActions([NAMESPACE + SAVE_PAGE, NAMESPACE + PAGE_DELETE_IMAGE]),
    submit() {
      if (this.$refs.form.validate()) {
        this.showSavingPageDialog = true;
        var that = this;
        this[NAMESPACE + SAVE_PAGE](this.state)
          .then(() => {
            this.showEditPageDialog = false;
            this.showSavingPageDialog = false;
            this.snackbarText = "your page has been successfully created";
            this.snackbar = true;
          })
          .catch(e => {
            this.showSavingPageDialog = false;
            this.snackbarText = "There was an error creating this page";
            this.snackbar = true;
          });
      }
    },
    clear() {
      this.$refs.form.reset();
    },
    setImageAlt(idx, $event) {
      this[NAMESPACE + MODIFY_IMAGE_ALT]({
        idx: idx,
        value: $event.target.value
      });
    },
    cancelEdit() {
      this[NAMESPACE + CANCEL_EDIT_PAGE]();
      this.showEditPageDialog = false;
    },
    hasClickOnMap(e) {
      if (this.mapConfig.hasSpotMarker || !e.latlng) return;
      this.mapConfig.hasSpotMarker = true;

      this.mapConfig.spotMarker = L.marker(e.latlng);
      this.mapConfig.spotMarker.addTo(this.map);

      this[NAMESPACE + MODIFY_PAGE]({
        pageID: this.ID,
        key: "Lat",
        value: e.latlng.lat
      });
      this[NAMESPACE + MODIFY_PAGE]({
        pageID: this.ID,
        key: "Lng",
        value: e.latlng.lng
      });
    },
    addImage(formData) {
      if (this.Images.length > 5) {
        this.snackbarText = "Maximum number of images allowed";
        return (this.snackbar = true);
      }

      var file = formData.get("file");
      if (file instanceof File) {
        if (file.size > 500000) {
          this.snackbarText = "This image is too big";
          return (this.snackbar = true);
        }

        var exists = this.Images.filter(
          i =>
            (i.URL && i.URL.indexOf(file.name) > -1) ||
            (i.File && i.File.indexOf(file.name) > -1)
        ).length;
        if (exists > 0) {
          this.snackbarText = "This images already exists in this page";
          return (this.snackbar = true);
        }

        var that = this;
        var reader = new FileReader();
        reader.onload = function(e) {
          that[NAMESPACE + PAGE_ADD_IMAGE]({
            URL: e.target.result,
            File: file.name
          });
        };

        reader.readAsDataURL(file);
      }
    },
    deleteImage(idx) {
      this[NAMESPACE + PAGE_DELETE_IMAGE](idx);
      this.snackbarText = "image successfully deleted";
      this.snackbar = true;
    }
  }
};
</script>


<style lang="scss">
.rounded {
  @include rounded(10px);
}

.image-alt-in {
  position: absolute;
  line-height: 27px;
  background-color: rgba($color: #fff, $alpha: 0.8);
  width: 100%;
  bottom: 0;
  padding: 8px;
  &:focus {
    outline: none;
  }
}
</style>
