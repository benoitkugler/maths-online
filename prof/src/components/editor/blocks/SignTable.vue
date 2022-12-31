<template>
  <base-sign-table
    :model-value="props.modelValue"
    @update:model-value="emit('update:modelValue', props.modelValue)"
  ></base-sign-table>
</template>

<script setup lang="ts">
import type { SignTableBlock, Variable } from "@/controller/api_gen";
import BaseSignTable from "./BaseSignTable.vue";

interface Props {
  modelValue: SignTableBlock;
  availableParameters: Variable[];
}
const props = defineProps<Props>();

const emit = defineEmits<{
  (event: "update:modelValue", value: SignTableBlock): void;
}>();
</script>

<style scoped></style>
