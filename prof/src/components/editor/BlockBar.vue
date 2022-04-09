<template>
  <v-list style="height: 70vh">
    <v-list-subheader>Enoncés</v-list-subheader>
    <v-list-item
      dense
      class="py-0"
      v-for="kind in staticKinds"
      @click="emit('add', kind)"
    >
      {{ labels[kind] }}
    </v-list-item>
    <v-divider></v-divider>
    <v-list-subheader>Champs de réponse</v-list-subheader>
    <v-list-item
      dense
      class="py-0"
      v-for="kind in fieldKinds"
      @click="emit('add', kind)"
    >
      {{ labels[kind] }}
    </v-list-item>
  </v-list>
</template>

<script setup lang="ts">
import { BlockKindLabels } from "@/controller/editor";
import { BlockKind } from "@/controller/exercice_gen";

const emit = defineEmits<{
  (e: "add", kind: BlockKind): void;
}>();

const labels = BlockKindLabels;

const staticKinds = [
  BlockKind.TextBlock,
  BlockKind.FormulaBlock,
  BlockKind.FigureBlock,
  BlockKind.FunctionGraphBlock,
  BlockKind.FunctionVariationGraphBlock,
  BlockKind.VariationTableBlock,
  BlockKind.SignTableBlock,
  BlockKind.TableBlock
];

const fieldKinds = [
  BlockKind.NumberFieldBlock,
  BlockKind.FormulaFieldBlock,
  BlockKind.RadioFieldBlock
  //TODO
];
</script>
