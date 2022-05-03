<template>
  <v-form class="mx-4 my-2">
    <v-text-field
      variant="outlined"
      density="comfortable"
      label="Nombre de joueurs par groupe"
      type="number"
      min="1"
      v-model.number="props.modelValue.MaxPlayersPerGroup"
      @update:model-value="emit('update:model-value', props.modelValue)"
    ></v-text-field>

    <v-text-field
      variant="outlined"
      density="comfortable"
      label="Nombre total de joueurs"
      type="number"
      min="0"
      v-model.number="props.modelValue.TotalPlayersNumber"
      @update:model-value="emit('update:model-value', props.modelValue)"
      hint="Le nombre total de joueurs est nécessaire pour déterminer correctement le lancement du dernier groupe."
    ></v-text-field>
  </v-form>
</template>

<script setup lang="ts">
import type { RandomGroupStrategy } from "@/controller/api_gen";

interface Props {
  modelValue: RandomGroupStrategy;
}

const emit = defineEmits<{
  (e: "update:model-value", v: RandomGroupStrategy): void;
}>();

const props = defineProps<Props>();
</script>

<style scoped></style>
