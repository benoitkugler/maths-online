<template>
  <TagChip :tag="props.modelValue || 'Aucune'" @click="onClick"></TagChip>
</template>

<script setup lang="ts">
import { DifficultyTag } from "@/controller/api_gen";
import TagChip from "./TagChip.vue";

interface Props {
  modelValue: DifficultyTag;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "update:model-value", v: DifficultyTag): void;
}>();

function onClick() {
  switch (props.modelValue) {
    case DifficultyTag.DiffEmpty:
      emit("update:model-value", DifficultyTag.Diff1);
      break;
    case DifficultyTag.Diff1:
      emit("update:model-value", DifficultyTag.Diff2);
      break;
    case DifficultyTag.Diff2:
      emit("update:model-value", DifficultyTag.Diff3);
      break;
    case DifficultyTag.Diff3:
      emit("update:model-value", DifficultyTag.DiffEmpty);
      break;
  }
}
</script>

<style scoped></style>
