<template>
  <base-sign-table
    :model-value="props.modelValue.Answer"
    @update:model-value="emit('update:modelValue', props.modelValue)"
  ></base-sign-table>
</template>

<script setup lang="ts">
import type { SignTableFieldBlock } from "@/controller/api_gen";
import BaseSignTable from "./BaseSignTable.vue";

interface Props {
  modelValue: SignTableFieldBlock;
}
const props = defineProps<Props>();

const emit = defineEmits<{
  (event: "update:modelValue", value: SignTableFieldBlock): void;
}>();
</script>

<style></style>
