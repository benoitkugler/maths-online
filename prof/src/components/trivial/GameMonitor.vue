<template>
  <v-dialog v-model="showConfirmStopGame" max-width="800">
    <v-card title="Terminer la partie">
      <v-card-text> Confirmez-vous l'interruption de la partie ? </v-card-text>
      <v-card-actions>
        <v-btn @click="showConfirmStopGame = false"> Retour </v-btn>
        <v-spacer></v-spacer>
        <v-btn @click="emitStopGame(true)" color="warning">
          Relancer la partie
        </v-btn>
        <v-btn @click="emitStopGame(false)" color="red">
          Terminer la partie
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
  <v-card class="ma-2">
    <v-card-text style="font-size: 16px" class="px-2">
      <v-row
        :justify="props.showID ? 'space-between' : 'center'"
        class="mb-2"
        no-gutters
      >
        <v-col cols="5" align-self="center">
          <div v-if="props.showID">
            Code :
            <v-chip>
              <b style="font-size: 26px">
                {{ props.summary.GameID }}
              </b>
            </v-chip>
          </div>
        </v-col>

        <v-col
          cols="3"
          v-if="props.summary.LatestQuestion.Id != 0"
          align-self="center"
        >
          <v-btn
            density="comfortable"
            rounded
            @click="emit('showQuestion', props.summary.LatestQuestion)"
            :color="colorsPerCategorie[props.summary.LatestQuestion.Categorie]"
            >Question {{ props.summary.LatestQuestion.Id }}</v-btn
          >
        </v-col>

        <v-col cols="auto" style="text-align: right" align-self="center">
          <v-chip color="info">
            {{ props.summary.RoomSize }} joueur{{
              props.summary.RoomSize > 1 ? "s" : ""
            }}
          </v-chip>
          <v-btn
            size="x-small"
            icon
            class="ml-1"
            title="Terminer la partie"
            @click="showConfirmStopGame = true"
          >
            <v-icon icon="mdi-close"></v-icon>
          </v-btn>
        </v-col>
      </v-row>
      <v-row no-gutters class="my-1 px-2">
        <v-col v-if="!props.summary.Players?.length">
          <i>En attente de joueurs...</i>
        </v-col>
        <v-col
          cols="6"
          v-for="(player, index) in props.summary.Players"
          :key="index"
        >
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
import type {
  GameSummary,
  QuestionContent,
  stopGame,
} from "@/controller/api_gen";
import { colorsPerCategorie } from "@/controller/trivial";
import { $ref } from "vue/macros";
import Pie from "./Pie.vue";

interface Props {
  summary: GameSummary;
  showID: boolean;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "stopGame", args: stopGame): void;
  (e: "showQuestion", args: QuestionContent): void;
}>();

let showConfirmStopGame = $ref(false);
function emitStopGame(restart: boolean) {
  emit("stopGame", { ID: props.summary.GameID, Restart: restart });
  showConfirmStopGame = false;
}
</script>

<style scoped></style>
