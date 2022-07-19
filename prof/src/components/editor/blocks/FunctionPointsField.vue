<template>
  <small>
    Les valeurs f(x) doivent être entière pour pouvoir être placée sur la
    grille.
  </small>
  <v-row class="mt-2 fix-input-width">
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
  <v-row class="fix-input-width">
    <v-col md="12">
      <ExpressionListField
        :model-value="props.modelValue.XGrid || []"
        @update:model-value="updateXGrid"
        label="Valeurs de X"
      ></ExpressionListField>
    </v-col>
  </v-row>
</template>

<script setup lang="ts">
import type { FunctionPointsFieldBlock } from "@/controller/api_gen";
import { ExpressionColor } from "@/controller/editor";
import ExpressionListField from "../utils/ExpressionListField.vue";
import VariableField from "../utils/VariableField.vue";

interface Props {
  modelValue: FunctionPointsFieldBlock;
}
const props = defineProps<Props>();

const emit = defineEmits<{
  (event: "update:modelValue", value: FunctionPointsFieldBlock): void;
}>();

function updateXGrid(g: string[]) {
  props.modelValue.XGrid = g;
  emit("update:modelValue", props.modelValue);
}

const color = ExpressionColor;
</script>

<style scoped>
.fix-input-width:deep(input) {
  width: 100%;
}
</style>
