<template>
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
        <v-col
          cols="12"
          md="6"
          lg="4"
          v-for="game in summaries"
          :key="game.GameID"
        >
          <GameMonitor
            :summary="game"
            :show-i-d="true"
            @stop-game="stopTrivGame"
          ></GameMonitor>
        </v-col>
      </v-row>
    </v-card-text>
  </v-card>
</template>

<script setup lang="ts">
import type { stopGame } from "@/controller/api_gen";
import { controller } from "@/controller/controller";
import type { gameSummary } from "@/controller/trivial_config_socket_gen";
import { onMounted } from "@vue/runtime-core";
import { $ref } from "vue/macros";
import GameMonitor from "./GameMonitor.vue";

// type Props = {};

// const props = defineProps<any>();

const emit = defineEmits<{
  (e: "closed"): void;
}>();

let summaries = $ref<gameSummary[]>([]);

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
</script>

<style scoped></style>
