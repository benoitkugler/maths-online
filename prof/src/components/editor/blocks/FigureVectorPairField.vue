<template>
  <v-card class="my-1">
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
            @update:model-value="
              s =>
                (props.modelValue.Criterion = selectItems.find(
                  e => e.text == s
                )!.value)
            "
          >
          </v-select>
        </v-col>
      </v-row>
    </v-card-text>
  </v-card>
  <figure-vue
    v-model="props.modelValue.Figure"
    :available-parameters="props.availableParameters"
  ></figure-vue>
</template>

<script setup lang="ts">
import { colorByKind } from "@/controller/editor";
import type {
  FigureVectorPairFieldBlock,
  Variable
} from "@/controller/exercice_gen";
import {
  TextKind,
  VectorPairCriterion,
  VectorPairCriterionLabels
} from "@/controller/exercice_gen";
import FigureVue from "./Figure.vue";

interface Props {
  modelValue: FigureVectorPairFieldBlock;
  availableParameters: Variable[];
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (event: "update:modelValue", value: FigureVectorPairFieldBlock): void;
}>();

const selectItems = Object.entries(VectorPairCriterionLabels).map(e => ({
  text: e[1],
  value: Number(e[0]) as VectorPairCriterion
}));

const expressionColor = colorByKind[TextKind.Expression];
</script>

<style scoped></style>
