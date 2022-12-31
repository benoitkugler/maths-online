<template>
  <base-variation-table
    :model-value="props.modelValue"
    @update:model-value="emit('update:modelValue', props.modelValue)"
    description="Le sens de variation est déduit des valeurs de f(x)."
  ></base-variation-table>
</template>

<script setup lang="ts">
import type { Variable, VariationTableBlock } from "@/controller/api_gen";
import BaseVariationTable from "./BaseVariationTable.vue";

interface Props {
  modelValue: VariationTableBlock;
  availableParameters: Variable[];
}
const props = defineProps<Props>();

const emit = defineEmits<{
  (event: "update:modelValue", value: VariationTableBlock): void;
}>();
</script>

<style></style>
