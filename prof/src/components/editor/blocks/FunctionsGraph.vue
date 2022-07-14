<template>
  <v-card color="secondary" class="my-1">
    <v-row no-gutters>
      <v-col align-self="center" md="9">
        <v-card-subtitle> Avec une expression </v-card-subtitle>
      </v-col>
      <v-col md="3" style="text-align: right">
        <v-btn
          icon
          @click="addFunctionExpr"
          title="Ajouter une fonction"
          size="x-small"
          class="mr-2 my-2"
        >
          <v-icon icon="mdi-plus" color="green" size="small"></v-icon>
        </v-btn>
      </v-col>
    </v-row>
    <v-list>
      <div v-for="(fn, index) in props.modelValue.FunctionExprs" :key="index">
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
              <v-btn
                icon
                size="x-small"
                flat
                @click="deleteFunctionExpr(index)"
              >
                <v-icon icon="mdi-delete" color="red"></v-icon>
              </v-btn>
            </v-col>
          </v-row>
        </v-list-item>
        <v-divider></v-divider>
      </div>
    </v-list>
  </v-card>

  <v-card color="secondary" class="my-1">
    <v-row no-gutters>
      <v-col align-self="center" md="9">
        <v-card-subtitle> Avec des variations </v-card-subtitle>
      </v-col>
      <v-col md="3" style="text-align: right">
        <v-btn
          icon
          @click="addFunctionVar"
          title="Ajouter une fonction"
          size="x-small"
          class="mr-2 my-2"
        >
          <v-icon icon="mdi-plus" color="green" size="small"></v-icon>
        </v-btn>
      </v-col>
    </v-row>
    <v-list>
      <div
        v-for="(fn, index) in props.modelValue.FunctionVariations"
        :key="index"
      >
        <v-list-item>
          <v-row>
            <v-col cols="10">
              <BaseVariationTable
                :model-value="fn"
                @update:model-value="
              (v) => (props.modelValue.FunctionVariations![index] = v)
            "
                description="Fonction définie par ses variations"
              ></BaseVariationTable>
            </v-col>
            <v-col cols="2" align-self="center" class="pr-0 pl-1">
              <v-btn icon size="x-small" flat @click="deleteFunctionVar(index)">
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
import type { FunctionsGraphBlock } from "@/controller/api_gen";
import {
  ExpressionColor,
  lastColorUsed,
  variableToString,
  xRune,
} from "@/controller/editor";
import ExpressionField from "../utils/ExpressionField.vue";
import BaseVariationTable from "./BaseVariationTable.vue";
import BtnColorPicker from "./BtnColorPicker.vue";

interface Props {
  modelValue: FunctionsGraphBlock;
}
const props = defineProps<Props>();

const emit = defineEmits<{
  (event: "update:modelValue", value: FunctionsGraphBlock): void;
}>();

function addFunctionExpr() {
  props.modelValue.FunctionExprs?.push({
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

function deleteFunctionExpr(index: number) {
  props.modelValue.FunctionExprs?.splice(index, 1);
  emit("update:modelValue", props.modelValue);
}

function addFunctionVar() {
  props.modelValue.FunctionVariations?.push({
    Label: "C_f",
    Xs: ["-5", "0", "5"],
    Fxs: ["-3", "2", "-1"],
  });
  emit("update:modelValue", props.modelValue);
}

function deleteFunctionVar(index: number) {
  props.modelValue.FunctionVariations?.splice(index, 1);
  emit("update:modelValue", props.modelValue);
}

const color = ExpressionColor;
</script>

<style scoped>
:deep(input) {
  width: 100%;
}
</style>
