<template>
  <v-row class="ma-2">
    <v-col cols="8">
      <keep-alive>
        <exercice-list
          v-if="currentExercice == null"
          @clicked="(ex) => (currentExercice = ex)"
        ></exercice-list>
      </keep-alive>
      <keep-alive>
        <exercice-skeleton
          v-if="currentExercice != null && editMode == 'skeleton'"
          @back="currentExercice = null"
          :exercice="currentExercice"
          :all-tags="allTags"
        ></exercice-skeleton>
      </keep-alive>
    </v-col>
    <v-col cols="4"></v-col>
  </v-row>
</template>

<script setup lang="ts">
import type { ExerciceExt } from "@/controller/api_gen";
import { controller } from "@/controller/controller";
import { onMounted } from "vue";
import { $ref } from "vue/macros";
import ExerciceList from "../components/exercices/ExerciceList.vue";
import ExerciceSkeleton from "../components/exercices/ExerciceSkeleton.vue";

let currentExercice = $ref<ExerciceExt | null>(null);
let editMode = $ref<"skeleton" | "questions">("skeleton");

let allTags = $ref<string[]>([]);
async function fetchTags() {
  const tags = await controller.EditorGetTags();
  allTags = tags || [];
}

onMounted(() => fetchTags());
</script>

<style></style>
