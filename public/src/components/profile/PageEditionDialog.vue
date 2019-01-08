<template>
  <v-dialog
    v-model="showEditPageDialog"
    @keydown.esc="showEditPageDialog = false"
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
          <slot name="pageTitle">Edit page : {{local_page.Title}}</slot>
        </v-toolbar-title>
        <v-spacer></v-spacer>
        <v-toolbar-items></v-toolbar-items>
        <v-btn dark flat @click.native="submit">
          <slot name="submitText">Save</slot>
        </v-btn>
        <v-btn icon @click.native="showEditPageDialog = false" dark>
          <v-icon>close</v-icon>
        </v-btn>
      </v-toolbar>
      <template>
        <v-container fluid pl-5 pr-5>
          <v-layout row wrap>
            <v-flex xs6>
              <v-form v-model="valid" ref="form" lazy-validation>
                <v-text-field
                  label="Name"
                  v-model="local_page.Name"
                  :rules="nameRules"
                  :counter="10"
                  required
                ></v-text-field>
                <v-text-field
                  label="Description"
                  v-model="local_page.Description"
                  counter="50"
                  required
                ></v-text-field>
                <v-textarea
                  name="LongDescription"
                  v-model="local_page.LongDescription"
                  auto-grow
                  counter="50"
                  placeholder="Describe the spot with more details"
                  row="1"
                  single-line
                ></v-textarea>
                <v-autocomplete
                  v-model="local_page.Activities"
                  :items="activities"
                  label="Activity"
                  item-text="Name"
                  return-object
                  multiple
                ></v-autocomplete>

                <!-- <v-checkbox
                  ma-0
                  color="primary"
                  label="Public"
                  v-model="local_page.Public"
                  required
                ></v-checkbox>-->
              </v-form>
            </v-flex>

            <v-flex pl-2 xs6 class="text-xs-center text-sm-center text-md-center text-lg-center">
              <v-card pa-5>
                <div>
                  <v-subheader color="warning">
                    <v-icon color="warning">help</v-icon>Left click to the determine the spot location, right click to remove it!
                  </v-subheader>
                </div>
                <l-map
                  :zoom="mapConfig.zoom"
                  :center="mapConfig.center"
                  ref="map"
                  style="height:40vh;width:100%;"
                >
                  <l-tile-layer :url="mapConfig.url" :attribution="mapConfig.attribution"></l-tile-layer>
                </l-map>
              </v-card>
            </v-flex>
            <v-flex xs12 mt-1>
              <upload-button
                label="Add photos of the spot!"
                :multiple="false"
                title="Browser"
                :disabled="local_page.Images.length > 5"
                @formData="addImage"
              ></upload-button>
              <v-layout v-if="local_page.Images && local_page.Images.length > 0" row wrap flex>
                <v-flex v-for="(i, idx) in local_page.Images" :key="idx" xs2 pl-4 pr-4>
                  <v-card class="rounded-card">
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
                      <v-btn
                        @click="deleteImage(idx, $event)"
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
      </template>
    </v-card>
  </v-dialog>
</template>


<script>
import AppSnackBar from "@/components/utils/AppSnackBar";
import UploadButton from "@/components/utils/UploadButton";
import { LMap, LTileLayer } from "vue2-leaflet";

import pageRepo from "@/repositories/page.js";
import activityRepo from "@/repositories/activity.js";
import imageRepo from "@/repositories/image.js";
import { AUTH_ERROR } from "@/store/actions/auth";

export default {
  name: "profile-page-edition-dialog",
  props: ["page", "state"],
  components: { UploadButton, LMap, LTileLayer, AppSnackBar },
  data() {
    return {
      snackbar: false,
      snackbarTimeout: 3000,
      snackbarText: "an error occured",

      local_page: this.page,

      isEditing: this.state === "edit",
      valid: true,
      nameRules: [
        v => !!v || "Name is required",
        v => (v && v.length <= 50) || "Name must be less than 10 characters"
      ],

      showEditPageDialog: false,
      map: null,
      mapConfig: {
        zoom: 3,
        center: [46, -1],
        url: "http://{s}.tile.osm.org/{z}/{x}/{y}.png",
        attribution:
          '&copy; <a href="http://openstreetmap.org/copyright">OpenStreetMap</a> contributors',
        showMarkers: true,
        hasSpotMarker: false,
        spotMarker: null
      }
    };
  },
  asyncComputed: {
    activities: async function() {
      return activityRepo.all().then(({ data }) => data);
    }
  },
  mounted() {
    var that = this;
    this.$nextTick(function() {
      that.map = this.$refs.map.mapObject;
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
    submit() {
      if (this.$refs.form.validate()) {
        this.showEditPageDialog = false;
        var that = this;
        console.log(this.local_page);
        pageRepo
          .createOrUpdate(this.local_page)
          .then(({ data }) => {
            that.$emit("NewPageCreated", data, that.state);
          })
          .catch(e => {
            that.$emit("NewPageCreated", e, "error");
          });
      }
    },
    clear() {
      this.$refs.form.reset();
    },
    hasClickOnMap(e) {
      if (this.mapConfig.hasSpotMarker || !e.latlng) return;
      this.mapConfig.hasSpotMarker = true;

      this.mapConfig.spotMarker = L.marker(e.latlng);
      this.mapConfig.spotMarker.addTo(this.map);

      this.local_page.Lat = e.latlng.lat;
      this.local_page.Lng = e.latlng.lng;
    },
    addImage(formData) {
      if (this.local_page.Images.length > 5) {
        this.snackbarText = "Maximum number of images allowed";
        return (this.snackbar = true);
      }

      var file = formData.get("file");
      if (file instanceof File) {
        if (file.size > 500000) {
          this.snackbarText = "This image is too big";
          return (this.snackbar = true);
        }
        var that = this;
        var reader = new FileReader();
        reader.onload = function(e) {
          that.page.Images.push({ URL: e.target.result, File: file.name });
        };

        reader.readAsDataURL(file);
      }
    },
    deleteImage(idx) {
      if (
        idx >= 0 &&
        this.local_page.Images &&
        idx < this.local_page.Images.length
      ) {
        !this.isEditing && this.local_page.Images.splice(idx, 1);
        this.isEditing &&
          imageRepo.delete(this.local_page.Images[idx]).then(({ data }) => {
            this.local_page.Images.splice(idx, 1);
            this.snackbarText = "Image successfully deleted";
            this.snackbar = true;
          });
      }
    }
  }
};
</script>


<style lang="scss">
.rounded-card {
  border-radius: 10px;
}

.custom-loader {
  animation: loader 1s infinite;
  display: flex;
}
@-moz-keyframes loader {
  from {
    transform: rotate(0);
  }
  to {
    transform: rotate(360deg);
  }
}
@-webkit-keyframes loader {
  from {
    transform: rotate(0);
  }
  to {
    transform: rotate(360deg);
  }
}
@-o-keyframes loader {
  from {
    transform: rotate(0);
  }
  to {
    transform: rotate(360deg);
  }
}
@keyframes loader {
  from {
    transform: rotate(0);
  }
  to {
    transform: rotate(360deg);
  }
}
</style>
