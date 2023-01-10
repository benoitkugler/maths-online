<template>
  <v-card>
    <v-row class="pa-2">
      <v-col align-self="center">
        <v-card-title> Ressources par chapitres </v-card-title>
      </v-col>
      <v-col align-self="center" cols="auto">
        <v-btn variant="flat" @click="emit('back')">
          <v-icon>mdi-format-list-bulleted</v-icon>
          Détails
        </v-btn>
      </v-col>
    </v-row>
    <v-card-text>
      <FolderLevelRow
        v-for="(level, index) in props.index"
        :key="index"
        class="my-4"
        :level="level"
        @clicked="(chapter) => emit('goTo', [level.Level, chapter])"
      ></FolderLevelRow>
    </v-card-text>
  </v-card>
</template>

<script setup lang="ts">
import type { Index, LevelTag } from "@/controller/api_gen";
import FolderLevelRow from "./FolderLevelRow.vue";

interface Props {
  index: Index;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "back"): void;
  (e: "goTo", query: [LevelTag, string]): void;
}>();
</script>

<style scoped></style>
