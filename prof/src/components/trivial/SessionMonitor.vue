<template>
  <v-dialog
    :model-value="questionToMonitor != null"
    @update:model-value="questionToMonitor = null"
    fullscreen
  >
    <QuestionMonitor
      v-if="questionToMonitor != null"
      @close="questionToMonitor = null"
      :question="questionToMonitor"
    ></QuestionMonitor>
  </v-dialog>

  <v-card>
    <v-row>
      <v-col>
        <v-card-title>Suivre la session de Triv'Maths</v-card-title>
      </v-col>
      <v-col style="text-align: right">
        <v-btn icon flat class="mx-2" @click="onClose">
          <v-icon icon="mdi-close"></v-icon>
        </v-btn>
      </v-col>
    </v-row>

    <v-card-text>
      <v-row justify="center">
        <v-col md="12" lg="6" v-for="game in summaries" :key="game.GameID">
          <GameMonitor
            :summary="game"
            :show-i-d="true"
            @stop-game="stopTrivGame"
            @show-question="showQuestion"
          ></GameMonitor>
        </v-col>
      </v-row>
    </v-card-text>
  </v-card>
</template>

<script setup lang="ts">
import type {
  GameSummary,
  QuestionContent,
  stopGame,
} from "@/controller/api_gen";
import { controller } from "@/controller/controller";
import { onMounted } from "@vue/runtime-core";
import { $ref } from "vue/macros";
import GameMonitor from "./GameMonitor.vue";
import QuestionMonitor from "./QuestionMonitor.vue";

// type Props = {};

// const props = defineProps<any>();

const emit = defineEmits<{
  (e: "closed"): void;
}>();

let summaries = $ref<GameSummary[]>([]);

let intervalHandle: number;
onMounted(() => {
  intervalHandle = setInterval(fetchMonitorData, 5000);
  fetchMonitorData();
});

async function fetchMonitorData() {
  const res = await controller.TrivialTeacherMonitor();
  if (res == undefined) return;
  summaries = res.Games || [];
}

async function stopTrivGame(params: stopGame) {
  await controller.StopTrivialGame(params);
  await fetchMonitorData();
  // automatically close an empty monitor dialog
  if (!summaries.length) {
    onClose();
  }
}

function onClose() {
  clearInterval(intervalHandle);
  emit("closed");
}

let questionToMonitor = $ref<QuestionContent | null>(null);

function showQuestion(question: QuestionContent) {
  questionToMonitor = question;
}
</script>

<style scoped></style>
