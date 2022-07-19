<template>
  <proof-sequence
    is-root
    :model-value="props.modelValue.Answer"
    @update:model-value="(v) => emit('update:modelValue', { Answer: v })"
  ></proof-sequence>
</template>

<script setup lang="ts">
import type { ProofFieldBlock } from "@/controller/api_gen";
import ProofSequence from "./proof/ProofSequence.vue";

interface Props {
  modelValue: ProofFieldBlock;
}
const props = defineProps<Props>();

const emit = defineEmits<{
  (event: "update:modelValue", value: ProofFieldBlock): void;
}>();
</script>

<style></style>
