<template>
  <!-- Fonctions -->
  <v-card
    color="secondary"
    class="my-1"
    subtitle="Fonctions (définies par une expression)"
  >
    <template v-slot:append>
      <v-btn
        icon
        @click="addFunctionExpr"
        title="Ajouter une fonction"
        size="x-small"
      >
        <v-icon icon="mdi-plus" color="green"></v-icon>
      </v-btn>
    </template>
    <v-list>
      <div v-for="(fn, index) in props.modelValue.FunctionExprs" :key="index">
        <v-list-item>
          <v-row>
            <v-col>
              <v-row>
                <v-col cols="3" class="mb-4">
                  <interpolated-text
                    v-model="fn.Decoration.Label"
                    @update:model-value="emitUpdate"
                    force-latex
                    center
                    label="Légende"
                  ></interpolated-text>
                </v-col>
                <v-col md="2" align-self="center">
                  ( {{ variableToString(fn.Variable) }} ) =
                </v-col>
                <v-col cols="7" align-self="center">
                  <v-text-field
                    variant="outlined"
                    density="compact"
                    v-model="fn.Function"
                    @update:model-value="emitUpdate"
                    label="Expression de la fonction"
                    hide-details
                    :color="expressionColor"
                  ></v-text-field>
                </v-col>
                <v-col md="4" align-self="center">
                  <ExpressionField
                    v-model="fn.From"
                    @update:model-value="emitUpdate"
                    label="Xmin"
                  ></ExpressionField>
                </v-col>
                <v-col md="4">
                  <ExpressionField
                    v-model="fn.To"
                    @update:model-value="emitUpdate"
                    label="Xmax"
                  ></ExpressionField>
                </v-col>
                <v-col>
                  <BtnColorPicker
                    v-model="fn.Decoration.Color"
                    @update:model-value="emitUpdate"
                  ></BtnColorPicker>
                </v-col>
              </v-row>
            </v-col>

            <v-col cols="auto" align-self="center">
              <v-btn icon size="x-small" @click="deleteFunctionExpr(index)">
                <v-icon icon="mdi-delete" color="red"></v-icon>
              </v-btn>
            </v-col>
          </v-row>
        </v-list-item>
        <v-divider></v-divider>
      </div>
    </v-list>
  </v-card>

  <!-- Function from variations -->
  <v-card
    color="secondary"
    class="my-1"
    subtitle="Fonctions (définies par des variations)"
  >
    <template v-slot:append>
      <v-btn
        icon
        @click="addFunctionVar"
        title="Ajouter une fonction"
        size="x-small"
      >
        <v-icon icon="mdi-plus" color="green"></v-icon>
      </v-btn>
    </template>
    <v-list>
      <div
        v-for="(fn, index) in props.modelValue.FunctionVariations"
        :key="index"
      >
        <v-list-item>
          <v-row>
            <v-col>
              <BaseVariationTable
                :model-value="fn"
                @update:model-value="(v) => updateVar(index, v)"
                description=""
              ></BaseVariationTable>
            </v-col>
            <v-col cols="auto" align-self="center">
              <v-btn icon size="x-small" @click="deleteFunctionVar(index)">
                <v-icon icon="mdi-delete" color="red"></v-icon>
              </v-btn>
            </v-col>
          </v-row>
        </v-list-item>
        <v-divider></v-divider>
      </div>
    </v-list>
  </v-card>

  <!-- Sequences -->
  <v-card color="secondary" class="my-1" subtitle="Suites">
    <template v-slot:append>
      <v-btn
        icon
        @click="addSequenceExpr"
        title="Ajouter une suite"
        size="x-small"
      >
        <v-icon icon="mdi-plus" color="green"></v-icon>
      </v-btn>
    </template>
    <v-list>
      <div v-for="(fn, index) in props.modelValue.SequenceExprs" :key="index">
        <v-list-item>
          <v-row>
            <v-col>
              <v-row>
                <v-col cols="3" class="mb-4">
                  <interpolated-text
                    v-model="fn.Decoration.Label"
                    @update:model-value="emitUpdate"
                    force-latex
                    center
                    label="Légende"
                  ></interpolated-text>
                </v-col>
                <v-col md="3" align-self="center" class="text-center">
                  Variable : {{ variableToString(fn.Variable) }}
                </v-col>
                <v-col cols="6" align-self="center">
                  <v-text-field
                    variant="outlined"
                    density="compact"
                    v-model="fn.Function"
                    @update:model-value="emitUpdate"
                    label="Expression de la suite"
                    hide-details
                    :color="expressionColor"
                  ></v-text-field>
                </v-col>
                <v-col md="4" align-self="center">
                  <ExpressionField
                    v-model="fn.From"
                    @update:model-value="emitUpdate"
                    label="Xmin"
                  ></ExpressionField>
                </v-col>
                <v-col md="4">
                  <ExpressionField
                    v-model="fn.To"
                    @update:model-value="emitUpdate"
                    label="Xmax"
                  ></ExpressionField>
                </v-col>
                <v-col>
                  <BtnColorPicker
                    v-model="fn.Decoration.Color"
                    @update:model-value="emitUpdate"
                  ></BtnColorPicker>
                </v-col>
              </v-row>
            </v-col>

            <v-col cols="auto" align-self="center">
              <v-btn icon size="x-small" @click="deleteSequenceExpr(index)">
                <v-icon icon="mdi-delete" color="red"></v-icon>
              </v-btn>
            </v-col>
          </v-row>
        </v-list-item>
        <v-divider></v-divider>
      </div>
    </v-list>
  </v-card>

  <!-- areas -->
  <v-card color="secondary" class="my-1" subtitle="Surfaces colorées">
    <template v-slot:append>
      <v-btn
        icon
        @click="addArea"
        title="Ajouter une surface colorée entre deux courbes."
        size="x-small"
        :disabled="functionsNamesItems.length < 2"
      >
        <v-icon icon="mdi-plus" color="green"></v-icon>
      </v-btn>
    </template>
    <v-list>
      <div v-for="(area, index) in props.modelValue.Areas" :key="index">
        <v-list-item>
          <v-row class="mt-1">
            <v-col cols="2" align-self="center">
              <btn-color-picker
                v-model="area.Color"
                @update:model-value="emitUpdate"
              ></btn-color-picker>
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
                    @update:model-value="   (s) => {
                        area.Top = nameFromSelection(s as string); 
                        emitUpdate()
                    }                    "
                    :color="expressionColor"
                  ></v-combobox>
                </v-col>
                <v-col cols="5">
                  <ExpressionField
                    label="Xmin"
                    v-model="area.Left"
                    @update:model-value="emitUpdate"
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
                      (s) => {
                        area.Bottom = nameFromSelection(s as string);
                         emitUpdate();
                        }                  "
                    :color="expressionColor"
                  ></v-combobox>
                </v-col>
                <v-col cols="5">
                  <ExpressionField
                    label="Xmax"
                    v-model="area.Right"
                    @update:model-value="emitUpdate"
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

  <!-- additional points -->
  <v-card color="secondary" class="my-1" subtitle="Points sur les courbes">
    <template v-slot:append>
      <v-btn
        icon
        @click="addPoint"
        title="Ajouter un point appartenant à une des courbes"
        size="x-small"
      >
        <v-icon icon="mdi-plus" color="green"></v-icon>
      </v-btn>
    </template>
    <v-list>
      <div v-for="(point, index) in props.modelValue.Points" :key="index">
        <v-list-item>
          <v-row>
            <v-col cols="2" align-self="center">
              <btn-color-picker
                v-model="point.Color"
                @update:model-value="emitUpdate"
              ></btn-color-picker>
            </v-col>
            <v-col cols="4" align-self="center">
              <v-combobox
                density="compact"
                variant="outlined"
                hide-details
                hide-no-data
                label="Fonction"
                :items="functionsNamesItems"
                item
                :model-value="nameToSelection(point.Function)"
                @update:model-value="(s) => {
                    point.Function = nameFromSelection(s as string);
                     emitUpdate();
                }"
                :color="expressionColor"
              ></v-combobox>
            </v-col>
            <v-col cols="2" align-self="center">
              <ExpressionField
                label="X"
                v-model="point.X"
                @update:model-value="emitUpdate"
              ></ExpressionField>
            </v-col>
            <v-col cols="3" align-self="center">
              <InterpolatedText
                v-model="point.Legend"
                @update:model-value="emitUpdate"
                label="Légende"
              >
              </InterpolatedText>
            </v-col>
            <v-col md="1" align-self="center" class="pl-1 pr-0">
              <v-btn icon size="x-small" @click="deletePoint(index)">
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
import type {
  FunctionsGraphBlock,
  Variable,
  VariationTableBlock,
} from "@/controller/api_gen";
import {
  ExpressionColor,
  lastColorUsed,
  nRune,
  variableToString,
  xRune,
} from "@/controller/editor";
import BtnColorPicker from "../utils/BtnColorPicker.vue";
import ExpressionField from "../utils/ExpressionField.vue";
import InterpolatedText from "../utils/InterpolatedText.vue";
import BaseVariationTable from "./BaseVariationTable.vue";
import { computed } from "vue";

interface Props {
  modelValue: FunctionsGraphBlock;
  availableParameters: Variable[];
}
const props = defineProps<Props>();

const emit = defineEmits<{
  (event: "update:modelValue", value: FunctionsGraphBlock): void;
}>();

function emitUpdate() {
  emit("update:modelValue", props.modelValue);
}

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
  emitUpdate();
}

function deleteFunctionExpr(index: number) {
  props.modelValue.FunctionExprs?.splice(index, 1);
  emitUpdate();
}

function addSequenceExpr() {
  props.modelValue.SequenceExprs?.push({
    Function: "n + 2",
    Decoration: {
      Label: "u_n",
      Color: lastColorUsed.color,
    },
    Variable: { Name: nRune, Indice: "" },
    From: "-4",
    To: "4",
  });
  emitUpdate();
}

function deleteSequenceExpr(index: number) {
  props.modelValue.SequenceExprs?.splice(index, 1);
  emitUpdate();
}

function addFunctionVar() {
  props.modelValue.FunctionVariations?.push({
    Label: "C_f",
    Xs: ["-5", "0", "5"],
    Fxs: ["-3", "2", "-1"],
  });
  emitUpdate();
}

function deleteFunctionVar(index: number) {
  props.modelValue.FunctionVariations?.splice(index, 1);
  emitUpdate();
}

function updateVar(index: number, v: VariationTableBlock) {
  props.modelValue.FunctionVariations![index] = v;
  emitUpdate();
}

/** includes the special [abscisseAxis] value */
const functionsNamesItems = computed(() => {
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

const abscisseAxis = "<Axe Ox>";

function nameFromSelection(s: string) {
  return s == abscisseAxis ? "" : s;
}
function nameToSelection(s: string) {
  return s == "" ? abscisseAxis : s;
}

function addArea() {
  props.modelValue.Areas?.push({
    Color: lastColorUsed.color,
    Top: nameFromSelection(functionsNamesItems.value[1]),
    Bottom: nameFromSelection(functionsNamesItems.value[0]),
    Left: "0",
    Right: "2",
  });
  emitUpdate();
}

function deleteArea(index: number) {
  props.modelValue.Areas?.splice(index, 1);
  emitUpdate();
}

function addPoint() {
  props.modelValue.Points?.push({
    Color: lastColorUsed.color,
    Function: nameFromSelection(functionsNamesItems.value[0]),
    X: "0",
    Legend: `P_${1 + (props.modelValue.Points.length || 0)}`,
  });
  emitUpdate();
}

function deletePoint(index: number) {
  props.modelValue.Points?.splice(index, 1);
  emitUpdate();
}

const expressionColor = ExpressionColor;
</script>

<style scoped>
:deep(input) {
  width: 100%;
}
</style>
