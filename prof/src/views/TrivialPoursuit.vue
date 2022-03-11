<template>
  <v-card class="ma-3">
    <v-card-title>Trivial Poursuit</v-card-title>
    <v-card-subtitle>Configurer une partie de Trivial Poursuit</v-card-subtitle>

    <v-form class="mx-4 my-2">
      <v-text-field
        label="Nombre de joueurs"
        type="number"
        min="1"
        v-model.number="launchOptions.NbPlayers"
      ></v-text-field>
      <v-text-field
        label="Durée limite pour une question"
        type="number"
        min="1"
        suffix="sec"
        v-model.number="launchOptions.TimeoutSeconds"
      ></v-text-field>
      <v-row>
        <v-col cols="4" sm="6" md="9"></v-col>
        <v-col cols="8" sm="6" md="3" class="text-right">
          <v-btn block @click="launchGame" :disabled="!isValid" color="success">
            Lancer la partie
          </v-btn>
        </v-col>
      </v-row>

      <v-alert class="px-4 my-2" :model-value="gameCode != ''" color="info">
        <v-row no-gutters>
          <v-col align-self="center">
            Code de
            <a href="/test-eleve" target="_blank">la partie à rejoindre</a> :
          </v-col>
          <v-col>
            <v-chip size="big" class="pa-3"
              ><b>{{ gameCode }}</b></v-chip
            >
          </v-col>
        </v-row>
      </v-alert>
    </v-form>
  </v-card>
</template>

<script setup lang="ts">
import type { LaunchGameIn } from "@/controller/api_gen";
import { controller } from "@/controller/controller";
import { computed, reactive } from "@vue/runtime-core";
import { $ref } from "vue/macros";

const pinLength = 3;

let launchOptions: LaunchGameIn = reactive({
  NbPlayers: 1,
  TimeoutSeconds: 60
});

const isValid = computed(() => !isSpinning && launchOptions.NbPlayers >= 1);
let isSpinning = $ref(false);

let gameCode = $ref("");

async function launchGame(): Promise<void> {
  isSpinning = true;
  const res = await controller.LaunchGame(launchOptions);
  isSpinning = false;

  if (res === undefined) {
    return;
  }

  gameCode = res.URL.slice(res.URL.length - pinLength);
}
</script>

<style></style>
