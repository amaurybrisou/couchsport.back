<template>
  <v-container>
    <div class="app-background" :style="{ 'background-image': 'url(' + backgroundImage + ')' }"></div>
    <v-layout
      :class="{'column mt-1': $vuetify.breakpoint.xsOnly, 'row wrap pr-2': $vuetify.breakpoint.smAndUp}"
    >
      <v-flex xs12 sm6>
        <!-- <v-card flat class="transparent" fill-height></v-card> -->
        <v-card v-if="page" class="flexcard fill-height">
          <v-card-title class="title font-weight-bold pb-0">
            <div class="font-weight-bold">{{ page.Name}}</div>
            <v-spacer></v-spacer>

            <v-tooltip bottom>
              <v-icon slot="activator" color="secondary">language</v-icon>
              <span>{{ talkedLanguages }}</span>
            </v-tooltip>
            <v-spacer></v-spacer>
            <v-tooltip bottom v-if="page.CouchNumber > 0">
              <div slot="activator">
                <v-chip
                  color="info"
                  text-color="white"
                  small
                >{{page.CouchNumber}} {{ $t('p.pd.avail_couch') }}</v-chip>
              </div>
              <span>{{ $t('p.pd.guests') }}</span>
            </v-tooltip>
            <div v-else>
              <v-chip color="primary" text-color="white" small>{{ $t('p.pd.no_guests') }}</v-chip>
            </div>
          </v-card-title>

          <v-card-text class="grow py-0">
            <v-list avatar>
              <v-list-tile v-if="page.Activities">
                <v-chip
                  v-for="(a, i) in page.Activities"
                  color="primary"
                  text-color="white"
                  small
                  :key="i"
                >{{ a.Name | capitalize }}</v-chip>
              </v-list-tile>
              <v-divider></v-divider>
              <div class="py-3 text-break subheading font-weight-regular">{{ page.LongDescription || page.Description }}</div>
            </v-list>
          </v-card-text>

          <v-card-actions class="ma-0 pa-0">
            <v-btn
              depressed
              @click="showContactDialog = true"
              color="primary"
              block
              :disabled="message.FromID == message.ToID || page.CouchNumber == 0"
            >{{ $t('contact') }}</v-btn>
          </v-card-actions>
        </v-card>
      </v-flex>
      <v-flex
        :class="{'xs12 mt-2': $vuetify.breakpoint.xsOnly, 'sm6 pl-2': $vuetify.breakpoint.smAndUp}"
      >
        <v-card pa-5>
          <l-map
            :zoom="mapConfig.zoom"
            :center="this.mapConfig.center"
            :maxBounds="mapConfig.maxBounds"
            :noWrap="mapConfig.noWrap"
            ref="map"
            style="height:45vh;width:100%;"
          >
            <l-tile-layer :url="mapConfig.url" :attribution="mapConfig.attribution"></l-tile-layer>
          </l-map>
        </v-card>
      </v-flex>
      <v-flex xs12 mt-2>
        <v-layout v-if="page && page.Images" row wrap flex>
          <v-flex v-for="(image, idx) in page.Images" :key="idx" align-content-space-between>
            <v-card class="rounded ma-1">
              <v-img
                max-height="250px"
                :src="image.URL"
                :lazy-src="image.URL"
                aspect-ratio="1"
                class="grey lighten-2"
                @click="showImageDialog = true"
              >
                <v-layout slot="placeholder" fill-height align-center justify-center ma-0>
                  <v-progress-circular indeterminate color="grey lighten-5"></v-progress-circular>
                </v-layout>
              </v-img>
            </v-card>
          </v-flex>
          <v-dialog id="image-dialog" v-model="showImageDialog">
            <v-carousel interval="700000000" height="80vh" hide-delimiters>
              <v-icon
                @click="showImageDialog = false"
                large
                class="right pa-0 ma-1 close-icon"
                icon
              >close</v-icon>
              <v-carousel-item v-for="(image, i) in page.Images" :key="i" :src="image.URL" lazy></v-carousel-item>
            </v-carousel>
          </v-dialog>
        </v-layout>
      </v-flex>
    </v-layout>
    <v-layout row justify-center>
      <v-dialog
        v-if="page && (message.FromID != message.ToID)"
        id="contact-dialog"
        v-model="showContactDialog"
        width="500"
      >
        <v-card>
          <v-toolbar color="primary">
            <v-card-title
              class="title font-weight-regular"
            >{{ $t('your') }} {{ $t('_message') }} {{ $t('to') }} {{ contactName }}</v-card-title>
          </v-toolbar>
          <v-form v-model="messageFormValid">
            <v-card-text>
              <v-text-field
                v-if="!email"
                name="Email"
                :label="$t('email')"
                autocomplete="email"
                v-model="message.Email"
                :rules="emailRules"
              ></v-text-field>
              <v-textarea
                name="Message"
                :label="$t('_message')"
                v-model="message.Text"
                :rules="textRules"
                row="1"
                maxlength="128"
                no-resize
                @keyup.ctrl.enter="sendMessage"
              ></v-textarea>
            </v-card-text>
            <v-card-actions>
              <v-spacer></v-spacer>
              <v-btn
                color="primary"
                flat
                @click.prevent.native="showContactDialog = false"
              >{{ $t('cancel') }}</v-btn>
              <v-btn
                color="primary"
                flat
                @click.native="sendMessage"
                :disabled="!messageFormValid"
              >{{ $t('send') }}</v-btn>
            </v-card-actions>
          </v-form>
        </v-card>
      </v-dialog>
    </v-layout>
    <app-snack-bar :state="snackbar" :text="snackbarText" @snackClose="snackbar = false"></app-snack-bar>
  </v-container>
</template>


<script>
import AppSnackBar from "@/components/utils/AppSnackBar";
import { LMap, LTileLayer } from "vue2-leaflet";

import { GET_PAGE } from "@/store/actions/pages";

import { mapGetters, mapState, mapActions } from "vuex";

export default {
  name: "page-details",
  components: { LMap, LTileLayer, AppSnackBar },
  data() {
    return {
      contactName: "",
      backgroundImage: "",
      page: null,
      message: {
        FromID: this.FromID,
        ToID: null,
        Email: this.email,
        Text: ""
      },

      messageFormValid: false,
      emailRules: [
        v => !!v || this.$t("message.required", ["", this.$t("email")]),
        v => /.+@.+/.test(v) || this.$t("message.invalid", [this.$t("email")])
      ],

      textRules: [
        v => !!v || this.$t("message.required", ["", this.$t("_message")]),
        v => (v && v.length >= 20) || this.$t("message.length_above", [20])
      ],

      snackbar: false,
      snackbarTimeout: 3000,
      snackbarText: "an error occured",

      showImageDialog: false,
      showContactDialog: false,

      map: null,
      mapConfig: {
        zoom: 11,
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
    ...mapState({
      email: state => state.auth.email,
      FromID: state => state.profile.profile.ID
    }),
    talkedLanguages() {
      const l = this.page.Owner.Languages;
      let m = this.$t("talk") + " ";
      for (var i in l) {
        m += l[i].Name;
        if (i < l.length - 2) {
          m += ", ";
        }
        if (i == l.length - 2) {
          m += " " + this.$t("and") + " ";
        }
      }
      return m;
    },
    imagesUrl() {
      return this.page.Images.map(e => e.URL);
    }
  },
  mounted() {
    this.map = this.$refs.map.mapObject;
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
  created: async function() {
    if (Number(this.$route.params.page_id) < 1)
      return this.$router.push({ name: "home" });

    this.page = await this["pages/" + GET_PAGE]({
      id: this.$route.params.page_id,
      profile: true
    })
      .then(data => {
        var page = data[0];
        this.message.ToID = page.OwnerID;
        this.contactName =
          page.Owner.Username ||
          page.Owner.Firstname ||
          page.Owner.Lastname ||
          page.Owner.email;

        this.backgroundImage =
          page.Images && page.Images.length > 0 ? page.Images[0].URL : "";

        this.map.setView([page.Lat, page.Lng]);
        L.marker([page.Lat, page.Lng]).addTo(this.map);
        return page;
      })
      .catch(err => {
        console.log(err);
        this.$router.push({ name: "home" });
      });
  },
  methods: {
    ...mapActions(["pages/" + GET_PAGE]),
    sendMessage: function(e) {
      if (!this.message.ToID) return;
      var that = this;
      this.$messenger
        .sendMessage(this.message)
        .then(() => {
          that.snackbarText = this.$t("message.success_sending", [
            this.$t("_message")
          ]);
          that.snackbar = true;
          that.showContactDialog = false;
        })
        .catch(() => {
          that.snackbarText = this.$t("message.error_sending", [
            this.$t("_message")
          ]);
          that.snackbar = true;
          that.showContactDialog = false;
        });
    }
  }
};
</script>


<style lang="scss">
#image-dialog,
#contact-dialog,
.close-icon {
  z-index: 1100;
}

.close-icon {
  position: absolute;
  right: 0;
}

.rounded {
  border-radius: 5px;
}

.round {
  border-radius: 50%;
}

.transparent {
  background-color: rgba($color: #fff, $alpha: 0.3);
}

.app-background {
  position: absolute;
  background-size: cover;
  opacity: 0.4;
  width: 100%;
  height: 100%;
  top: 0;
  left: 0;
}

.flexcard {
  display: flex;
  flex-direction: column;
  .text-break {
    word-break: break-all;
    overflow-y: auto
  }
}
</style>
