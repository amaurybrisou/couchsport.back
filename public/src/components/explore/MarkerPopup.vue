<template>
  <div class="marker-popup">
    <v-card flat elevation-0 flex d-flex column>
      <span v-if="image">
        <img
          :src="image.URL"
          :alt="image.Alt"
          :max-height="image.height"
          :max-width="image.width"
          :height="image.height"
          :width="image.width"
    
          aspect-ratio="1"
          class="grey lighten-2"
    
        />
      </span>
      <v-card-title class="pa-0">
        <div>
          <h3 class="headline mb-0">{{name}}</h3>
        </div>
      </v-card-title>
      <v-divider></v-divider>
      <v-card-text class="grow pa-0 pt-2 pb-2">
        <div class="pa-0 mb-2">{{desc}}</div>

        <div v-if="activities.length > 0">
          <v-chip
            v-for="(a, i) in activities"
            color="primary"
            text-color="white"
            small
            :key="i"
          >{{ a.Name }}</v-chip>
        </div>
      </v-card-text>
      <v-divider></v-divider>
      <v-card-actions>
        <v-btn class="mr-2" color="secondary" small flat :href="url">view details</v-btn>
      </v-card-actions>
    </v-card>
  </div>
</template>

<script>
export default {
  name: "marker-popup",
  props: ["name", "desc", "image", "id", "url", "activities"]
};
</script>
