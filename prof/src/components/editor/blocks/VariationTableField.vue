<template>
  <base-variation-table
    :model-value="props.modelValue.Answer"
    @update:model-value="emit('update:modelValue', props.modelValue)"
    description="Le sens de variation est déduit des valeurs de f(x)."
  ></base-variation-table>
</template>

<script setup lang="ts">
import type { Variable, VariationTableFieldBlock } from "@/controller/api_gen";
import BaseVariationTable from "./BaseVariationTable.vue";

interface Props {
  modelValue: VariationTableFieldBlock;
  availableParameters: Variable[];
}
const props = defineProps<Props>();

const emit = defineEmits<{
  (event: "update:modelValue", value: VariationTableFieldBlock): void;
}>();
</script>

<style></style>
