<template>
  <InterpolatedText
    :model-value="props.modelValue.Terms"
    @update:model-value="(s) => emit('update:modelValue', { Terms: s })"
  ></InterpolatedText>
</template>

<script setup lang="ts">
import type { ProofEquality } from "@/controller/api_gen";
import InterpolatedText from "../../utils/InterpolatedText.vue";

interface Props {
  modelValue: ProofEquality;
}
const props = defineProps<Props>();

const emit = defineEmits<{
  (event: "update:modelValue", value: ProofEquality): void;
}>();
</script>

<style></style>
