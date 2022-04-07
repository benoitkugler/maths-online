<template>
  <container
    @delete="emit('delete')"
    @swap="(o, t) => emit('swap', o, t)"
    :index="props.Pos.index"
    :nb-blocks="props.Pos.nbBlocks"
  >
    <template v-slot:toolbar>
      <v-checkbox
        class="pr-2"
        hide-details
        label="Comme conseil"
        v-model="props.Data.IsHint"
      ></v-checkbox>
    </template>
    <interpolated-text
      v-model="props.Data.Parts"
      class="px-2"
    ></interpolated-text>
    <small class="text-grey mt-1"
      >Insérer du code LaTeX avec : $\frac{a}{b}$. Insérer une expression avec :
      #{2x + 1}</small
    >
  </container>
</template>

<script setup lang="ts">
import type { TextBlock } from "@/controller/exercice_gen";
import InterpolatedText from "../InterpolatedText.vue";
import type { ContainerProps } from "./Container.vue";
import Container from "./Container.vue";

interface Props {
  Data: TextBlock;
  Pos: ContainerProps;
}
const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "delete"): void;
  (e: "swap", origin: number, target: number): void;
}>();
</script>

<style></style>
