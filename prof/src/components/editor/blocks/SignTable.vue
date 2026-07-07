<template>
  <base-sign-table
    :model-value="props.modelValue"
    @update:model-value="(v) => emit('update:modelValue', v)"
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
