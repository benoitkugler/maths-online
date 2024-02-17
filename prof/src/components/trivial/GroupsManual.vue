<template>
  <v-row class="my-2">
    <v-col>
      <v-text-field
        variant="outlined"
        density="compact"
        type="number"
        label="Nombre de groupes"
        hint="Une partie est créée par groupe."
        :model-value="props.modelValue.NbGroups"
        @update:model-value="
          (v) => emit('update:model-value', { NbGroups: Number(v) as Int})
        "
      ></v-text-field>
    </v-col>
  </v-row>
</template>

<script setup lang="ts">
import type { GroupsStrategyManual, Int } from "@/controller/api_gen";

interface Props {
  modelValue: GroupsStrategyManual;
}

const emit = defineEmits<{
  (e: "update:model-value", v: GroupsStrategyManual): void;
}>();

const props = defineProps<Props>();
</script>

<style scoped></style>
