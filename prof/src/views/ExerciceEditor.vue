<template>
  <div class="ma-2">
    <v-row>
      <v-col cols="8">
        <keep-alive>
          <exercicegroup-list
            v-if="currentExercicegroup == null"
            @edit="editExercice"
            :tags="allTags"
          ></exercicegroup-list>
        </keep-alive>
        <exercice-skeleton
          v-if="currentExercicegroup != null && editMode == 'skeleton'"
          @back="backToList"
          @next="editMode = 'questions'"
          :exercice="currentExercicegroup"
          @update="(v) => (currentExercicegroup = v)"
          :all-tags="allTags"
          :session_id="sessionID"
        ></exercice-skeleton>
        <exercice-editor-pannel
          v-else-if="currentExercicegroup != null && editMode == 'questions'"
          :session_id="sessionID"
          :exercice="currentExercicegroup"
          @update="(v) => (currentExercicegroup = v)"
          @back="editMode = 'skeleton'"
        ></exercice-editor-pannel>
      </v-col>
      <v-col cols="auto">
        <keep-alive>
          <client-preview :session_id="sessionID"></client-preview>
        </keep-alive>
      </v-col>
    </v-row>
  </div>
</template>

<script setup lang="ts">
import type { ExercicegroupExt } from "@/controller/api_gen";
import { controller } from "@/controller/controller";
import { onMounted } from "vue";
import { $ref } from "vue/macros";
import ClientPreview from "../components/editor/ClientPreview.vue";
import ExerciceEditorPannel from "../components/editor/exercices/ExerciceEditorPannel.vue";
import ExercicegroupList from "../components/editor/exercices/ExercicegroupList.vue";
import ExerciceSkeleton from "../components/editor/exercices/ExerciceSkeleton.vue";

let sessionID = $ref("");

let currentExercicegroup = $ref<ExercicegroupExt | null>(null);
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

async function editExercice(ex: ExercicegroupExt) {
  currentExercicegroup = ex;
  editMode = "skeleton";
}

async function backToList() {
  currentExercicegroup = null;
  controller.EditorPausePreview({ sessionID: sessionID });
}
</script>

<style></style>
