<template>
  <v-list style="height: 70vh">
    <v-list-subheader><h3 class="text-purple">Enoncés</h3></v-list-subheader>
    <v-list-item
      rounded
      dense
      class="py-0 bg-purple-lighten-4 ma-1"
      v-for="kind in staticKinds"
      @click="emit('add', kind)"
      :key="kind"
    >
      {{ labels[kind].label }}
    </v-list-item>
    <v-divider></v-divider>
    <v-list-subheader
      ><h3 class="text-pink">Champs de réponse</h3></v-list-subheader
    >
    <v-list-item
      rounded
      dense
      class="py-0 bg-pink-lighten-4 ma-1"
      v-for="kind in fieldKinds"
      :key="kind"
      @click="emit('add', kind)"
    >
      {{ labels[kind].label }}
    </v-list-item>
  </v-list>
</template>

<script setup lang="ts">
import { BlockKindLabels, sortedBlockKindLabels } from "@/controller/editor";
import type { BlockKind } from "@/controller/exercice_gen";

const emit = defineEmits<{
  (e: "add", kind: BlockKind): void;
}>();

const labels = BlockKindLabels;

const staticKinds = sortedBlockKindLabels
  .filter((k) => !k[1].isAnswerField)
  .map((k) => k[0]);
const fieldKinds = sortedBlockKindLabels
  .filter((k) => k[1].isAnswerField)
  .map((k) => k[0]);
</script>
