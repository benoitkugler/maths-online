<template>
  <div class="ma-2">
    <v-row>
      <v-col>
        <QuestiongroupPannel
          v-if="viewKind == 'editor'"
          :session_id="sessionID"
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
          ></QuestiongroupList>
        </keep-alive>
      </v-col>
      <v-col cols="auto">
        <keep-alive>
          <ClientPreview :session_id="sessionID"></ClientPreview>
        </keep-alive>
      </v-col>
    </v-row>
  </div>
</template>

<script setup lang="ts">
import type { Question, QuestiongroupExt } from "@/controller/api_gen";
import { controller } from "@/controller/controller";
import { onMounted } from "@vue/runtime-core";
import { $ref } from "vue/macros";
import ClientPreview from "../components/editor/ClientPreview.vue";
import QuestiongroupList from "../components/editor/questions/QuestiongroupList.vue";
import QuestiongroupPannel from "../components/editor/questions/QuestiongroupPannel.vue";

let sessionID = $ref("");
let allKnownTags = $ref<string[]>([]);

let viewKind: "questions" | "editor" = $ref("questions");
let currentGroup: QuestiongroupExt | null = $ref(null);
let currentVariants: Question[] = $ref([]);

onMounted(async () => {
  if (!controller.editorSessionID.length) {
    await controller.EditorStartSession();
  }
  sessionID = controller.editorSessionID;

  fetchTags();
});

async function fetchTags() {
  const tags = await controller.EditorGetTags();
  allKnownTags = tags || [];
}

function backToList() {
  currentGroup = null;
  currentVariants = [];

  fetchTags(); // required since the tags may have changed

  controller.EditorPausePreview({ sessionID: sessionID });
  viewKind = "questions";
}

function editQuestion(group: QuestiongroupExt, variants: Question[]) {
  currentGroup = group;
  currentVariants = variants;
  viewKind = "editor";
}
</script>

<style></style>
