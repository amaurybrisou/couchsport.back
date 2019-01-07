<template>
  <v-container id="explore" fluid fill-height pa-0>
    <v-layout column>
      <v-flex>
        <v-toolbar dark color="secondary">
          <v-toolbar-title>Start your trip</v-toolbar-title>
          <v-autocomplete
            :loading="loading"
            :items="items"
            :search-input.sync="search"
            @input="filterMarkers('spots')"
            v-model="select"
            return-object
            item-text="Name"
            cache-items
            class="mx-3"
            flat
            hide-no-data
            hide-details
            label="Look for a place or an activity ?"
            solo-inverted
            :menu-props="{zIndex:'2000'}"
          ></v-autocomplete>
          <v-btn icon @click="select = null; filterMarkers('spots');">
            <v-icon>clear</v-icon>
          </v-btn>
        </v-toolbar>
      </v-flex>
      <v-flex xs12>
        <l-map :zoom="mapConfig.zoom" :center="mapConfig.center" ref="map">
          <l-tile-layer :url="mapConfig.url" :attribution="mapConfig.attribution"></l-tile-layer>
        </l-map>
      </v-flex>
    </v-layout>
  </v-container>
</template>
    
<script>
import { LMap, LTileLayer, LMarker } from "vue2-leaflet";

import pageRepo from "@/repositories/page.js";
import MarkerPopup from "@/components/explore/MarkerPopup";
import Vue from "vue";

export default {
  name: "Explore",
  components: { LMap, LTileLayer, LMarker, MarkerPopup },
  watch: {
    search(val) {
      val && val !== this.select && this.querySelections(val);
    }
  },
  computed: {
    mapConfig: function() {
      return {
        zoom: this.$route.query.zoom || 5,
        center: [
          this.$route.query.lat || 47.41322,
          this.$route.query.lng || -1.219482
        ],
        url: "http://{s}.tile.osm.org/{z}/{x}/{y}.png",
        attribution:
          '&copy; <a href="http://openstreetmap.org/copyright">OpenStreetMap</a> contributors'
      };
    }
  },
  data() {
    return {
      loading: false,
      items: [],
      search: null,
      select: null,
      pages: [
        //test & dev purpose object
        {
          ID: 1,
          CreatedAt: "2019-01-03T00:00:00Z",
          UpdatedAt: "2019-01-03T00:00:00Z",
          DeletedAt: null,
          Name: "Les Conches",
          Description: "Oh La belle plage!",
          LongDescription:
            "Lorem ipsum dolor, sit amet consectetur adipisicing elit. Suscipit dolor dolorem iusto consequatur dignissimos sunt a molestias vero similique sequi praesentium assumenda dolores blanditiis, laborum aperiam laboriosam saepe architecto magni.",
          Images: [
            {
              ID: 1,
              CreatedAt: "2019-01-03T00:00:00Z",
              UpdatedAt: "2019-01-03T00:00:00Z",
              DeletedAt: null,
              URL: "static/img/uploads/les-conches.jpg",
              Alt: "Un surfer aux conches!",
              PageID: 1
            },
            {
              ID: 2,
              CreatedAt: "2019-01-03T00:00:00Z",
              UpdatedAt: "2019-01-03T00:00:00Z",
              DeletedAt: null,
              URL: "static/img/uploads/les-conches1.jpg",
              Alt: "Un surfer aux conches!",
              PageID: 1
            }
          ],
          Lat: 46.381062,
          Lng: -1.487591,
          Followers: null,
          Owner: {
            ID: 0,
            CreatedAt: "0001-01-01T00:00:00Z",
            UpdatedAt: "0001-01-01T00:00:00Z",
            DeletedAt: null,
            Username: "",
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
              Country: "",
              Sex: "",
              City: "",
              ZipCode: "",
              StreetName: "",
              Phone: "",
              Firstname: "",
              Lastname: "",
              Avatar: {
                ID: 0,
                CreatedAt: "0001-01-01T00:00:00Z",
                UpdatedAt: "0001-01-01T00:00:00Z",
                DeletedAt: null,
                URL: "",
                Alt: "",
                PageID: 0
              },
              AvatarID: 0,
              StreetNumber: 0,
              Activities: null,
              Languages: null
            },
            ProfileID: 0,
            Type: "",
            New: false
          },
          OwnerID: 16,
          Public: false,
          Activities: [
            {
              ID: 74,
              CreatedAt: "2019-01-02T20:40:17Z",
              UpdatedAt: "2019-01-02T20:40:17Z",
              DeletedAt: null,
              Name: "Apnée",
              Profiles: null,
              Pages: null
            },
            {
              ID: 84,
              CreatedAt: "2019-01-02T20:40:17Z",
              UpdatedAt: "2019-01-02T20:40:17Z",
              DeletedAt: null,
              Name: "Course d'orientation",
              Profiles: null,
              Pages: null
            }
          ]
        }
      ],
      autocompleteFeed: [],
      layers: {
        spots: {
          id: 0,
          Name: "Spots",
          active: true,
          markers: [],
          popupImage: {
            //height: "150px",
            width: "255px"
          },
          popupOptions: {
            maxHeight: 450,
            maxWidth: 450
          }
        }
      }
    };
  },
  mounted() {
    pageRepo.all().then(resp => {
      this.pages = resp.data;
      this.$nextTick(function() {
        this.extractAutoCompleteItems();
        this.map = this.$refs.map.mapObject;
        this.initLayer("spots");
        this.filterMarkers("spots");
      });
    });
  },
  methods: {
    initLayer(layer) {
      var latlng;
      for (let index = 0; index < this.pages.length; index++) {
        const p = this.pages[index];
        if (!p.Public) continue;

        MarkerPopup.router = this.$router;
        const MarkerPopupConst = Vue.extend(MarkerPopup);
        const comp = new MarkerPopupConst({
          propsData: {
            url: "/users/page/" + p.ID,
            name: p.Name,
            image:
              p.Images.length > 0
                ? {
                    URL: p.Images[0].URL,
                    Alt: p.Images[0].Alt,
                    width: this.layers[layer].popupImage.width,
                    height: this.layers[layer].popupImage.height
                  }
                : {},
            activities: p.Activities == null ? {} : p.Activities,
            desc: p.Description
          }
        }).$mount(document.createElement("div"));

        latlng = [p.Lat, p.Lng];

        this.layers[layer].markers.push({
          id: p.ID,
          name: p.Name,
          tags:
            p.Activities instanceof Array
              ? [p.Name].concat(p.Activities.map(e => e.Name))
              : [],
          type: "marker",
          latlng: latlng,
          show: false,
          obj: L.marker(latlng).bindPopup(
            comp.$el.innerHTML,
            this.layers[layer].popupOptions
          )
        });
      }
    },
    filterMarkers(layer) {
      for (let index = 0; index < this.layers[layer].markers.length; index++) {
        const m = this.layers[layer].markers[index];

        if (this.select == null) {
          m.show = true;
          m.obj.addTo(this.map);
          continue;
        }
        if (m.tags.length === 0) {
          m.show = false;
          m.obj.removeFrom(this.map);
          continue;
        }

        for (let j = 0; j < m.tags.length; j++) {
          const markerTag = m.tags[j];
          if (
            (markerTag || "")
              .toLowerCase()
              .indexOf((this.select.Name || "").toLowerCase()) > -1
          ) {
            m.show = true;
            m.obj.addTo(this.map);
            break;
          }
          m.show = false;
          m.obj.removeFrom(this.map);
        }
      }
    },
    extractAutoCompleteItems() {
      if (this.pages.length > 0) {
        var that = this;
        this.pages.forEach(page => {
          if (page.Public) {
            that.autocompleteFeed = that.autocompleteFeed
              .concat(page.Activities || [])
              .concat(page);
          }
        });
      }
    },
    querySelections(v) {
      this.items = this.autocompleteFeed.filter(e => {
        return (e.Name || "").toLowerCase().indexOf(v.toLowerCase()) > -1;
      });
    }
  }
};
</script>

<style lang="scss">
</style>
