<template>
  <v-dialog
    v-model="showEditPageDialog"
    @keydown.esc="cancelEdit()"
    fullscreen
    transition="dialog-bottom-transition"
    :overlay="false"
    style="z-index:600;"
  >
    <template slot="activator">
      <slot name="open-btn"></slot>
    </template>

    <v-card>
      <v-toolbar dark color="primary">
        <v-toolbar-title>
          <slot name="pageTitle">{{ $t('edit_page') }} : {{Name}}</slot>
        </v-toolbar-title>
        <v-spacer></v-spacer>
        <v-toolbar-items></v-toolbar-items>
        <v-btn dark flat @click.native="submit">
          <slot name="submitText">{{ $t('save') }}</slot>
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
                <v-text-field
                  :label="$t('p.ped.name')"
                  v-model="Name"
                  :rules="nameRules"
                  required
                  hide-details
                ></v-text-field>
                <v-text-field
                  :label="$t('p.ped.description')"
                  v-model="Description"
                  required
                  hide-details
                ></v-text-field>
                <v-textarea
                  name="LongDescription"
                  v-model="LongDescription"
                  maxlength="512"
                  :placeholder="$t('p.ped.long_desc_ph')"
                  row="1"
                  hide-details
                  no-resize
                ></v-textarea>
                <v-autocomplete
                  v-model="Activities"
                  :items="allActivities"
                  :label="$t('activities') | capitalize"
                  item-text="Name"
                  return-object
                  multiple
                ></v-autocomplete>
                <v-slider
                  v-model="CouchNumber"
                  :rules="couchNumberRules"
                  color="primary"
                  :label="$t('p.ped.couch_number')"
                  :hint="$t('p.ped.couch_number')"
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
                    <v-icon color="warning">help</v-icon>
                    {{ $t('p.ped.map_help') }}
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
                :label="$t('p.ped.upload_image_hint')"
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
            <v-card-text>
              {{ $t('message.stand_by') }}
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
  props: ["state"],
  components: { UploadButton, LMap, LMarker, LTileLayer, AppSnackBar },
  data() {
    return {
      snackbar: false,
      snackbarTimeout: 3000,
      snackbarText: "an error occured",

      isEditing: this.state === "edit",

      valid: true,
      nameRules: [
        v => !!v || this.$t("message.auth.required", ["name"]),
        v => (v && v.length <= 50) || this.$t("message.length_below", [50])
      ],

      couchNumberRules: [val => val < 15 || `Really ?!`],

      showEditPageDialog: false,
      showSavingPageDialog: false,

      map: null,
      mapConfig: {
        zoom: 12,
        center: [46, -1],
        maxBounds: [[-90, -180], [90, 180]],
        noWrap: true,
        url: "http://{s}.tile.osm.org/{z}/{x}/{y}.png",
        attribution:
          '&copy; <a href="http://openstreetmap.org/copyright">OpenStreetMap</a> contributors',
        showMarkers: true,
        hasSpotMarker: false,
        spotMarker: null
      }
    };
  },
  computed: {
    ...mapGetters({ allActivities: "activities" }),
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
      that.map.on("click", that.hasClickOnMap);
      that.map.on("contextmenu", that.delMarker);
    });
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
    ...mapMutations([
      NAMESPACE + MODIFY_IMAGE_ALT,
      NAMESPACE + MODIFY_PAGE,
      NAMESPACE + CANCEL_EDIT_PAGE,
      NAMESPACE + PAGE_ADD_IMAGE
    ]),
    ...mapActions([NAMESPACE + SAVE_PAGE, NAMESPACE + PAGE_DELETE_IMAGE]),
    initEditPageDialog(v) {
      if (!v) return;
      this.showEditPageDialog = true;
      this.addMarker([this.Lat, this.Lng]);
      this.mapConfig.zoom = 5;
      this.mapConfig.center = [this.Lat, this.Lng];

      var that = this;
      setTimeout(function() {
        that.map.invalidateSize();
      }, 200);
    },
    submit() {
      if (this.$refs.form.validate()) {
        this.showSavingPageDialog = true;
        var that = this;
        this[NAMESPACE + SAVE_PAGE](this.state)
          .then(() => {
            this.showSavingPageDialog = false;
            this.showEditPageDialog = false;
            that.$emit("page_saved", true);
          })
          .catch(e => {
            this.showSavingPageDialog = false;
            that.$emit("page_saved_error", false);
          });
      }
    },
    addMarker(latlng) {
      this.mapConfig.spotMarker = L.marker(latlng);
      this.mapConfig.spotMarker.addTo(this.map);
      this.mapConfig.hasSpotMarker = true;
    },
    delMarker() {
      this.mapConfig.spotMarker.removeFrom(this.map);
      this.mapConfig.hasSpotMarker = false;
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
      this.addMarker(e.latlng);

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
        this.snackbarText = this.$t("p.pde.max_images");
        return (this.snackbar = true);
      }

      var file = formData.get("file");
      if (file instanceof File) {
        if (file.size > 500000) {
          this.snackbarText = this.$t("message.too_big", ["image"]);
          return (this.snackbar = true);
        }

        var exists = this.Images.filter(
          i =>
            (i.URL && i.URL.indexOf(file.name) > -1) ||
            (i.File && i.File.indexOf(file.name) > -1)
        ).length;
        if (exists > 0) {
          this.snackbarText = this.$t("message.exist", ["image"]);
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
      this.snackbarText = this.$t("message.success_deleting", ["image"]);
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
