<template>
  <v-container>
    <div class="app-background" :style="{ 'background-image': 'url(' + backgroundImage + ')' }"></div>
    <v-layout row wrap>
      <v-flex xs6>
        <v-card flat class="fill-height transparent"></v-card>
        <v-card v-if="page" class="page-detail-text">
          <v-card-title class="title font-weight-bold">
            <v-list-tile>{{ page.Name}}</v-list-tile>
            <v-spacer></v-spacer>
            <div v-if="page.Activities">
              <v-chip
                v-for="(a, i) in page.Activities"
                color="primary"
                text-color="white"
                small
                :key="i"
              >{{ a.Name }}</v-chip>
            </div>
            <v-tooltip bottom v-if="page.CouchNumber > 0">
              <div slot="activator">
                <v-chip color="info" text-color="white" small>{{page.CouchNumber}} Couch available</v-chip>
              </div>
              <span>This spot accepts guests</span>
            </v-tooltip>
            <div v-else>
              <v-chip color="primary" text-color="white" small>Does not accept guests</v-chip>
            </div>
          </v-card-title>
          <v-divider></v-divider>
          <v-card-text class="font-weight-regular body-2">{{ page.LongDescription }}</v-card-text>
          <v-btn
            @click="showContactDialog = true"
            class="contact-btn"
            color="primary"
            block
            absolute
            :disabled="message.FromID == message.ToID"
          >Contact</v-btn>
        </v-card>
      </v-flex>
      <v-flex xs6 pl-2>
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
          <v-flex v-for="(image, idx) in page.Images" :key="idx" align-content-space-between ml-2>
            <v-card class="rounded">
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
            <v-card-title class="title font-weight-regular">Your message to {{ contactName }}</v-card-title>
          </v-toolbar>
          <v-form v-model="messageFormValid">
            <v-card-text>
              <v-text-field
                v-if="!email"
                name="Email"
                label="Your email"
                autocomplete="email"
                v-model="message.Email"
                :rules="emailRules"
              ></v-text-field>
              <v-textarea
                name="Message"
                label="Your Message"
                v-model="message.Text"
                :rules="textRules"
                row="1"
                maxlength="128"
                hide-details
                no-resize
              ></v-textarea>
            </v-card-text>
            <v-card-actions>
              <v-spacer></v-spacer>
              <v-btn color="primary" flat @click.prevent.native="showContactDialog = false">Cancel</v-btn>
              <v-btn
                color="primary"
                flat
                @click.native="sendMessage"
                :disabled="!messageFormValid"
              >Send</v-btn>
            </v-card-actions>
          </v-form>
        </v-card>
      </v-dialog>
    </v-layout>
    <app-snack-bar
      :state="snackbar"
      :text="snackbarText"
      @snackClose="snackbar = false"
      @snackOpen.once="setTimeout"
    ></app-snack-bar>
  </v-container>
</template>


<script>
import AppSnackBar from "@/components/utils/AppSnackBar";
import { LMap, LTileLayer } from "vue2-leaflet";
import pageRepo from "@/repositories/page.js";
import conversationRepo from "@/repositories/conversation.js";
import { mapGetters, mapState } from "vuex";

export default {
  name: "page-details",
  components: { LMap, LTileLayer, AppSnackBar },
  data() {
    return {
      messageFormValid: false,
      emailRules: [
        v => !!v || "E-mail is required",
        v => /.+@.+/.test(v) || "E-mail must be valid"
      ],

      textRules: [
        v => !!v || "Message is required",
        v => (v && v.length >= 20) || "Message must be more than 20 characters"
      ],

      snackbar: false,
      snackbarTimeout: 3000,
      snackbarText: "an error occured",

      showImageDialog: false,
      showContactDialog: false,

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
      FromID: state => state.user.profile.ID
    }),
    message() {
      return { FromID: this.FromID, ToID: null, Email: this.email, Text: "" };
    },
    imagesUrl() {
      return this.page.Images.map(e => e.URL);
    },
    backgroundImage() {
      return this.page && this.page.Images && this.page.Images.length > 0
        ? this.page.Images[0].URL
        : "";
    },
    map: function() {
      return this.$refs.map.mapObject;
    },
    contactName: function() {
      return (
        this.page.Owner.Username ||
        this.page.Owner.Firstname ||
        this.page.Owner.Lastname
      );
    }
  },
  asyncComputed: {
    page: async function() {
      if (Number(this.$route.params.page_id) < 1) return this.$router.push("/");
      return pageRepo
        .get({ id: this.$route.params.page_id, profile: true })
        .then(({ data }) => {
          var page = data[0];

          this.map.setView([page.Lat, page.Lng]);
          L.marker([page.Lat, page.Lng]).addTo(this.map);

          this.message.ToID = page.OwnerID;

          return page;
        });
    }
  },
  watch: {
    snackbar: function(v) {
      !!v && setTimeout((this.snackbar = false), this.snackbarTimeout);
    }
  },
  methods: {
    sendMessage: function(e) {
      conversationRepo
        .sendMessage(this.message)
        .then(() => {
          this.snackbarText = "Your messages has been sent";
          this.snackbar = true;
          this.showContactDialog = false;
        })
        .catch(() => {
          this.snackbarText = "An error occured while sending your message";
          this.snackbar = true;
          this.showContactDialog = false;
        });
    }
  }
};
</script>


<style lang="scss">
.page-detail-text {
  position: relative;
  top: -45vh;
  margin-left: 1vh;
  height: 45vh;
  opacity: 0.8;
  .contact-btn {
    bottom: 0vh;
    margin: 0;
  }
}

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
</style>
