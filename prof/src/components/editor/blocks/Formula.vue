<template>
  <Container
    @delete="emit('delete')"
    @swap="(o, t) => emit('swap', o, t)"
    :index="props.Pos.index"
    :nb-blocks="props.Pos.nbBlocks"
  >
    <InterpolatedText v-model="props.Data.Parts"></InterpolatedText>
    <small class="text-grey mt-1"
      >Interprété comme du code LaTeX. Insérer une expression avec : #{2x +
      1}</small
    >
  </Container>
</template>

<script setup lang="ts">
import { itemize } from "@/controller/editor";
import type { FormulaBlock } from "@/controller/exercice_gen";
import { computed } from "@vue/runtime-core";
import InterpolatedText from "../InterpolatedText.vue";
import type { ContainerProps } from "./Container.vue";
import Container from "./Container.vue";

interface Props {
  Data: FormulaBlock;
  Pos: ContainerProps;
}
const props = defineProps<Props>();

const colorByKind = ["", "green", "orange"];

const emit = defineEmits<{
  (e: "delete"): void;
  (e: "swap", origin: number, target: number): void;
}>();

const spans = computed(() =>
  itemize(props.Data.Parts).map(chunk => ({
    color: colorByKind[chunk.Kind],
    text: chunk.Content
  }))
);
</script>

<style></style>
