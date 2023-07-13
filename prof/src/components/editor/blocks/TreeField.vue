<template>
  <TreeB
    :available-parameters="props.availableParameters"
    :model-value="props.modelValue.Answer"
    @update:model-value="
      (v) => {
        props.modelValue.Answer = v;
        emit('update:modelValue', props.modelValue);
      }
    "
  ></TreeB>
</template>

<script setup lang="ts">
import type { TreeFieldBlock, Variable } from "@/controller/api_gen";
import TreeB from "./TreeB.vue";

interface Props {
  modelValue: TreeFieldBlock;
  availableParameters: Variable[];
}
const props = defineProps<Props>();

const emit = defineEmits<{
  (event: "update:modelValue", value: TreeFieldBlock): void;
}>();
</script>

<style></style>
