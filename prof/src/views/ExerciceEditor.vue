<template>
  <folder-view
    v-if="viewMode == 'folder'"
    :index="exercicesIndex"
    @back="viewMode = 'details'"
  >
  </folder-view>
  <div class="ma-2" v-else>
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
import type { ExercicegroupExt, Index } from "@/controller/api_gen";
import { controller } from "@/controller/controller";
import { onMounted } from "vue";
import { $ref } from "vue/macros";
import ClientPreview from "../components/editor/ClientPreview.vue";
import ExercicegroupList from "../components/editor/exercices/ExercicegroupList.vue";
import ExercicegroupPannel from "../components/editor/exercices/ExercicegroupPannel.vue";
import FolderView from "../components/editor/FolderView.vue";

let viewMode = $ref<"details" | "folder">("folder");

let exercicesIndex = $ref<Index>([]);
async function fetchIndex() {
  const res = await controller.EditorGetExercicesIndex();
  exercicesIndex = res || [];
}

let currentExercicegroup = $ref<ExercicegroupExt | null>(null);

let preview = $ref<InstanceType<typeof ClientPreview> | null>(null);

let allTags = $ref<string[]>([]);
async function fetchTags() {
  const tags = await controller.EditorGetTags();
  allTags = tags || [];
}

onMounted(async () => {
  fetchIndex();
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
