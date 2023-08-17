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
    <v-row>
      <v-col>
        <QuestiongroupPannel
          v-if="viewKind == 'editor'"
          :group="currentGroup!"
          :variants="currentVariants"
          :all-tags="allKnownTags"
          @back="backToList"
          @preview="(qu: LoopbackShowQuestion) => preview?.showQuestion(qu)"
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
      </v-col>
      <v-col cols="auto">
        <keep-alive>
          <ClientPreview
            ref="preview"
            :hide="viewKind == 'questions'"
          ></ClientPreview>
        </keep-alive>
      </v-col>
    </v-row>
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
  type LoopbackShowQuestion
} from "@/controller/api_gen";
import { controller } from "@/controller/controller";
import { onMounted } from "@vue/runtime-core";
import { $ref } from "vue/macros";
import ClientPreview from "../components/editor/ClientPreview.vue";
import QuestiongroupList from "../components/editor/questions/QuestiongroupList.vue";
import QuestiongroupPannel from "../components/editor/questions/QuestiongroupPannel.vue";
import FolderView from "../components/editor/FolderView.vue";
import { emptyTagsDB } from "@/controller/editor";
import { ref } from "vue";
import { nextTick } from "vue";

let viewMode = $ref<"details" | "folder">("folder");

let questionsIndex = $ref<Index>([]);
async function fetchIndex() {
  const res = await controller.EditorGetQuestionsIndex();
  questionsIndex = res || [];
}

let allKnownTags = $ref<TagsDB>(emptyTagsDB());
let preview = $ref<InstanceType<typeof ClientPreview> | null>(null);

let viewKind: "questions" | "editor" = $ref("questions");
let currentGroup: QuestiongroupExt | null = $ref(null);
let currentVariants: Question[] = $ref([]);

onMounted(async () => {
  fetchIndex();
  fetchTags();
});

async function fetchTags() {
  const tags = await controller.EditorGetTags();
  if (tags) {
    allKnownTags = tags;
  }
}

let initialQuery = $ref<Query | null>(null);
function showFolder(index: [LevelTag, string]) {
  initialQuery = {
    TitleQuery: "",
    LevelTags: [index[0]],
    ChapterTags: [index[1]],
    SubLevelTags: [],
    Origin: OriginKind.All
  };
  viewMode = "details";
}

function backToList() {
  currentGroup = null;
  currentVariants = [];

  fetchTags(); // required since the tags may have changed

  viewKind = "questions";
  preview?.pause();
}

function editQuestion(group: QuestiongroupExt, variants: Question[]) {
  currentGroup = group;
  currentVariants = variants;
  viewKind = "editor";
}

const list = ref<InstanceType<typeof QuestiongroupList> | null>(null);
function goAndCreateQuestiongroup() {
  viewMode = "details";
  nextTick(() => {
    if (list.value == null) return;
    list.value.createQuestiongroup();
  });
}
</script>

<style></style>
