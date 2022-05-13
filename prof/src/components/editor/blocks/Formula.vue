<template>
  <InterpolatedText
    v-model="props.modelValue.Parts"
    @update:model-value="emit('update:modelValue', props.modelValue)"
  ></InterpolatedText>
  <small class="text-grey mt-1 d-block"
    >Interprété comme du code LaTeX. Insérer une expression avec : &2x +
    1&</small
  >
</template>

<script setup lang="ts">
import type { FormulaBlock } from "@/controller/exercice_gen";
import InterpolatedText from "../utils/InterpolatedText.vue";

interface Props {
  modelValue: FormulaBlock;
}
const props = defineProps<Props>();

const emit = defineEmits<{
  (event: "update:modelValue", value: FormulaBlock): void;
}>();
</script>

<style></style>
