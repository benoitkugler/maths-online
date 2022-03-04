<template>
  <v-card class="ma-3">
    <v-card-title>Trivial Poursuit</v-card-title>
    <v-card-subtitle>Configurer une partie de Trivial Poursuit</v-card-subtitle>

    <v-form class="mx-4 my-2">
      <v-text-field
        label="Nombre de joueurs"
        type="number"
        v-model.number="nbPlayers"
      ></v-text-field>
      <v-row>
        <v-spacer></v-spacer>
        <v-col>
          <v-btn block @click="launchGame" :disabled="!isValid" color="success">
            Lancer la partie
          </v-btn>
        </v-col>
      </v-row>

      <v-alert class="px-4 my-2" :model-value="gameURL != ''" color="info">
        Rejoindre la partie à cette adresse :
        <v-btn class="mx-2" :href="gameURL" target="_blank">{{
          gameURL
        }}</v-btn>
      </v-alert>
    </v-form>
  </v-card>
</template>

<script setup lang="ts">
import { controller } from "@/controller/controller";
import { ref } from "@vue/reactivity";
import { computed } from "@vue/runtime-core";

const nbPlayers = ref(0);

const isValid = computed(() => !isSpinning.value && nbPlayers.value > 0);
const isSpinning = ref(false);

const gameURL = ref("");

async function launchGame(): Promise<void> {
  isSpinning.value = true;
  const res = await controller.LaunchGame({ NbPlayers: nbPlayers.value });
  isSpinning.value = false;

  if (res === undefined) {
    return;
  }

  gameURL.value = res.URL;
}
</script>

<style></style>
