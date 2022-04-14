<template>
  <v-app>
    <v-navigation-drawer app v-model="showSideBar">
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
      <v-list-item>
        <v-btn link :to="{ name: 'editor' }">Editeur de question</v-btn>
      </v-list-item>
    </v-navigation-drawer>

    <v-app-bar app dense color="secondary">
      <v-app-bar-nav-icon
        @click="showSideBar = !showSideBar"
      ></v-app-bar-nav-icon>
      <v-app-bar-title tag="h5">
        Maths online -
        <b>{{ $route.meta.Label }}</b>
      </v-app-bar-title>
      <v-spacer></v-spacer>
      <small>(Version {{ version }})</small>
    </v-app-bar>

    <v-main>
      <router-view v-slot="{ Component }">
        <transition>
          <keep-alive>
            <component :is="Component" />
          </keep-alive>
        </transition>
      </router-view>
    </v-main>
  </v-app>
</template>

<script lang="ts">
import { defineComponent } from "vue";

export default defineComponent({
  name: "App",
  data() {
    return {
      showSideBar: false,
      version: process.env.VERSION,
      tmp: ""
    };
  },
  components: {}
});
</script>
