<template>
  <div class="ma-2">
    <v-row>
      <v-col>
        <Editor
          v-if="viewKind == 'editor'"
          :session_id="sessionID"
          @back="backToQuestions"
          :question="currentQuestion!"
          :tags="currentTags"
          :all-tags="allKnownTags"
        ></Editor>
        <QuestionList
          v-if="viewKind == 'questions'"
          :tags="allKnownTags"
          @edit="editQuestion"
        ></QuestionList>
      </v-col>
      <v-col md="auto">
        <keep-alive>
          <Preview :session_id="sessionID"></Preview>
        </keep-alive>
      </v-col>
    </v-row>
  </div>
</template>

<script setup lang="ts">
import { controller } from "@/controller/controller";
import type { Question } from "@/controller/exercice_gen";
import { onMounted } from "@vue/runtime-core";
import { $ref } from "vue/macros";
import Editor from "../components/editor/Editor.vue";
import Preview from "../components/editor/Preview.vue";
import QuestionList from "../components/editor/QuestionList.vue";

let sessionID = $ref("");
let allKnownTags = $ref<string[]>([]);

let viewKind: "questions" | "editor" = $ref("questions");
let currentQuestion: Question | null = $ref(null);
let currentTags: string[] = $ref([]);

onMounted(async () => {
  const session = await controller.EditorStartSession(null);
  if (session === undefined) {
    return;
  }
  sessionID = session.ID;

  document.addEventListener("visibilitychange", () => {
    if (document.visibilityState === "hidden") {
      navigator.sendBeacon(
        controller.getURL(`/prof/editor/api/end-preview/${sessionID}`)
      );
    }
  });

  fetchTags();
});

async function fetchTags() {
  const tags = await controller.EditorGetTags();
  allKnownTags = tags || [];
}

function backToQuestions() {
  fetchTags();
  controller.EditorPausePreview({ sessionID: sessionID });
  viewKind = "questions";
}

function editQuestion(question: Question, tags: string[]) {
  currentQuestion = question;
  currentTags = tags;
  viewKind = "editor";
}
</script>

<style></style>
