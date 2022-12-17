<template>
  <div class="ma-2">
    <v-row>
      <v-col>
        <keep-alive>
          <exercicegroup-list
            v-if="currentExercicegroup == null"
            @edit="editExercice"
            :tags="allTags"
          ></exercicegroup-list>
        </keep-alive>
        <exercicegroup-pannel
          v-if="currentExercicegroup != null"
          :group="currentExercicegroup"
          :all-tags="allTags"
          @back="backToList"
          @preview="(ex) => preview?.showExercice(ex)"
        ></exercicegroup-pannel>
      </v-col>
      <v-col cols="auto">
        <keep-alive>
          <client-preview ref="preview"></client-preview>
        </keep-alive>
      </v-col>
    </v-row>
  </div>
</template>

<script setup lang="ts">
import type { ExercicegroupExt } from "@/controller/api_gen";
import { controller } from "@/controller/controller";
import { LoopbackServerEventKind } from "@/controller/loopback_gen";
import { onMounted } from "vue";
import { $ref } from "vue/macros";
import ClientPreview from "../components/editor/ClientPreview.vue";
import ExercicegroupList from "../components/editor/exercices/ExercicegroupList.vue";
import ExercicegroupPannel from "../components/editor/exercices/ExercicegroupPannel.vue";

let currentExercicegroup = $ref<ExercicegroupExt | null>(null);

let preview = $ref<InstanceType<typeof ClientPreview> | null>(null);

let allTags = $ref<string[]>([]);
async function fetchTags() {
  const tags = await controller.EditorGetTags();
  allTags = tags || [];
}

onMounted(async () => {
  fetchTags();
});

async function editExercice(ex: ExercicegroupExt) {
  currentExercicegroup = ex;
}

async function backToList() {
  currentExercicegroup = null;

  preview?.pause();
}
</script>

<style></style>
