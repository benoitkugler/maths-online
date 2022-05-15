<template>
  <v-card class="ma-2">
    <v-card-text>
      <v-row justify="center" class="mb-2">
        <v-col cols="8">
          <div v-if="props.showID">
            Code de la partie :
            <v-chip>
              <b>
                {{ props.summary.GameID }}
              </b>
            </v-chip>
          </div>
        </v-col>
        <v-col cols="4" style="text-align: right">
          <v-chip color="info">
            {{ props.summary.RoomSize }} joueur{{
              props.summary.RoomSize > 1 ? "s" : ""
            }}
          </v-chip>
        </v-col>
      </v-row>
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
  showID: boolean;
}

const props = defineProps<Props>();
</script>

<style scoped></style>
