<template>
  <folder-view
    v-if="viewMode == 'folder'"
    :index="exercicesIndex"
    @back="viewMode = 'details'"
    @go-to="showFolder"
  >
    <v-btn class="mx-1" @click="goAndCreateExercicegroup">
      <v-icon color="green">mdi-plus</v-icon>
      Créer un exercice
    </v-btn>
  </folder-view>
  <div class="ma-2" v-else>
    <v-row>
      <v-col>
        <keep-alive>
          <exercicegroup-list
            v-if="currentExercicegroup == null"
            @edit="editExercice"
            :tags="allKnownTags"
            @back="viewMode = 'folder'"
            :initial-query="initialQuery"
            ref="list"
          ></exercicegroup-list>
        </keep-alive>
        <exercicegroup-pannel
          v-if="currentExercicegroup != null"
          :group="currentExercicegroup"
          :all-tags="allKnownTags"
          @back="backToList"
          @preview="(ex: LoopbackShowExercice) => preview?.showExercice(ex)"
        ></exercicegroup-pannel>
      </v-col>
      <v-col cols="auto">
        <keep-alive>
          <client-preview
            ref="preview"
            :hide="currentExercicegroup == null"
          ></client-preview>
        </keep-alive>
      </v-col>
    </v-row>
  </div>
</template>

<script setup lang="ts">
import {
  LevelTag,
  OriginKind,
  type ExercicegroupExt,
  type Index,
  type Query,
  type TagsDB,
  type LoopbackShowExercice
} from "@/controller/api_gen";
import { controller } from "@/controller/controller";
import { emptyTagsDB } from "@/controller/editor";
import { onMounted } from "vue";
import { $ref } from "vue/macros";
import ClientPreview from "../components/editor/ClientPreview.vue";
import ExercicegroupList from "../components/editor/exercices/ExercicegroupList.vue";
import ExercicegroupPannel from "../components/editor/exercices/ExercicegroupPannel.vue";
import FolderView from "../components/editor/FolderView.vue";
import { ref } from "vue";
import { nextTick } from "vue";

let viewMode = $ref<"details" | "folder">("folder");

let exercicesIndex = $ref<Index>([]);
async function fetchIndex() {
  const res = await controller.EditorGetExercicesIndex();
  exercicesIndex = res || [];
}

let currentExercicegroup = $ref<ExercicegroupExt | null>(null);

let preview = $ref<InstanceType<typeof ClientPreview> | null>(null);

let allKnownTags = $ref<TagsDB>(emptyTagsDB());
async function fetchTags() {
  const tags = await controller.EditorGetTags();
  if (tags) {
    allKnownTags = tags;
  }
}

onMounted(async () => {
  controller.ensureSettings();
  fetchIndex();
  fetchTags();
});

let initialQuery = $ref<Query | null>(null);
function showFolder(index: [LevelTag, string]) {
  initialQuery = {
    TitleQuery: "",
    LevelTags: [index[0]],
    ChapterTags: [index[1]],
    SubLevelTags: [],
    Origin: OriginKind.All,
    Matiere: controller.settings.FavoriteMatiere
  };
  viewMode = "details";
}

async function editExercice(ex: ExercicegroupExt) {
  currentExercicegroup = ex;
}

const list = ref<InstanceType<typeof ExercicegroupList> | null>(null);
function goAndCreateExercicegroup() {
  viewMode = "details";
  nextTick(() => {
    if (list.value == null) return;
    list.value.createExercicegroup();
  });
}

async function backToList() {
  currentExercicegroup = null;

  preview?.pause();
}
</script>

<style></style>
