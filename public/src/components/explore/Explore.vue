<template>
  <v-container
    fluid
    fill-height
    pa-0
  >
    <v-layout column>
      <v-flex>
        <v-toolbar
          dark
          color="secondary"
        >
          <v-toolbar-title class="hidden-xs-only">
            {{ $t('p.explore.sb_title') | capitalize }}
          </v-toolbar-title>
          <v-autocomplete
            v-model="select"
            :data-intro="$t('p.explore.help.first_step')"
            data-step="1"
            :loading="loading"
            :items="items"
            :search-input.sync="search"
            return-object
            item-text="Name"
            cache-items
            class="mx-3"
            flat
            hide-no-data
            hide-details
            :label="$t('p.explore.sb_placeholder') | capitalize"
            solo-inverted
            :menu-props="{zIndex:'2000'}"
            @input="filterMarkers('spots')"
          />
          <v-btn
            icon
            @click="select = null; filterMarkers('spots');"
          >
            <v-icon>clear</v-icon>
          </v-btn>
        </v-toolbar>
      </v-flex>
      <v-flex xs12>
        <l-map
          ref="map"
          :data-intro="$t('p.explore.help.second_step')"
          data-step="2"
          data-position="top"
          :zoom="mapConfig.zoom"
          :center="mapConfig.center"
          :max-bounds="mapConfig.maxBounds"
          :no-wrap="mapConfig.noWrap"
        >
          <l-tile-layer
            :url="mapConfig.url"
            :attribution="mapConfig.attribution"
            :no-wrap="mapConfig.noWrap"
          />
        </l-map>
      </v-flex>
    </v-layout>
  </v-container>
</template>

<script>
import { LMap, LTileLayer, L } from 'vue2-leaflet'

import pageRepo from 'repositories/pages.js'
import MarkerPopup from 'components/explore/MarkerPopup'
import Vue from 'vue'
import IntroJS from 'mixins/intro'

export default {
  name: 'Explore',
  components: { LMap, LTileLayer },
  mixins: [IntroJS],
  data () {
    return {
      loading: false,
      items: [],
      search: null,
      select: null,
      autocompleteFeed: [],
      mapConfig: {
        zoom: 5,
        center: {
          lat: 47.41322,
          lng: -1.219482
        },
        // maxBounds: [[-120, -210], [120, 210]],
        noWrap: true,
        url: 'http://{s}.tile.osm.org/{z}/{x}/{y}.png',
        attribution:
          '&copy; <a href="http://openstreetmap.org/copyright">OpenStreetMap</a> contributors'
      },
      layers: {
        spots: {
          id: 0,
          Name: 'Spots',
          active: true,
          markers: [],
          popupImage: {
            height: '255px',
            width: '255px'
          },
          popupOptions: {
            maxHeight: 450,
            maxWidth: 450
          }
        }
      }
    }
  },
  watch: {
    search (val) {
      val && val !== this.select && this.querySelections(val)
    }
  },
  created () {
    pageRepo.all().then(resp => {
      this.pages = resp.data
      this.$nextTick(function () {
        this.extractAutoCompleteItems()
        this.map = this.$refs.map.mapObject
        this.initLayer('spots')
        this.filterMarkers('spots')
      })
    })
  },
  mounted () {
    if (this.$route.query.zoom) this.mapConfig.zoom = this.$route.query.zoom
    if (this.$route.query.lat && this.$route.query.lng) {
      this.mapConfig.center.lat = this.$route.query.lat
      this.mapConfig.center.lng = this.$route.query.lng
    } else {
      if (navigator.geolocation) {
        var self = this
        navigator.geolocation.getCurrentPosition(function (position) {
          self.mapConfig.center = {
            lat: position.coords.latitude,
            lng: position.coords.longitude
          }
        })
      }
    }

    var that = this
    this.help.setOption('doneLabel', this.$t('help.next_page'))
    this.help.oncomplete(function () {
      that.$router.push({ name: 'page-details', params: { page_name: 'random' } })
    })
    this.help.start()
  },
  methods: {
    initLayer (layer) {
      var latlng
      for (let index = 0; index < this.pages.length; index++) {
        const p = this.pages[index]
        if (!p.Public) continue

        MarkerPopup.router = this.$router
        MarkerPopup.i18n = this.$i18n
        const MarkerPopupConst = Vue.extend(MarkerPopup)
        const comp = new MarkerPopupConst({
          propsData: {
            id: p.ID,
            url: '/' + this.$i18n.locale + '/pages/' + p.Name,
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
        }).$mount(document.createElement('div'))

        latlng = [p.Lat, p.Lng]

        this.layers[layer].markers.push({
          id: p.ID,
          name: p.Name,
          tags:
            p.Activities instanceof Array
              ? [p.Name].concat(p.Activities.map(e => e.Name))
              : [],
          type: 'marker',
          latlng: latlng,
          show: false,
          obj: L.marker(latlng).bindPopup(
            comp.$el.innerHTML,
            this.layers[layer].popupOptions
          )
        })
      }
    },
    filterMarkers (layer) {
      for (let index = 0; index < this.layers[layer].markers.length; index++) {
        const m = this.layers[layer].markers[index]

        if (this.select == null) {
          m.show = true
          m.obj.addTo(this.map)
          this.map.setZoom(2)
          continue
        }
        if (m.tags.length === 0) {
          m.show = false
          m.obj.removeFrom(this.map)
          continue
        }

        for (let j = 0; j < m.tags.length; j++) {
          const markerTag = m.tags[j]
          if (
            (markerTag || '')
              .toLowerCase()
              .indexOf((this.select.Name || '').toLowerCase()) > -1
          ) {
            m.show = true
            m.obj.addTo(this.map)
            this.map.setZoom(2)
            break
          }
          m.show = false
          m.obj.removeFrom(this.map)
        }
      }
    },
    extractAutoCompleteItems () {
      if (this.pages.length > 0) {
        var that = this
        this.pages.forEach(page => {
          if (page.Public) {
            that.autocompleteFeed = that.autocompleteFeed
              .concat(page.Activities || [])
              .concat(page)
          }
        })
      }
    },
    querySelections (v) {
      this.items = this.autocompleteFeed.filter(e => {
        return (e.Name || '').toLowerCase().indexOf(v.toLowerCase()) > -1
      })
    }
  }
}
</script>
