<template>
  <Container @delete="emit('delete')">
    <v-text-field
      v-model="props.Data.Parts"
      hint="Interprété comme du code LaTeX. Insérer une expression avec : #{2x + 1}"
    ></v-text-field>
  </Container>
</template>

<script setup lang="ts">
import { itemize } from "@/controller/editor";
import type { FormulaBlock } from "@/controller/exercice_gen";
import { computed } from "@vue/runtime-core";
import Container from "./Container.vue";

interface Props {
  Data: FormulaBlock;
}
const props = defineProps<Props>();

const colorByKind = ["", "green", "orange"];

const emit = defineEmits<{
  (e: "delete"): void;
}>();

const spans = computed(() =>
  itemize(props.Data.Parts).map(chunk => ({
    color: colorByKind[chunk.Kind],
    text: chunk.Content
  }))
);
</script>

<style>
.v-field__field {
  padding-top: 0;
}

.v-field__input {
  padding-left: 4px;
}
</style>
