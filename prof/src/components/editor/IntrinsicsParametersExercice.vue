<template>
  <intrinsics-parameters-help v-model="showHelp"> </intrinsics-parameters-help>

  <parameters-container
    title="Fonctions spéciales"
    add-title="Ajouter une fonction spéciale"
    :is-validated="props.isValidated"
    :is-loading="props.isLoading"
    @add="add"
    @show-help="showHelp = true"
  >
    <v-list v-if="parameters.length" @dragend="showDropZone = false">
      <drop-zone
        v-if="showDropZone"
        @drop="(origin) => swap(origin, 0)"
      ></drop-zone>
      <div v-for="(param, index) in parameters" :key="index">
        <v-list-item class="pr-0 pl-1">
          <v-row no-gutters :class="param.isShared ? 'text-primary' : ''">
            <v-col cols="1" align-self="center">
              <drag-icon
                color="green-lighten-3"
                @start="(e) => onItemDragStart(e, index)"
              ></drag-icon>
            </v-col>
            <v-col class="pl-1">
              <v-text-field
                class="small-input"
                hide-details
                variant="underlined"
                density="compact"
                :model-value="param.parameter"
                @update:model-value="(v) => autocomplete(index, v)"
                @blur="emit('done')"
              ></v-text-field>
            </v-col>
            <v-col cols="auto" align-self="center">
              <v-btn
                icon
                size="small"
                flat
                @click="toogleShared(index)"
                title="Partager les variables entre les questions"
                class="pl-1"
              >
                <v-icon
                  icon="mdi-format-list-numbered"
                  :color="param.isShared ? 'primary' : 'black'"
                  size="small"
                ></v-icon>
              </v-btn>
              <v-btn
                icon
                size="small"
                flat
                @click="remove(index)"
                title="Supprimer cette fonction spéciale"
              >
                <v-icon icon="mdi-delete" color="red" size="small"></v-icon>
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
import DropZone from "@/components/DropZone.vue";
import { onDragListItemStart, swapItems } from "@/controller/utils";
import { ref } from "@vue/reactivity";
import { $computed, $ref } from "vue/macros";
import DragIcon from "../DragIcon.vue";
import IntrinsicsParametersHelp from "./parameters/IntrinsicsParametersHelp.vue";
import ParametersContainer from "./parameters/ParametersContainer.vue";

interface Props {
  sharedParameters: string[];
  questionParameters: string[];
  isLoading: boolean;
  isValidated: boolean;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (
    e: "update",
    sharedParameters: string[],
    questionParameters: string[],
    shouldCheck: boolean
  ): void;
  (e: "done"): void;
}>();

const showHelp = ref(false);

interface withShared {
  parameter: string;
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

let showDropZone = $ref(false);
function onItemDragStart(payload: DragEvent, index: number) {
  onDragListItemStart(payload, index);
  setTimeout(() => (showDropZone = true), 100); // workaround bug
}

// to keep sync with the server
const intrinsics = ["pythagorians", "projection"];

function autocomplete(index: number, text: string) {
  for (const it of intrinsics) {
    if (text.endsWith(it.substring(0, 4))) {
      text += it.substring(4) + "()";
      break;
    }
  }

  update(index, text);
}

function add() {
  const l = parameters.map((v) => v);
  l.push({
    parameter: "a,b,c = pythagorians()",
    isShared: true,
  });
  emitSeparateList(l, false);
}

function update(index: number, param: string) {
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

function toogleShared(index: number) {
  const l = parameters.map((v) => v);
  const p = l[index];
  p.isShared = !p.isShared;
  emitSeparateList(l, true);
}
</script>

<style scoped>
.small-input:deep(input) {
  font-size: 14px;
  width: 100%;
}
</style>
