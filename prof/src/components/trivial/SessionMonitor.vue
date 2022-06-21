<template>
  <v-card>
    <v-row>
      <v-col>
        <v-card-title>Suivre la session de Triv'Maths</v-card-title>
      </v-col>
      <v-col style="text-align: right">
        <v-btn icon flat class="mx-2" @click="emit('closed')">
          <v-icon icon="mdi-close"></v-icon>
        </v-btn>
      </v-col>
    </v-row>

    <v-card-text>
      <v-row justify="center">
        <v-col cols="4" v-for="game in summaries" :key="game.GameID">
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
import { TrivialMonitorController } from "@/controller/trivial";
import type {
  gameSummary,
  teacherSocketData,
} from "@/controller/trivial_config_socket_gen";
import { onMounted } from "@vue/runtime-core";
import { $ref } from "vue/macros";
import GameMonitor from "./GameMonitor.vue";

// type Props = {};

// const props = defineProps<any>();

const emit = defineEmits<{
  (e: "closed"): void;
}>();

let summaries = $ref<gameSummary[]>([]);

let ct: TrivialMonitorController;
onMounted(() => {
  ct = new TrivialMonitorController(onServerData);
});

function onServerData(data: teacherSocketData) {
  summaries = data.Games || [];
}

async function stopTrivGame(params: stopGame) {
  await controller.StopTrivialGame(params);
}
</script>

<style scoped></style>
