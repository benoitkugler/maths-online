<template>
  <folder-view
    v-if="viewMode == 'folder'"
    :index="exercicesIndex"
    @back="viewMode = 'details'"
    @go-to="showFolder"
  >
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
          ></exercicegroup-list>
        </keep-alive>
        <exercicegroup-pannel
          v-if="currentExercicegroup != null"
          :group="currentExercicegroup"
          :all-tags="allKnownTags"
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
import {
  LevelTag,
  OriginKind,
  type ExercicegroupExt,
  type Index,
  type Query,
  type TagsDB,
} from "@/controller/api_gen";
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

let allKnownTags = $ref<TagsDB>({
  Levels: [],
  ChaptersByLevel: {},
  TrivByChapters: {},
});
async function fetchTags() {
  const tags = await controller.EditorGetTags();
  if (tags) {
    allKnownTags = tags;
  }
}

onMounted(async () => {
  fetchIndex();
  fetchTags();
});

let initialQuery = $ref<Query | null>(null);
function showFolder(index: [LevelTag, string]) {
  initialQuery = {
    TitleQuery: "",
    LevelTags: index[0] ? [index[0]] : [],
    ChapterTags: index[1] ? [index[1]] : [],
    Origin: OriginKind.All,
  };
  viewMode = "details";
}

async function editExercice(ex: ExercicegroupExt) {
  currentExercicegroup = ex;
}

async function backToList() {
  currentExercicegroup = null;

  preview?.pause();
}
</script>

<style></style>
