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
            :items="selectItems.map((e) => e.text)"
            :model-value="selectItems.find(e => e.value == props.modelValue.Criterion)!.text"
            @update:model-value="onCriterionUpdate"
          >
          </v-select>
        </v-col>
      </v-row>
    </v-card-text>
  </v-card>
  <figure-block-vue
    v-model="props.modelValue.Figure"
    @update:model-value="emit('update:modelValue', props.modelValue)"
    :available-parameters="props.availableParameters"
  ></figure-block-vue>
</template>

<script setup lang="ts">
import type {
  FigureVectorPairFieldBlock,
  Variable,
} from "@/controller/api_gen";
import {
  VectorPairCriterion,
  VectorPairCriterionLabels,
} from "@/controller/api_gen";
import FigureBlockVue from "./FigureBlock.vue";

interface Props {
  modelValue: FigureVectorPairFieldBlock;
  availableParameters: Variable[];
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (event: "update:modelValue", value: FigureVectorPairFieldBlock): void;
}>();

const selectItems = Object.entries(VectorPairCriterionLabels).map((e) => ({
  text: e[1],
  value: Number(e[0]) as VectorPairCriterion,
}));

function onCriterionUpdate(s: string) {
  props.modelValue.Criterion = selectItems.find((e) => e.text == s)!.value;
  emit("update:modelValue", props.modelValue);
}
</script>

<style scoped></style>
