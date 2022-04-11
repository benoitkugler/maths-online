<template>
  <small class="text-grey mt-1">
    {{ props.description }}
  </small>
  <v-row>
    <v-col md="10" align-self="center">
      <v-table style="overflow-x: auto; max-width: 70vh">
        <tr>
          <th></th>
          <td
            v-for="(_, index) in props.modelValue.Xs"
            style="text-align: center; width: 40px"
          >
            <v-btn
              icon
              size="x-small"
              flat
              @click="removeColumn(index)"
              title="Supprimer la colonne"
            >
              <v-icon icon="mdi-close" color="red"></v-icon>
            </v-btn>
          </td>
        </tr>
        <tr>
          <th>x</th>
          <td
            v-for="(x, index) in props.modelValue.Xs"
            style="text-align: center; width: 80px"
          >
            <v-text-field
              class="centered-input"
              variant="outlined"
              density="compact"
              :model-value="x"
              @update:model-value="s => props.modelValue.Xs![index] = s"
              hide-details
              :color="expressionColor"
            ></v-text-field>
          </td>
        </tr>
        <tr>
          <th class="px-2">f(x)</th>
          <td
            v-for="(fx, index) in props.modelValue.Fxs"
            style="text-align: center; width: 80px"
          >
            <v-text-field
              class="centered-input"
              variant="outlined"
              density="compact"
              :model-value="fx"
              @update:model-value="s => props.modelValue.Fxs![index] = s"
              hide-details
              :color="expressionColor"
            ></v-text-field>
          </td>
        </tr>
      </v-table>
    </v-col>
    <v-col md="2" align-self="center">
      <v-btn
        icon
        @click="addColumn"
        title="Ajouter une colonne"
        size="x-small"
        class="mr-2 my-2"
      >
        <v-icon icon="mdi-plus" color="green"></v-icon>
      </v-btn>
    </v-col>
  </v-row>
</template>

<script setup lang="ts">
import { ExpressionColor } from "@/controller/editor";
import type { FunctionVariationGraphBlock } from "@/controller/exercice_gen";

interface Props {
  modelValue: FunctionVariationGraphBlock;
  description: string;
}
const props = defineProps<Props>();

const emit = defineEmits<{
  (event: "update:modelValue", value: FunctionVariationGraphBlock): void;
}>();

const expressionColor = ExpressionColor;

function addColumn() {
  props.modelValue.Xs?.push("5");
  props.modelValue.Fxs?.push("5");
}

function removeColumn(index: number) {
  props.modelValue.Xs?.splice(index, 1);
  props.modelValue.Fxs?.splice(index, 1);
}
</script>

<style>
.centered-input:deep(input) {
  text-align: center;
}
</style>
