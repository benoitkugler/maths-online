<template>
  <folder-view
    class="ma-2"
    v-if="viewMode == 'folder'"
    :index="questionsIndex"
    @back="viewMode = 'details'"
    @go-to="showFolder"
  >
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
          @preview="(qu) => preview?.showQuestion(qu)"
        ></QuestiongroupPannel>
        <keep-alive>
          <QuestiongroupList
            v-if="viewKind == 'questions'"
            :tags="allKnownTags"
            @edit="editQuestion"
            @back="viewMode = 'folder'"
            :initial-query="initialQuery"
          ></QuestiongroupList>
        </keep-alive>
      </v-col>
      <v-col cols="auto">
        <keep-alive>
          <ClientPreview ref="preview"></ClientPreview>
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
} from "@/controller/api_gen";
import { controller } from "@/controller/controller";
import { onMounted } from "@vue/runtime-core";
import { $ref } from "vue/macros";
import ClientPreview from "../components/editor/ClientPreview.vue";
import QuestiongroupList from "../components/editor/questions/QuestiongroupList.vue";
import QuestiongroupPannel from "../components/editor/questions/QuestiongroupPannel.vue";
import FolderView from "../components/editor/FolderView.vue";

let viewMode = $ref<"details" | "folder">("folder");

let questionsIndex = $ref<Index>([]);
async function fetchIndex() {
  const res = await controller.EditorGetQuestionsIndex();
  questionsIndex = res || [];
}

let allKnownTags = $ref<TagsDB>({
  Levels: [],
  ChaptersByLevel: {},
  TrivByChapters: {},
});
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
    LevelTags: index[0] ? [index[0]] : [],
    ChapterTags: index[1] ? [index[1]] : [],
    Origin: OriginKind.All,
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
</script>

<style></style>
