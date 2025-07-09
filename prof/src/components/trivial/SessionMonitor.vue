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
        <v-card-title>Suivre les parties Isy'Triv</v-card-title>
      </v-col>
      <v-col style="text-align: right">
        <v-btn icon flat class="mx-2" @click="onClose">
          <v-icon icon="mdi-close"></v-icon>
        </v-btn>
      </v-col>
    </v-row>

    <v-card-text>
      <v-row justify="center">
        <v-col cols="12" lg="6" v-for="game in summaries" :key="game.GameID">
          <GameMonitor
            :summary="game"
            :show-i-d="true"
            @start-game="startTrivGame"
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
  RoomID,
  stopGame,
} from "@/controller/api_gen";
import { controller } from "@/controller/controller";
import { ref, onMounted } from "vue";
import GameMonitor from "./GameMonitor.vue";
import QuestionMonitor from "./QuestionMonitor.vue";

// type Props = {};

// const props = defineProps<any>();

const emit = defineEmits<{
  (e: "closed"): void;
}>();

const summaries = ref<GameSummary[]>([]);

const refreshDelay = 5000; // milliseconds

let intervalHandle: ReturnType<typeof setInterval>;
onMounted(() => {
  intervalHandle = setInterval(fetchMonitorData, refreshDelay);
  fetchMonitorData();
});

async function fetchMonitorData() {
  const res = await controller.TrivialTeacherMonitor();
  if (res == undefined) return;
  summaries.value = res.Games || [];
}

async function startTrivGame(id: RoomID) {
  const ok = await controller.StartTrivialGame({ "game-id": id });
  if (!ok) return;
  controller.showMessage("Partie lancée avec succés.");

  await fetchMonitorData();
}

async function stopTrivGame(params: stopGame) {
  const res = await controller.StopTrivialGame(params);
  if (res === undefined) return;
  controller.showMessage("Partie interrompue avec succès");

  await fetchMonitorData();
  // automatically close an empty monitor dialog
  if (!summaries.value.length) {
    onClose();
  }
}

function onClose() {
  clearInterval(intervalHandle);
  emit("closed");
}

const questionToMonitor = ref<QuestionContent | null>(null);

function showQuestion(question: QuestionContent) {
  questionToMonitor.value = question;
}
</script>

<style scoped></style>
