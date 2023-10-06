<template>
  <v-card subtitle="Détails des succès">
    <v-card-text>
      <v-list density="compact" style="column-count: 2">
        <v-list-item
          v-for="(oc, event) in stats.Occurences"
          :key="event"
          :title="events[event].title"
          style="break-inside: avoid-column"
        >
          <template v-slot:append>
            <v-chip class="mx-2" :color="colors[events[event].kind]">{{
              oc
            }}</v-chip>
          </template>
        </v-list-item>
      </v-list>
    </v-card-text>
  </v-card>
</template>

<script setup lang="ts">
import type { Stats } from "@/controller/api_gen";

interface Props {
  stats: Stats;
}

const props = defineProps<Props>();

const enum EK {
  isytriv,
  homework,
  misc
}

const events = [
  { kind: EK.isytriv, title: "Créer une partie d'IsyTriv" }, // E_IsyTriv_Create
  { kind: EK.isytriv, title: "Réussir trois questions IsyTriv d'affilée" }, // E_IsyTriv_Streak3
  { kind: EK.isytriv, title: "Remporter une partie IsyTriv" }, // E_IsyTriv_Win
  { kind: EK.homework, title: "Terminer un exercice" }, // E_Homework_TaskDone
  { kind: EK.homework, title: "Terminer une feuille d'exercices" }, // E_Homework_TravailDone
  { kind: EK.misc, title: "Répondre correctement à une question" }, // E_All_QuestionRight
  { kind: EK.misc, title: "Répondre incorrectement à une question" }, // E_All_QuestionWrong
  { kind: EK.misc, title: "Modifier sa playlist" }, // E_Misc_SetPlaylist
  { kind: EK.misc, title: "Se connecter 3 jours de suite" }, // E_ConnectStreak3
  { kind: EK.misc, title: "Se connecter 7 jours de suite" }, // E_ConnectStreak7
  { kind: EK.misc, title: "Se connecter 30 jours de suite" } // E_ConnectStreak30
] as const;

const colors = ["blue", "green", "grey"];
</script>
