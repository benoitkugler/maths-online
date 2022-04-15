<template>
  Les valeurs f(x) sont arrondies à l'unité avant d'être comparées.
  <v-row class="mt-2">
    <v-col md="2" align-self="center">
      <v-text-field
        draggable="false"
        variant="outlined"
        density="compact"
        v-model="props.modelValue.Label"
        label="Nom"
        hide-details
      ></v-text-field>
    </v-col>
    <v-col md="3" align-self="center">
      <VariableField
        prefix="("
        suffix=") ="
        v-model="props.modelValue.Variable"
        label="Variable"
      >
      </VariableField>
    </v-col>
    <v-col md="7" align-self="center">
      <v-text-field
        variant="outlined"
        density="compact"
        v-model="props.modelValue.Function"
        label="Expression de la fonction"
        hide-details
        :color="color"
      ></v-text-field>
    </v-col>
  </v-row>
  <v-row>
    <v-col md="12">
      <IntListField
        :model-value="props.modelValue.XGrid || []"
        @update:model-value="g => (props.modelValue.XGrid = g)"
        label="Valeurs de X"
      ></IntListField>
    </v-col>
  </v-row>
</template>

<script setup lang="ts">
import { ExpressionColor } from "@/controller/editor";
import type { FunctionPointsFieldBlock } from "@/controller/exercice_gen";
import VariableField from "../utils/VariableField.vue";
import IntListField from "./IntListField.vue";

interface Props {
  modelValue: FunctionPointsFieldBlock;
}
const props = defineProps<Props>();

const emit = defineEmits<{
  (event: "update:modelValue", value: FunctionPointsFieldBlock): void;
}>();

const color = ExpressionColor;
</script>

<style></style>
