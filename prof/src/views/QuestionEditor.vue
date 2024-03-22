<template>
  <folder-view
    class="ma-2"
    v-if="viewMode == 'folder'"
    :index="questionsIndex"
    @back="viewMode = 'details'"
    @go-to="showFolder"
  >
    <v-btn class="mx-1" @click="goAndCreateQuestiongroup">
      <v-icon color="green">mdi-plus</v-icon>
      Créer une question
    </v-btn>
  </folder-view>
  <div class="ma-2" v-else>
    <QuestiongroupPannel
      v-if="viewKind == 'editor'"
      :group="currentGroup!"
      :variants="currentVariants"
      :all-tags="allKnownTags"
      @back="backToList"
    ></QuestiongroupPannel>
    <keep-alive>
      <QuestiongroupList
        v-if="viewKind == 'questions'"
        :tags="allKnownTags"
        @edit="editQuestion"
        @back="viewMode = 'folder'"
        :initial-query="initialQuery"
        ref="list"
      ></QuestiongroupList>
    </keep-alive>
  </div>
</template>

<script setup lang="ts">
import {
  OriginKind,
  type Index,
  type LevelTag,
  type Query,
  type Question,
  type QuestiongroupExt,
  type TagsDB,
  type LoopbackShowQuestion,
} from "@/controller/api_gen";
import { controller } from "@/controller/controller";
import ClientPreview from "../components/editor/ClientPreview.vue";
import QuestiongroupList from "../components/editor/questions/QuestiongroupList.vue";
import QuestiongroupPannel from "../components/editor/questions/QuestiongroupPannel.vue";
import FolderView from "../components/editor/FolderView.vue";
import { emptyTagsDB } from "@/controller/editor";
import { ref, onMounted, nextTick } from "vue";

const viewMode = ref<"details" | "folder">("folder");

const questionsIndex = ref<Index>([]);
async function fetchIndex() {
  const res = await controller.EditorGetQuestionsIndex();
  questionsIndex.value = res || [];
}

const allKnownTags = ref<TagsDB>(emptyTagsDB());
const preview = ref<InstanceType<typeof ClientPreview> | null>(null);

const viewKind = ref<"questions" | "editor">("questions");
const currentGroup = ref<QuestiongroupExt | null>(null);
const currentVariants = ref<Question[]>([]);

onMounted(async () => {
  controller.ensureSettings();
  fetchIndex();
  fetchTags();
});

async function fetchTags() {
  const tags = await controller.EditorGetTags();
  if (tags) {
    allKnownTags.value = tags;
  }
}

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

function backToList() {
  currentGroup.value = null;
  currentVariants.value = [];

  fetchTags(); // required since the tags may have changed

  viewKind.value = "questions";
  preview.value?.pause();
}

function editQuestion(group: QuestiongroupExt, variants: Question[]) {
  currentGroup.value = group;
  currentVariants.value = variants;
  viewKind.value = "editor";
}

const list = ref<InstanceType<typeof QuestiongroupList> | null>(null);
function goAndCreateQuestiongroup() {
  viewMode.value = "details";
  nextTick(() => {
    if (list.value == null) return;
    list.value.createQuestiongroup();
  });
}
</script>

<style></style>
