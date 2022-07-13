<template>
  <base-variation-table
    :model-value="props.modelValue"
    @update:model-value="emit('update:modelValue', props.modelValue)"
    description="Fonction définie par les expressions pour x et f(x) :"
  ></base-variation-table>
</template>

<script setup lang="ts">
import type { FunctionVariationGraphBlock } from "@/controller/api_gen";
import BaseVariationTable from "./BaseVariationTable.vue";

interface Props {
  modelValue: FunctionVariationGraphBlock;
}
const props = defineProps<Props>();

const emit = defineEmits<{
  (event: "update:modelValue", value: FunctionVariationGraphBlock): void;
}>();
</script>

<style></style>
