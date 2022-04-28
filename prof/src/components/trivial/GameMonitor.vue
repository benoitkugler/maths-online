<template>
  <v-card class="ma-2">
    <v-card-text>
      <v-row no-gutters class="my-1">
        <v-col v-if="!props.summary.Players?.length">
          <i>En attente de joueurs...</i>
        </v-col>
        <v-col cols="6" v-for="player in props.summary.Players">
          <pie
            :label="player.Player"
            :success="player.Successes"
            :highlight="player.Player == props.summary.CurrentPlayer"
          ></pie>
        </v-col>
      </v-row>
    </v-card-text>
  </v-card>
</template>

<script setup lang="ts">
import type { gameSummary } from "@/controller/trivial_config_socket_gen";
import Pie from "./Pie.vue";

interface Props {
  summary: gameSummary;
}

const props = defineProps<Props>();
</script>

<style scoped></style>
