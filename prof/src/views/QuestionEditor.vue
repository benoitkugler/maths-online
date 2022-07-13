<template>
  <div class="ma-2">
    <v-row>
      <v-col>
        <QuestionEditorPannel
          v-if="viewKind == 'editor'"
          :session_id="sessionID"
          @back="backToQuestions"
          @duplicated="onDuplicated"
          :question="currentQuestion!"
          :origin="currentOrigin"
          :tags="currentTags"
          :all-tags="allKnownTags"
        ></QuestionEditorPannel>
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
import type { Origin, Question } from "@/controller/api_gen";
import { controller } from "@/controller/controller";
import { personnalOrigin } from "@/controller/editor";
import { onMounted } from "@vue/runtime-core";
import { $ref } from "vue/macros";
import ClientPreview from "../components/editor/ClientPreview.vue";
import QuestionEditorPannel from "../components/editor/QuestionEditorPannel.vue";
import QuestionList from "../components/editor/QuestionList.vue";

let sessionID = $ref("");
let allKnownTags = $ref<string[]>([]);

let viewKind: "questions" | "editor" = $ref("questions");
let currentQuestion: Question | null = $ref(null);
let currentTags: string[] = $ref([]);
let currentOrigin: Origin = $ref(personnalOrigin());

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

function backToQuestions() {
  fetchTags();
  controller.EditorPausePreview({ sessionID: sessionID });
  viewKind = "questions";
}

function onDuplicated(question: Question) {
  currentQuestion = question;
  // copy to avoid potential side effects
  currentTags = currentTags.map((v) => v);
  currentOrigin = personnalOrigin();
}

function editQuestion(question: Question, tags: string[], origin: Origin) {
  currentQuestion = question;
  currentTags = tags;
  currentOrigin = origin;
  viewKind = "editor";
}
</script>

<style></style>
