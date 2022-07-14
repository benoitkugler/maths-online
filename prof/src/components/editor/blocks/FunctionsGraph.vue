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
                    :color="expressionColor"
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

  <v-card color="secondary" class="my-1">
    <v-row no-gutters>
      <v-col align-self="center" md="9">
        <v-card-subtitle> Surfaces colorées </v-card-subtitle>
      </v-col>
      <v-col md="3" style="text-align: right">
        <v-btn
          icon
          @click="addArea"
          title="Ajouter une surface colorée entre deux courbes."
          size="x-small"
          class="mr-2 my-2"
          :disabled="functionsNamesItems.length < 2"
        >
          <v-icon icon="mdi-plus" color="green" size="small"></v-icon>
        </v-btn>
      </v-col>
    </v-row>
    <v-list>
      <div v-for="(area, index) in props.modelValue.Areas" :key="index">
        <v-list-item>
          <v-row>
            <v-col cols="2" align-self="center">
              <btn-color-picker v-model="area.Color"></btn-color-picker>
            </v-col>
            <v-col cols="9" align-self="center">
              <v-row>
                <v-col cols="7">
                  <v-combobox
                    density="compact"
                    variant="outlined"
                    hide-details
                    hide-no-data
                    label="Fonction 1"
                    :items="functionsNamesItems"
                    item
                    :model-value="nameToSelection(area.Top)"
                    @update:model-value="
                      (s) => (area.Top = nameFromSelection(s))
                    "
                    :color="expressionColor"
                  ></v-combobox>
                </v-col>
                <v-col cols="5">
                  <ExpressionField
                    label="Xmin"
                    v-model="area.Left"
                  ></ExpressionField>
                </v-col>
              </v-row>
              <v-row>
                <v-col cols="7">
                  <v-combobox
                    density="compact"
                    variant="outlined"
                    hide-details
                    hide-no-data
                    label="Fonction 2"
                    :items="functionsNamesItems"
                    :model-value="nameToSelection(area.Bottom)"
                    @update:model-value="
                      (s) => (area.Bottom = nameFromSelection(s))
                    "
                    :color="expressionColor"
                  ></v-combobox>
                </v-col>
                <v-col cols="5">
                  <ExpressionField
                    label="Xmax"
                    v-model="area.Right"
                  ></ExpressionField>
                </v-col>
              </v-row>
            </v-col>
            <v-col md="1" align-self="center" class="pl-1 pr-0">
              <v-btn icon size="x-small" @click="deleteArea(index)">
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
import { $computed } from "vue/macros";
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

const functionsNamesItems = $computed(() => {
  const set: { [key: string]: boolean } = {};
  props.modelValue.FunctionExprs?.forEach(
    (fn) => (set[fn.Decoration.Label] = true)
  );
  props.modelValue.FunctionVariations?.forEach((fn) => (set[fn.Label] = true));
  set[abscisseAxis] = true; // add empty string as horizontal axis
  const out = Object.keys(set);
  out.sort((a, b) => a.localeCompare(b));
  return out;
});

const abscisseAxis = "<Axe des abscisses>";

function nameFromSelection(s: string) {
  return s == abscisseAxis ? "" : s;
}
function nameToSelection(s: string) {
  return s == "" ? abscisseAxis : s;
}

function addArea() {
  props.modelValue.Areas?.push({
    Color: lastColorUsed.color,
    Top: nameFromSelection(functionsNamesItems[1]),
    Bottom: nameFromSelection(functionsNamesItems[0]),
    Left: "0",
    Right: "2",
  });
  emit("update:modelValue", props.modelValue);
}

function deleteArea(index: number) {
  props.modelValue.Areas?.splice(index, 1);
  emit("update:modelValue", props.modelValue);
}

const expressionColor = ExpressionColor;
</script>

<style scoped>
:deep(input) {
  width: 100%;
}
</style>
