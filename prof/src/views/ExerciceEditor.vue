<template>
  <v-row class="ma-2">
    <v-col cols="8">
      <keep-alive>
        <exercice-list
          v-if="currentExercice == null"
          @clicked="showExercice"
        ></exercice-list>
      </keep-alive>
      <keep-alive>
        <exercice-skeleton
          v-if="currentExercice != null && editMode == 'skeleton'"
          @back="currentExercice = null"
          @next="editMode = 'questions'"
          :exercice="currentExercice"
          :all-tags="allTags"
        ></exercice-skeleton>
        <exercice-editor-pannel
          v-else-if="currentExercice != null && editMode == 'questions'"
          :session_id="sessionID"
          :exercice="currentExercice"
          @back="editMode = 'skeleton'"
        ></exercice-editor-pannel>
      </keep-alive>
    </v-col>
    <v-col cols="4"></v-col>
  </v-row>
</template>

<script setup lang="ts">
import type { ExerciceExt, ExerciceHeader } from "@/controller/api_gen";
import { controller } from "@/controller/controller";
import { onMounted } from "vue";
import { $ref } from "vue/macros";
import ExerciceEditorPannel from "../components/exercices/ExerciceEditorPannel.vue";
import ExerciceList from "../components/exercices/ExerciceList.vue";
import ExerciceSkeleton from "../components/exercices/ExerciceSkeleton.vue";

let sessionID = $ref("");

let currentExercice = $ref<ExerciceExt | null>(null);
let editMode = $ref<"skeleton" | "questions">("skeleton");

let allTags = $ref<string[]>([]);
async function fetchTags() {
  const tags = await controller.EditorGetTags();
  allTags = tags || [];
}

onMounted(() => fetchTags());

async function showExercice(ex: ExerciceHeader) {
  const res = await controller.ExerciceGetContent({ id: ex.Exercice.Id });
  if (res == undefined) {
    return;
  }
  currentExercice = res;
}
</script>

<style></style>
