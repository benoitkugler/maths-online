<template>
  <v-card>
    <v-row>
      <v-col>
        <v-card-title>Suivre la session de TrivialPoursuit</v-card-title>
      </v-col>
      <v-col style="text-align: right">
        <v-btn icon flat class="mx-2" @click="emit('close')">
          <v-icon icon="mdi-close"></v-icon>
        </v-btn>
      </v-col>
    </v-row>

    <v-card-text>
      <v-alert color="info" class="mb-2">
        Code de la session :
        <v-chip>
          {{ sessionID }}
        </v-chip>
      </v-alert>

      <v-row>
        <v-col cols="4" v-for="game in summaries">
          <GameMonitor :summary="game"></GameMonitor>
        </v-col>
      </v-row>
    </v-card-text>
  </v-card>
</template>

<script setup lang="ts">
import { TrivialMonitorController } from "@/controller/trivial";
import type {
  gameSummary,
  teacherSocketData
} from "@/controller/trivial_config_socket_gen";
import { onMounted } from "@vue/runtime-core";
import { $ref } from "vue/macros";
import GameMonitor from "./GameMonitor.vue";

interface Props {
  sessionID: string;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "close"): void;
}>();

let summaries = $ref<gameSummary[]>([]);

let ct: TrivialMonitorController;
onMounted(() => {
  ct = new TrivialMonitorController(props.sessionID, onServerData);
});

function onServerData(data: teacherSocketData) {
  summaries = data.Games || [];
}
</script>

<style scoped></style>
