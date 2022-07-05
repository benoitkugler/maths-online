<template>
  <random-parameters-help v-model="showHelp"></random-parameters-help>

  <parameters-container
    title="Paramètres aléatoires"
    add-title="Ajouter un paramètre aléatoire"
    :is-validated="props.isValidated"
    :is-loading="props.isLoading"
    @add="add"
    @show-help="showHelp = true"
    class="mb-2"
  >
    <v-list v-if="parameters.length" @dragend="onDragend">
      <drop-zone
        v-if="showDropZone"
        @drop="(origin) => swap(origin, 0)"
      ></drop-zone>
      <div v-for="(param, index) in parameters" :key="index">
        <v-list-item class="pr-0 pl-1">
          <v-row no-gutters :class="param.isShared ? 'text-pink' : ''">
            <v-col cols="1" align-self="center">
              <v-icon
                size="large"
                class="pr-2"
                style="cursor: grab"
                @dragstart="(e) => onItemDragStart(e, index)"
                draggable="true"
                color="green-lighten-3"
                icon="mdi-drag-vertical"
              ></v-icon>
            </v-col>
            <v-col cols="3">
              <variable-field
                v-model="param.parameter.variable"
                @update:model-value="update(index, param.parameter)"
                @blur="emit('done')"
              >
              </variable-field>
            </v-col>
            <v-col cols="6">
              <v-text-field
                class="ml-2 small-input"
                variant="underlined"
                density="compact"
                hide-details
                :model-value="param.parameter.expression"
                @update:model-value="(s) => onExpressionChange(s, index)"
                @blur="emit('done')"
                :color="expressionColor"
              ></v-text-field>
            </v-col>
            <v-col cols="2">
              <v-btn icon size="small" flat @click="remove(index)">
                <v-icon icon="mdi-delete" color="red"></v-icon>
              </v-btn>
            </v-col>
          </v-row>
        </v-list-item>
        <drop-zone
          v-if="showDropZone"
          @drop="(origin) => swap(origin, index + 1)"
        ></drop-zone>
      </div>
    </v-list>
  </parameters-container>
</template>

<script setup lang="ts">
import {
  ExpressionColor,
  onDragListItemStart,
  swapItems,
  xRune,
} from "@/controller/editor";
import type {
  RandomParameter,
  RandomParameters,
} from "@/controller/exercice_gen";
import { $computed, $ref } from "vue/macros";
import DropZone from "./DropZone.vue";
import ParametersContainer from "./parameters/ParametersContainer.vue";
import RandomParametersHelp from "./parameters/RandomParametersHelp.vue";
import VariableField from "./utils/VariableField.vue";

interface Props {
  sharedParameters: RandomParameters;
  questionParameters: RandomParameters;
  isLoading: boolean;
  isValidated: boolean;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (
    e: "update",
    sharedParameters: RandomParameter[],
    questionParameters: RandomParameter[],
    shouldCheck: boolean
  ): void;
  (e: "done"): void;
}>();

const expressionColor = ExpressionColor;

let showHelp = $ref(false);

interface withShared {
  parameter: RandomParameter;
  isShared: boolean;
}

let parameters = $computed(() =>
  (props.sharedParameters || [])
    .map((p) => ({ parameter: p, isShared: true }))
    .concat(
      (props.questionParameters || []).map((p) => ({
        parameter: p,
        isShared: false,
      }))
    )
);

function emitSeparateList(merged: withShared[], shouldCheck: boolean) {
  const shared = merged
    .filter((item) => item.isShared)
    .map((item) => item.parameter);
  const question = merged
    .filter((item) => !item.isShared)
    .map((item) => item.parameter);
  emit("update", shared, question, shouldCheck);
}

function onExpressionChange(s: string, index: number) {
  const l = parameters.map((v) => v);
  const param = l[index];
  param.parameter.expression = s;
  emitSeparateList(l, false);
}

let showDropZone = $ref(false);
function onItemDragStart(payload: DragEvent, index: number) {
  onDragListItemStart(payload, index);
  setTimeout(() => (showDropZone = true), 100); // workaround bug
}

function onDragend() {
  showDropZone = false;
}

// add as shared parameter by default
function add() {
  const l = parameters.map((v) => v);
  l.push({
    parameter: {
      variable: { Name: xRune, Indice: "" },
      expression: "randint(1;10)",
    },
    isShared: true,
  });
  emitSeparateList(l, false);
}

function update(index: number, param: RandomParameter) {
  const l = parameters.map((v) => v);
  l[index].parameter = param;
  emitSeparateList(l, false);
}

function remove(index: number) {
  const l = parameters.map((v) => v);
  l.splice(index, 1);
  emitSeparateList(l, true);
}

function swap(origin: number, target: number) {
  const l = swapItems(origin, target, parameters);
  emitSeparateList(l, false);
}
</script>

<style scoped>
.small-input:deep(input) {
  font-size: 14px;
  width: 100%;
}
</style>
