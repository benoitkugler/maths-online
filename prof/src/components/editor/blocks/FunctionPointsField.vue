<template>
  <small>
    Les valeurs prises doivent être entières pour pouvoir être placées sur la
    grille.
  </small>
  <v-row class="mt-2 fix-input-width">
    <v-col md="2" align-self="center">
      <v-text-field
        draggable="false"
        variant="outlined"
        density="compact"
        v-model="props.modelValue.Label"
        @update:model-value="emitUpdate"
        label="Nom"
        hide-details
      ></v-text-field>
    </v-col>
    <v-col md="2" align-self="center">
      <VariableField
        v-model="props.modelValue.Variable"
        @update:model-value="emitUpdate"
        label="Variable"
      >
      </VariableField>
    </v-col>
    <v-col md="8" align-self="center">
      <v-text-field
        variant="outlined"
        density="compact"
        v-model="props.modelValue.Function"
        @update:model-value="emitUpdate"
        label="Expression de la fonction ou de la suite"
        hide-details
        :color="color"
      ></v-text-field>
    </v-col>
  </v-row>
  <v-row class="fix-input-width">
    <v-col md="12">
      <ExpressionListField
        :model-value="props.modelValue.XGrid || []"
        @update:model-value="updateXGrid"
        label="Valeurs de X"
      ></ExpressionListField>
    </v-col>
  </v-row>
  <v-row no-gutters>
    <v-col md="12">
      <v-checkbox
        density="compact"
        v-model="props.modelValue.IsDiscrete"
        @update:model-value="emitUpdate"
        label="Afficher comme une suite"
        hide-details
      ></v-checkbox>
    </v-col>
  </v-row>
</template>

<script setup lang="ts">
import type { FunctionPointsFieldBlock, Variable } from "@/controller/api_gen";
import { ExpressionColor } from "@/controller/editor";
import ExpressionListField from "../utils/ExpressionListField.vue";
import VariableField from "../utils/VariableField.vue";

interface Props {
  modelValue: FunctionPointsFieldBlock;
  availableParameters: Variable[];
}
const props = defineProps<Props>();

const emit = defineEmits<{
  (event: "update:modelValue", value: FunctionPointsFieldBlock): void;
}>();

function emitUpdate() {
  emit("update:modelValue", props.modelValue);
}

function updateXGrid(g: string[]) {
  props.modelValue.XGrid = g;
  emitUpdate();
}

const color = ExpressionColor;
</script>

<style scoped>
.fix-input-width:deep(input) {
  width: 100%;
}
</style>
