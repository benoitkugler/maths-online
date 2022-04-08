<template>
  <v-row>
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
    <v-col md="4" align-self="center">
      <v-text-field
        variant="outlined"
        density="compact"
        v-model="props.modelValue.Function"
        label="Expression de la fonction"
        hide-details
      ></v-text-field>
    </v-col>
    <v-col md="3" align-self="center">
      <v-row>
        <v-col md="12">
          <v-text-field
            variant="outlined"
            density="compact"
            v-model.number="props.modelValue.Range[0]"
            label="Xmin"
            hide-details
          ></v-text-field>
        </v-col>
        <v-col md="12">
          <v-text-field
            variant="outlined"
            density="compact"
            v-model.number="props.modelValue.Range[1]"
            label="Xmax"
            hide-details
          ></v-text-field>
        </v-col>
      </v-row>
    </v-col>
  </v-row>
</template>

<script setup lang="ts">
import type { FunctionGraphBlock } from "@/controller/exercice_gen";
import VariableField from "../VariableField.vue";

interface Props {
  modelValue: FunctionGraphBlock;
}
const props = defineProps<Props>();

const emit = defineEmits<{
  (event: "update:modelValue", value: FunctionGraphBlock): void;
}>();
</script>

<style></style>
