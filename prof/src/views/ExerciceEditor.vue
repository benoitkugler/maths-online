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
  type LoopbackShowExercice,
} from "@/controller/api_gen";
import { controller } from "@/controller/controller";
import { emptyTagsDB } from "@/controller/editor";
import { nextTick, ref, onMounted } from "vue";
import ClientPreview from "../components/editor/ClientPreview.vue";
import ExercicegroupList from "../components/editor/exercices/ExercicegroupList.vue";
import ExercicegroupPannel from "../components/editor/exercices/ExercicegroupPannel.vue";
import FolderView from "../components/editor/FolderView.vue";

const viewMode = ref<"details" | "folder">("folder");

const exercicesIndex = ref<Index>([]);
async function fetchIndex() {
  const res = await controller.EditorGetExercicesIndex();
  exercicesIndex.value = res || [];
}

const currentExercicegroup = ref<ExercicegroupExt | null>(null);

const preview = ref<InstanceType<typeof ClientPreview> | null>(null);

const allKnownTags = ref<TagsDB>(emptyTagsDB());
async function fetchTags() {
  const tags = await controller.EditorGetTags();
  if (tags) {
    allKnownTags.value = tags;
  }
}

onMounted(async () => {
  controller.ensureSettings();
  fetchIndex();
  fetchTags();
});

const initialQuery = ref<Query | null>(null);
function showFolder(index: [LevelTag, string]) {
  initialQuery.value = {
    TitleQuery: "",
    LevelTags: [index[0]],
    ChapterTags: [index[1]],
    SubLevelTags: [],
    Origin: OriginKind.All,
    Matiere: controller.settings.FavoriteMatiere,
  };
  viewMode.value = "details";
}

async function editExercice(ex: ExercicegroupExt) {
  currentExercicegroup.value = ex;
}

const list = ref<InstanceType<typeof ExercicegroupList> | null>(null);
function goAndCreateExercicegroup() {
  viewMode.value = "details";
  nextTick(() => {
    if (list.value == null) return;
    list.value.createExercicegroup();
  });
}

async function backToList() {
  currentExercicegroup.value = null;

  preview.value?.pause();
}
</script>

<style></style>
