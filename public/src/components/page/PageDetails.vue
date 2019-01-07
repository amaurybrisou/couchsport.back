<template>
  <v-container v-if="page">
    <div class="app-background" :style="{ 'background-image': 'url(' + backgroundImage + ')' }"></div>
    <v-layout row wrap>
      <v-flex xs6>
        <v-card flat class="fill-height transparent"></v-card>
        <v-card class="page-detail-text">
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
          </v-card-title>
          <v-divider></v-divider>
          <v-card-text class="font-weight-regular body-2">
            Lorem ipsum dolor sit amet consectetur adipisicing elit. Vitae quisquam impedit, cumque suscipit modi ex. Mollitia, molestias repellat, ut nihil modi, ipsa totam eos aperiam in quam optio quos quisquam.
          </v-card-text>
        </v-card>
      </v-flex>
      <v-flex xs6 pl-2>
        <v-card pa-5>
          <l-map
            :zoom="mapConfig.zoom"
            :center="[page.Lat, page.Lng]"
            ref="map"
            style="height:45vh;width:100%;"
          >
            <l-tile-layer :url="mapConfig.url" :attribution="mapConfig.attribution"></l-tile-layer>
          </l-map>
        </v-card>
      </v-flex>
      <v-flex xs12 mt-2>
        <v-layout v-if="page.Images && page.Images.length > 0" row wrap flex>
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
          <v-dialog class="image-dialog" v-model="showImageDialog">
            <v-carousel interval="700000000" height="80vh" hide-delimiters>
              <v-icon
                @click="showImageDialog = false"
                large
                class="right pa-0 ma-1 close-icon"
                icon
              >close</v-icon>
              <v-carousel-item
                v-for="(image) in page.Images"
                :key="image.URL"
                :src="image.URL"
                lazy
              ></v-carousel-item>
            </v-carousel>
          </v-dialog>
        </v-layout>
      </v-flex>
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

export default {
  name: "page-details",
  components: { LMap, LTileLayer, AppSnackBar },
  data() {
    return {
      snackbar: false,
      snackbarTimeout: 3000,
      snackbarText: "an error occured",
      map: null,
      showImageDialog: false,
      mapConfig: {
        zoom: 11,
        center: [46, -1],
        url: "http://{s}.tile.osm.org/{z}/{x}/{y}.png",
        attribution:
          '&copy; <a href="http://openstreetmap.org/copyright">OpenStreetMap</a> contributors',
        showMarkers: true,
        hasSpotMarker: false,
        spotMarker: null
      },
      page: {
        ID: 1,
        CreatedAt: "2019-01-07T14:19:28Z",
        UpdatedAt: "2019-01-07T14:19:34Z",
        DeletedAt: null,
        Name: "Plage des Conches",
        Description: "Surf",
        LongDescription: "lorem",
        Images: [
          {
            ID: 1,
            CreatedAt: "2019-01-07T14:19:34Z",
            UpdatedAt: "2019-01-07T14:19:34Z",
            DeletedAt: null,
            URL: "/static/uploads/1/IMG_6409.JPG",
            Alt: "/static/uploads/1/IMG_6409.JPG",
            File: "",
            PageID: 1
          },
          {
            ID: 2,
            CreatedAt: "2019-01-07T14:19:34Z",
            UpdatedAt: "2019-01-07T14:19:34Z",
            DeletedAt: null,
            URL: "/static/uploads/1/IMG_6410.JPG",
            Alt: "/static/uploads/1/IMG_6410.JPG",
            File: "",
            PageID: 1
          },
          {
            ID: 3,
            CreatedAt: "2019-01-07T14:19:34Z",
            UpdatedAt: "2019-01-07T14:19:34Z",
            DeletedAt: null,
            URL: "/static/uploads/1/IMG_6411.JPG",
            Alt: "/static/uploads/1/IMG_6411.JPG",
            File: "",
            PageID: 1
          }
        ],
        Lat: 46.371536674384544,
        Lng: -1.4821243286132812,
        Followers: null,
        Owner: {
          ID: 0,
          CreatedAt: "0001-01-01T00:00:00Z",
          UpdatedAt: "0001-01-01T00:00:00Z",
          DeletedAt: null,
          Email: "",
          Password: "",
          OwnedPages: null,
          FollowingPages: null,
          Friends: null,
          Profile: {
            ID: 0,
            CreatedAt: "0001-01-01T00:00:00Z",
            UpdatedAt: "0001-01-01T00:00:00Z",
            DeletedAt: null,
            Username: "",
            Country: "",
            Gender: "",
            City: "",
            ZipCode: "",
            StreetName: "",
            Phone: "",
            Firstname: "",
            Lastname: "",
            Avatar: "",
            AvatarFile: "",
            StreetNumber: 0,
            UserID: 0,
            Activities: null,
            Languages: null
          },
          ProfileID: 0,
          Type: "",
          New: false
        },
        OwnerID: 1,
        Public: true,
        Activities: [
          {
            ID: 18,
            Name: "Surf",
            Profiles: null,
            Pages: null
          }
        ],
        New: false
      }
    };
  },
  computed: {
    imagesUrl() {
      return this.page.Images.map(e => e.URL);
    },
    backgroundImage() {
      return this.page && this.page.Images && this.page.Images.length > 0
        ? this.page.Images[0].URL
        : "";
    }
  },
  mounted() {
    var that = this;
    this.$nextTick(function() {
      that.map = this.$refs.map.mapObject;
      if (Number(that.$route.params.page_id) < 1) return that.$router.push("/");
      pageRepo.get(this.$route.params.page_id).then(({ data }) => {
        that.page = data[0];
        L.marker([that.page.Lat, that.page.Lng]).addTo(that.map);
      });
    });
  },
  methods: {
    setTimeout() {
      var that = this;
      setTimeout(function() {
        that.snackbar = false;
      }, that.snackbarTimeout);
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
}
.image-dialog,
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
  opacity: 0.8;
  width: 100%;
  height: 100%;
  top: 0;
  left: 0;
}
</style>
