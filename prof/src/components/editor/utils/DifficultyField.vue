<template>
  <v-select
    style="max-width: 200px"
    :items="[
      DifficultyTag.DiffEmpty,
      DifficultyTag.Diff1,
      DifficultyTag.Diff2,
      DifficultyTag.Diff3,
    ]"
    label="Difficulté"
    variant="outlined"
    density="compact"
    :model-value="props.modelValue"
    @update:model-value="(tag) => emit('update:model-value', tag)"
    :color="tagColor(undefined)"
    :disabled="props.readonly"
    hint="Permet de selectionner les ressources par difficulté dans les activités."
    persistent-hint
  >
  </v-select>
</template>

<script setup lang="ts">
import { DifficultyTag } from "@/controller/api_gen";
import { tagColor } from "@/controller/editor";

interface Props {
  modelValue: DifficultyTag;
  readonly: boolean;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "update:model-value", v: DifficultyTag): void;
}>();
</script>

<style scoped></style>
