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
          @update="(v) => (currentExercice = v)"
          :all-tags="allTags"
        ></exercice-skeleton>
        <exercice-editor-pannel
          v-else-if="currentExercice != null && editMode == 'questions'"
          :session_id="sessionID"
          :exercice="currentExercice"
          @update="(v) => (currentExercice = v)"
          @back="editMode = 'skeleton'"
        ></exercice-editor-pannel>
      </keep-alive>
    </v-col>
    <v-col cols="auto">
      <keep-alive>
        <client-preview :session_id="sessionID"></client-preview>
      </keep-alive>
    </v-col>
  </v-row>
</template>

<script setup lang="ts">
import type { ExerciceExt, ExerciceHeader } from "@/controller/api_gen";
import { controller } from "@/controller/controller";
import { onMounted } from "vue";
import { $ref } from "vue/macros";
import ClientPreview from "../components/editor/ClientPreview.vue";
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

onMounted(async () => {
  if (!controller.editorSessionID.length) {
    await controller.EditorStartSession();
  }
  sessionID = controller.editorSessionID;

  fetchTags();
});

async function showExercice(ex: ExerciceHeader) {
  const res = await controller.ExerciceGetContent({ id: ex.Exercice.Id });
  if (res == undefined) {
    return;
  }
  currentExercice = res;
}
</script>

<style></style>
