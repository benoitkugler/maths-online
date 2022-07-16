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
    <v-list v-if="props.parameters?.length" @dragend="showDropZone = false">
      <drop-zone
        v-if="showDropZone"
        @drop="(origin) => swapIntrinsics(origin, 0)"
      ></drop-zone>
      <div v-for="(param, index) in props.parameters" :key="index">
        <v-list-item class="pr-0 pl-1">
          <v-row no-gutters>
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
            <v-col cols="9" class="pl-1">
              <v-text-field
                class="small-input"
                hide-details
                variant="underlined"
                density="compact"
                :model-value="param"
                @update:model-value="(v) => autocomplete(index, v)"
                @blur="emit('done')"
              ></v-text-field>
            </v-col>
            <v-col cols="2" align-self="center" style="text-align: right">
              <v-btn
                icon
                size="x-small"
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
          @drop="(origin) => swapIntrinsics(origin, index + 1)"
        ></drop-zone>
      </div>
    </v-list>
  </parameters-container>
</template>

<script setup lang="ts">
import { onDragListItemStart, swapItems } from "@/controller/editor";
import { ref } from "@vue/reactivity";
import { $ref } from "vue/macros";
import DropZone from "./DropZone.vue";
import IntrinsicsParametersHelp from "./parameters/IntrinsicsParametersHelp.vue";
import ParametersContainer from "./parameters/ParametersContainer.vue";

interface Props {
  parameters: string[];
  isLoading: boolean;
  isValidated: boolean;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "update", list: string[], shouldCheck: boolean): void;
  (e: "done"): void;
}>();

const showHelp = ref(false);

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
  const l = props.parameters || [];
  l.push("a,b,c = pythagorians()");
  emit("update", l, false);
}

function update(index: number, param: string) {
  props.parameters![index] = param;
  emit("update", props.parameters, false);
}

function remove(index: number) {
  props.parameters!.splice(index, 1);
  emit("update", props.parameters, true);
}

function swapIntrinsics(origin: number, target: number) {
  const l = swapItems(origin, target, props.parameters!);
  emit("update", l, false);
}
</script>

<style scoped>
.small-input:deep(input) {
  font-size: 14px;
  width: 100%;
}
</style>
