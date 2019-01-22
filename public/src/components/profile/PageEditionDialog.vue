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
        <v-btn dark flat @click.native="validate">
          <slot name="submitText">{{ $t('save') }}</slot>
        </v-btn>
        <v-btn icon @click.native.prevent="cancelEdit()" dark>
          <v-icon>close</v-icon>
        </v-btn>
      </v-toolbar>
      <template>
        <v-form v-model="rules.valid" @keypress.enter="validate" ref="form" lazy-validation>
          <v-container
            fluid
            :class="{ 'sm4 px-5 pb-0': $vuetify.breakpoint.smAndUp, 'xs12 pa-1': $vuetify.breakpoint.xsOnly }"
          >
            <v-layout row wrap>
              <v-flex
                :class="{ 'sm6 pr-1 pb-0': $vuetify.breakpoint.smAndUp, 'xs12 pa-1': $vuetify.breakpoint.xsOnly }"
              >
                <v-text-field
                  :label="$t('name') | capitalize"
                  v-model="Name"
                  :rules="rules['Name']"
                  autofocus
                  @keypress.enter="validate"
                  required
                ></v-text-field>
                <v-text-field
                  :label="$t('description') | capitalize"
                  v-model="Description"
                  :rules="rules['Description']"
                  @keypress.enter="validate"
                  required
                ></v-text-field>
                <v-textarea
                  name="LongDescription"
                  v-model="LongDescription"
                  :rules="rules['LongDescription']"
                  @keypress.ctrl.enter="validate"
                  maxlength="512"
                  :placeholder="$t('p.ped.long_desc_ph') | capitalize"
                  rows="3"
                  no-resize
                ></v-textarea>
                <v-autocomplete
                  v-model="Activities"
                  :items="allActivities"
                  :label="$t('activities') | capitalize"
                  :rules="rules['Activities']"
                  item-text="Name"
                  return-object
                  multiple
                >
                  <template slot="selection" slot-scope="data">
                    <v-chip
                      :selected="data.selected"
                      close
                      color="secondary"
                      @input="removeActivity(data.item)"
                    >
                      <v-subheader
                        class="body-2"
                      >{{ $t('allActivities.'+`${data.item.Name}`) | capitalize }}</v-subheader>
                    </v-chip>
                  </template>
                </v-autocomplete>
                <v-slider
                  v-model="CouchNumber"
                  :rules="rules['CouchNumber']"
                  color="primary"
                  :label="$t('p.ped.couch_number')"
                  min="0"
                  max="15"
                  thumb-label
                ></v-slider>
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
                  <div v-if="rules.invalidLocation">
                    <v-alert
                      color="error"
                      :value="$t('message.invalid', [$t('the_m') + ' ' + $t('location')])"
                    >{{ $t('message.invalid', [$t('the_m') + ' ' +$t('location')]) }}</v-alert>
                  </div>
                </v-card>
              </v-flex>
              <v-flex v-if="Images" xs12 mt-1>
                <upload-button
                  :label="$t('p.ped.upload_image_hint') | capitalize"
                  :multiple="false"
                  :accept="rules.imageFormatsAllowed"
                  title="Browser"
                  :disabled="Images.length > 5"
                  :errors="imagesErrors"
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
                        <v-text-field
                          :placeholder="$t('p.ped.image_alt_ph')"
                          :value="i.Alt"
                          height="5"
                          solo
                          single-line
                          append-icon="clear"
                          hide-details
                          color="grey"
                          background-color="rgba(255,255,255,0.7)"
                          :rules="rules['Alt']"
                          @change="setImageAlt(idx, $event)"
                          @click:append="deleteImage(idx)"
                          @keypress.enter="validate"
                          @mouseenter:append="alert('ok')"
                        ></v-text-field>
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
                {{ $t('message.stand_by') | capitalize }}
                <v-progress-linear indeterminate color="white" class="mb-0"></v-progress-linear>
              </v-card-text>
            </v-card>
          </v-dialog>
        </v-form>
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
  CANCEL_EDIT_PAGE,
  REMOVE_ACTIVITY
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

      showEditPageDialog: false,
      showSavingPageDialog: false,

      imagesErrors: [],
      map: null,
      mapConfig: {
        zoom: 1,
        center: [46, -1],
        maxBounds: [[-90, -180], [90, 180]],
        noWrap: true,
        url: "http://{s}.tile.osm.org/{z}/{x}/{y}.png",
        attribution:
          '&copy; <a href="http://openstreetmap.org/copyright">OpenStreetMap</a> contributors',
        showMarkers: true,
        hasSpotMarker: false,
        spotMarker: null
      },

      maxActivitiesAllowed: 3,
      errors: [],
      rules: {
        valid: true,
        invalidLocation: false,
        imageFormatsAllowed: "image/jpeg, image:jpg, image/png, image/gif",
        Name: [
          v => !!v || this.$t("message.required", ["", this.$t("name")]),
          v => (v && v.length <= 50) || this.$t("message.length_below", [50])
        ],
        Description: [
          v =>
            !!v || this.$t("message.required", ["e", this.$t("description")]),
          v =>
            /^[0-9a-zA-ZàáâäãåąčćęèéêëėįìíîïłńòóôöõøùúûüųūÿýżźñçčšžÀÁÂÄÃÅĄĆČĖĘÈÉÊËÌÍÎÏĮŁŃÒÓÔÖÕØÙÚÛÜŲŪŸÝŻŹÑßÇŒÆČŠŽ∂ð ,!?.'-]{0,75}$/i.test(
              v
            ) ||
            this.$t("message.invalid", [
              this.$t("the_f") + " " + this.$t("description")
            ])
        ],
        LongDescription: [
          v =>
            /^[0-9a-zA-ZàáâäãåąčćęèéêëėįìíîïłńòóôöõøùúûüųūÿýżźñçčšžÀÁÂÄÃÅĄĆČĖĘÈÉÊËÌÍÎÏĮŁŃÒÓÔÖÕØÙÚÛÜŲŪŸÝŻŹÑßÇŒÆČŠŽ∂ð ,!?.'-]{0,512}$/i.test(
              v
            ) || this.$t("message.invalid", ["description"])
        ],
        CouchNumber: [v => v < 15 || this.$t("p.ped.wow")],
        Activities: [
          v => !!v || this.$t("message.required", ["e", this.$t("activity")]),
          v =>
            v.length > 0 ||
            this.$t("message.required", ["e", this.$t("activity")]),
          v =>
            v.length <= this.maxActivitiesAllowed ||
            this.$t("message.too_much", [
              this.maxActivitiesAllowed,
              this.$t("activities")
            ])
        ],
        Alt: [
          v =>
            /^[a-zA-Z0-9!? ]{0,15}$/.test(v) ||
            this.$t("message.valid_chars_hint", ["a-zA-Z0-9!? "]),
          v => v.length < 15 || this.$t("p.ped.invalid_image_alt")
        ]
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
    },
    showEditPageDialog(v) {
      if (!v) return;
      if (this.state === "edit" && !!this.Lat && !!this.Lng) {
        this.showEditPageDialog = true;
        this.addMarker([this.Lat, this.Lng]);
        this.mapConfig.zoom = 5;
        this.mapConfig.center = [this.Lat, this.Lng];
      }

      var that = this;
      setTimeout(function() {
        that.map.invalidateSize();
      }, 200);
    }
  },
  methods: {
    ...mapMutations([
      NAMESPACE + MODIFY_IMAGE_ALT,
      NAMESPACE + MODIFY_PAGE,
      NAMESPACE + CANCEL_EDIT_PAGE,
      NAMESPACE + PAGE_ADD_IMAGE,
      NAMESPACE + REMOVE_ACTIVITY
    ]),
    ...mapActions([NAMESPACE + SAVE_PAGE, NAMESPACE + PAGE_DELETE_IMAGE]),
    validate() {
      this.rules.invalidLocation = false;
      if (!this.$refs.form.validate()) {
        return;
      }
      if (!this.Lat || this.Lat < -90 || this.Lat > 90)
        return (this.rules.invalidLocation = true);
      if (!this.Lng || this.Lng < -180 || this.Lng > 180)
        return (this.rules.invalidLocation = true);

      if (this.Images.length === 0) {
        return (this.imagesErrors = [
          this.$t("message.required", ["e", this.$t("image")])
        ]);
      }

      this.submit();
    },
    submit() {
      this.showSavingPageDialog = true;
      this[NAMESPACE + SAVE_PAGE](this.state)
        .then(() => {
          this.showSavingPageDialog = false;
          this.showEditPageDialog = false;
          this.$emit("page_saved", true);
          this.delMarker();
        })
        .catch(e => {
          this.showSavingPageDialog = false;
          this.delMarker();
          this.$emit("page_saved", false);
        });
    },
    removeActivity(activity) {
      this[NAMESPACE + REMOVE_ACTIVITY](activity);
    },
    addMarker(latlng) {
      this.mapConfig.spotMarker = L.marker(latlng);
      this.mapConfig.spotMarker.addTo(this.map);
      this.mapConfig.hasSpotMarker = true;
      this[NAMESPACE + MODIFY_PAGE]({ key: "Lat", value: latlng[0] });
      this[NAMESPACE + MODIFY_PAGE]({ key: "Lng", value: latlng[1] });
    },
    delMarker() {
      this.mapConfig.spotMarker.removeFrom(this.map);
      this.mapConfig.hasSpotMarker = false;
      this[NAMESPACE + MODIFY_PAGE]({ key: "Lat", value: null });
      this[NAMESPACE + MODIFY_PAGE]({ key: "Lng", value: null });
    },
    clear() {
      this.$refs.form.reset();
    },
    setImageAlt(idx, value) {
      console.log(arguments);
      this[NAMESPACE + MODIFY_IMAGE_ALT]({
        idx: idx,
        value: value
      });
    },
    cancelEdit() {
      this[NAMESPACE + CANCEL_EDIT_PAGE]();
      this.showEditPageDialog = false;
      this.imagesErrors = [];
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
          this.snackbarText = this.$t("message.too_big", [
            this.$t("image"),
            "500ko"
          ]);
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
          that.imagesErrors = [];
        };

        reader.readAsDataURL(file);
      }
    },
    deleteImage(idx) {
      this[NAMESPACE + PAGE_DELETE_IMAGE](idx)
        .then(() => {
          this.snackbarText = this.$t("message.success_deleting", [
            this.$t("image")
          ]);
          this.snackbar = true;
        })
        .catch(() => {
          this.snackbarText = this.$t("message.error_deleting", [
            this.$t("image")
          ]);
          this.snackbar = true;
        });
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
