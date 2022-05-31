<template>
  <v-card>
    <v-row>
      <v-col>
        <v-card-title>Suivre la session de Triv'Maths</v-card-title>
      </v-col>
      <v-col style="text-align: right">
        <v-btn icon flat class="mx-2" @click="emit('close')">
          <v-icon icon="mdi-close"></v-icon>
        </v-btn>
      </v-col>
    </v-row>

    <v-card-text>
      <v-alert
        color="info"
        class="mb-2"
        v-if="showSessionID"
        style="text-align: center"
      >
        Code de la session :
        <v-chip>
          <b style="font-size: 22px">
            {{ props.runningSession.SessionID }}
          </b>
        </v-chip>
      </v-alert>

      <v-row justify="center">
        <v-col cols="4" v-for="game in summaries">
          <GameMonitor :summary="game" :show-i-d="showGamesID"></GameMonitor>
        </v-col>
      </v-row>
    </v-card-text>
  </v-card>
</template>

<script setup lang="ts">
import type { LaunchSessionOut } from "@/controller/api_gen";
import { GroupStrategyKind } from "@/controller/api_gen";
import { TrivialMonitorController } from "@/controller/trivial";
import type {
  gameSummary,
  teacherSocketData,
} from "@/controller/trivial_config_socket_gen";
import { computed, onMounted } from "@vue/runtime-core";
import { $ref } from "vue/macros";
import GameMonitor from "./GameMonitor.vue";

interface Props {
  runningSession: LaunchSessionOut;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "close"): void;
}>();

const showSessionID = computed(
  () =>
    props.runningSession.GroupStrategyKind ==
    GroupStrategyKind.RandomGroupStrategy
);

const showGamesID = computed(
  () =>
    props.runningSession.GroupStrategyKind ==
    GroupStrategyKind.FixedSizeGroupStrategy
);

let summaries = $ref<gameSummary[]>([]);

let ct: TrivialMonitorController;
onMounted(() => {
  ct = new TrivialMonitorController(
    props.runningSession.SessionID,
    onServerData
  );
});

function onServerData(data: teacherSocketData) {
  summaries = data.Games || [];
}
</script>

<style scoped></style>
