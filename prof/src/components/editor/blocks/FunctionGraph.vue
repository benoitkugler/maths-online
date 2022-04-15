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
      <div v-for="(fn, index) in props.modelValue.Functions">
        <v-list-item>
          <v-row>
            <v-col cols="10">
              <v-row>
                <v-col md="2" align-self="center">
                  <v-text-field
                    draggable="false"
                    variant="outlined"
                    density="compact"
                    v-model="fn.Decoration.Label"
                    label="Nom"
                    hide-details
                  ></v-text-field>
                </v-col>
                <v-col md="4" align-self="center">
                  <VariableField
                    prefix="("
                    suffix=") ="
                    v-model="fn.Variable"
                    label="Variable"
                  >
                  </VariableField>
                </v-col>
                <v-col cols="6" align-self="center">
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
                  <v-text-field
                    variant="outlined"
                    density="compact"
                    v-model.number="fn.Range[0]"
                    label="Xmin"
                    hide-details
                  ></v-text-field>
                </v-col>
                <v-col md="4">
                  <v-text-field
                    variant="outlined"
                    density="compact"
                    v-model.number="fn.Range[1]"
                    label="Xmax"
                    hide-details
                  ></v-text-field>
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
import { ExpressionColor, xRune } from "@/controller/editor";
import type { FunctionGraphBlock } from "@/controller/exercice_gen";
import VariableField from "../VariableField.vue";
import BtnColorPicker from "./BtnColorPicker.vue";

interface Props {
  modelValue: FunctionGraphBlock;
}
const props = defineProps<Props>();

const emit = defineEmits<{
  (event: "update:modelValue", value: FunctionGraphBlock): void;
}>();

// TODO: add color field

function addFunction() {
  props.modelValue.Functions?.push({
    Function: "x^2",
    Decoration: {
      Label: "f",
      Color: ""
    },
    Variable: { Name: xRune, Indice: "" },
    Range: [-4, 4]
  });
  emit("update:modelValue", props.modelValue);
}

function deleteFunction(index: number) {
  props.modelValue.Functions?.splice(index, 1);
  emit("update:modelValue", props.modelValue);
}

const color = ExpressionColor;
</script>

<style></style>
