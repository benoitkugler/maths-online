<template>
  <InterpolatedText
    :model-value="props.modelValue.Content"
    @update:model-value="(s) => emit('update:modelValue', { Content: s })"
  ></InterpolatedText>
</template>

<script setup lang="ts">
import type { ProofStatement } from "@/controller/api_gen";
import InterpolatedText from "../../utils/InterpolatedText.vue";

interface Props {
  modelValue: ProofStatement;
}
const props = defineProps<Props>();

const emit = defineEmits<{
  (event: "update:modelValue", value: ProofStatement): void;
}>();
</script>

<style></style>
