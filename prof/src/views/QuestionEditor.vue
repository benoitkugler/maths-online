<template>
  <div class="ma-2">
    <v-row>
      <v-col>
        <GroupPannel
          v-if="viewKind == 'editor'"
          :session_id="sessionID"
          :group="currentGroup!"
          :variants="currentVariants"
          :all-tags="allKnownTags"
          @back="backToList"
        ></GroupPannel>
        <keep-alive>
          <QuestionList
            v-if="viewKind == 'questions'"
            :tags="allKnownTags"
            @edit="editQuestion"
          ></QuestionList>
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
import GroupPannel from "../components/editor/questions/GroupPannel.vue";
import QuestionList from "../components/editor/questions/QuestionList.vue";

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

// function onDuplicated(question: Question) {
//   currentQuestion = question;
//   // copy to avoid potential side effects
//   currentTags = currentTags.map((v) => v);
//   currentOrigin = personnalOrigin();
// }

function editQuestion(group: QuestiongroupExt, variants: Question[]) {
  currentGroup = group;
  currentVariants = variants;
  viewKind = "editor";
}
</script>

<style></style>
