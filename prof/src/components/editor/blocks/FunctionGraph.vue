<template>
  <v-card color="secondary" class="my-1">
    <v-row no-gutters>
      <v-col align-self="center" md="9">
        <v-card-subtitle> Fonctions </v-card-subtitle>
      </v-col>
      <v-col md="3" style="text-align: right">
        <v-btn
          icon
          @click="addFunction"
          title="Ajouter une fonction"
          size="x-small"
          class="mr-2 my-2"
        >
          <v-icon icon="mdi-plus" color="green" size="small"></v-icon>
        </v-btn>
      </v-col>
    </v-row>
    <v-list>
      <div v-for="(fn, index) in props.modelValue.Functions" :key="index">
        <v-list-item>
          <v-row>
            <v-col cols="10">
              <v-row>
                <v-col cols="3" align-self="center">
                  <v-text-field
                    variant="outlined"
                    density="compact"
                    v-model="fn.Decoration.Label"
                    label="Légende"
                    hide-details
                  ></v-text-field>
                </v-col>
                <v-col md="2" align-self="center">
                  ( {{ variableToString(fn.Variable) }} ) =
                </v-col>
                <v-col cols="7" align-self="center">
                  <v-text-field
                    variant="outlined"
                    density="compact"
                    v-model="fn.Function"
                    label="Expression de la fonction"
                    hide-details
                    :color="color"
                  ></v-text-field>
                </v-col>
                <v-col md="4" align-self="center">
                  <ExpressionField
                    v-model="fn.From"
                    label="Xmin"
                  ></ExpressionField>
                </v-col>
                <v-col md="4">
                  <ExpressionField
                    v-model="fn.To"
                    label="Xmax"
                  ></ExpressionField>
                </v-col>
                <v-col>
                  <BtnColorPicker
                    v-model="fn.Decoration.Color"
                  ></BtnColorPicker>
                </v-col>
              </v-row>
            </v-col>

            <v-col cols="2" align-self="center">
              <v-btn icon size="x-small" flat @click="deleteFunction(index)">
                <v-icon icon="mdi-delete" color="red"></v-icon>
              </v-btn>
            </v-col>
          </v-row>
        </v-list-item>
        <v-divider></v-divider>
      </div>
    </v-list>
  </v-card>
</template>

<script setup lang="ts">
import type { FunctionGraphBlock } from "@/controller/api_gen";
import {
ExpressionColor,
lastColorUsed,
variableToString,
xRune
} from "@/controller/editor";
import ExpressionField from "../utils/ExpressionField.vue";
import BtnColorPicker from "./BtnColorPicker.vue";

interface Props {
  modelValue: FunctionGraphBlock;
}
const props = defineProps<Props>();

const emit = defineEmits<{
  (event: "update:modelValue", value: FunctionGraphBlock): void;
}>();

function addFunction() {
  props.modelValue.Functions?.push({
    Function: "x^2",
    Decoration: {
      Label: "f",
      Color: lastColorUsed.color,
    },
    Variable: { Name: xRune, Indice: "" },
    From: "-4",
    To: "4",
  });
  emit("update:modelValue", props.modelValue);
}

function deleteFunction(index: number) {
  props.modelValue.Functions?.splice(index, 1);
  emit("update:modelValue", props.modelValue);
}

const color = ExpressionColor;
</script>

<style scoped>
:deep(input) {
  width: 100%;
}
</style>
