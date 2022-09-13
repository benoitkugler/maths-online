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
        <exercicegroup-pannel
          v-if="currentExercicegroup != null"
          :session_id="sessionID"
          :group="currentExercicegroup"
          :all-tags="allTags"
          @back="backToList"
        ></exercicegroup-pannel>
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
import ExercicegroupList from "../components/editor/exercices/ExercicegroupList.vue";
import ExercicegroupPannel from "../components/editor/exercices/ExercicegroupPannel.vue";

let sessionID = $ref("");

let currentExercicegroup = $ref<ExercicegroupExt | null>(null);

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
}

async function backToList() {
  currentExercicegroup = null;
  controller.EditorPausePreview({ sessionID: sessionID });
}
</script>

<style></style>
