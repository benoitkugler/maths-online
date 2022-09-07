<template>
  <v-app>
    <v-app-bar app :title="meta.title" color="primary"></v-app-bar>

    <v-main>
      <v-container>
        <router-view></router-view>
      </v-container>
    </v-main>
  </v-app>
</template>

<script setup lang="ts">
import { computed } from "@vue/reactivity";
import { useRoute } from "vue-router";

const route = useRoute();
const meta = computed(() => route.meta);
</script>
