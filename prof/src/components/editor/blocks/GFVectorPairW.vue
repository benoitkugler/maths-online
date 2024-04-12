<template>
  <v-card class="my-2">
    <v-card-subtitle class="bg-secondary py-3"
      >Type de vecteurs</v-card-subtitle
    >
    <v-card-text>
      <v-row>
        <v-col cols="12">
          <v-select
            persistent-hint
            label="Critère"
            hint="Caractéristique attendue pour les deux vecteurs à construire"
            variant="outlined"
            density="compact"
            :items="selectItems.map(e => e.text)"
            :model-value="selectItems.find(e => e.value == props.modelValue.Criterion)!.text"
            @update:model-value="onCriterionUpdate"
          >
          </v-select>
        </v-col>
      </v-row>
    </v-card-text>
  </v-card>
</template>

<script setup lang="ts">
import type { GFVectorPair } from "@/controller/api_gen";
import {
  VectorPairCriterion,
  VectorPairCriterionLabels
} from "@/controller/api_gen";

interface Props {
  modelValue: GFVectorPair;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (event: "update:modelValue", value: GFVectorPair): void;
}>();

function emitUpdate() {
  emit("update:modelValue", props.modelValue);
}

const selectItems = Object.entries(VectorPairCriterionLabels).map(e => ({
  text: e[1],
  value: Number(e[0]) as VectorPairCriterion
}));

function onCriterionUpdate(s: string) {
  props.modelValue.Criterion = selectItems.find(e => e.text == s)!.value;
  emitUpdate();
}
</script>

<style scoped></style>
