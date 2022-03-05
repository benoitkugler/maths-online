<template>
  <v-app>
    <v-navigation-drawer app permanent>
      <v-list-item>
        <v-list-item-title class="title"> Pages </v-list-item-title>
      </v-list-item>
      <v-divider></v-divider>

      <v-list-item>
        <v-btn link :to="{ name: 'home' }">Accueil</v-btn>
      </v-list-item>
      <v-list-item>
        <v-btn link :to="{ name: 'trivial' }">Trivial Poursuit</v-btn>
      </v-list-item>
    </v-navigation-drawer>

    <v-app-bar app>
      <v-app-bar-title> Maths online </v-app-bar-title>
      <v-spacer></v-spacer>
      <small>(Version {{ version }})</small>
    </v-app-bar>

    <v-main>
      <router-view />
    </v-main>
  </v-app>
</template>

<script lang="ts">
import { defineComponent } from "vue";

export default defineComponent({
  name: "App",

  data() {
    return {
      version: process.env.VERSION
    };
  }
});
</script>
