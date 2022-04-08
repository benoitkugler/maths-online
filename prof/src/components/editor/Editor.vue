<template>
  <v-card class="ma-1">
    <v-row>
      <v-col md="auto">
        <v-card-title>Editeur de question</v-card-title>
        <v-card-subtitle
          >Editer et visualiser une question pour l'élève</v-card-subtitle
        >
      </v-col>
      <v-spacer></v-spacer>
      <v-col align-self="center" style="text-align: right">
        <v-menu offset-y close-on-content-click>
          <template v-slot:activator="{ isActive, props }">
            <v-btn
              icon
              title="Ajouter un contenu"
              v-on="{ isActive }"
              v-bind="props"
            >
              <v-icon icon="mdi-plus" color="green"></v-icon>
            </v-btn>
          </template>
          <block-bar @add="addBlock"></block-bar>
        </v-menu>

        <v-divider vertical></v-divider>
        <v-btn icon class="mx-2" @click="save" :disabled="!session_id">
          <v-icon icon="mdi-content-save"></v-icon>
        </v-btn>
      </v-col>
    </v-row>
    <v-row no-gutters>
      <v-col md="4" class="mx-2">
        <random-parameters
          :parameters="question.random_parameters"
          @add="addRandomParameter"
          @update="updateRandomParameter"
          @delete="deleteRandomParameter"
        ></random-parameters>
      </v-col>
      <v-col class="mr-2">
        <div
          style="height: 70vh; overflow-y: auto"
          @dragstart="onDragStart"
          @dragend="onDragEnd"
        >
          <drop-zone
            v-if="showDropZone"
            @drop="origin => swapBlocks(origin, 0)"
          ></drop-zone>
          <div v-for="(row, index) in rows" :key="index">
            <container
              @delete="removeBlock(index)"
              @swap="swapBlocks"
              :index="index"
              :nb-blocks="rows.length"
              :kind="row.Kind"
            >
              <component
                :model-value="row.Data"
                @update:model-value="(v: any) => updateBlock(index, v)"
                :is="component(row.Kind)"
              ></component>
            </container>
            <drop-zone
              v-if="showDropZone"
              @drop="origin => swapBlocks(origin, index + 1)"
            ></drop-zone>
          </div>
        </div>
      </v-col>
    </v-row>
  </v-card>
</template>

<script setup lang="ts">
import type { randomParameter } from "@/controller/api_gen";
import { controller } from "@/controller/controller";
import type {
  Block,
  FigureBlock,
  FormulaBlock,
  Question,
  TextBlock
} from "@/controller/exercice_gen";
import { BlockKind } from "@/controller/exercice_gen";
import { reactive, ref } from "@vue/reactivity";
import type { Component } from "@vue/runtime-core";
import { $ref } from "vue/macros";
import BlockBar from "./BlockBar.vue";
import Container from "./blocks/Container.vue";
import FigureVue from "./blocks/Figure.vue";
import FormulaVue from "./blocks/Formula.vue";
import TextVue from "./blocks/Text.vue";
import DropZone from "./DropZone.vue";
import RandomParameters from "./RandomParameters.vue";

const props = defineProps({
  session_id: { type: String, required: true }
});

const question: Question = reactive({
  title: "Nouvelle question",
  enonce: [],
  random_parameters: []
});

const rows = ref(<Block[]>[]); // TODO

function component(kind: BlockKind): Component {
  switch (kind) {
    case BlockKind.TextBlock:
      return TextVue;
    case BlockKind.FormulaBlock:
      return FormulaVue;
    case BlockKind.FigureBlock:
      return FigureVue;
    default:
      throw "Unexpected Kind";
  }
}

function newBlock(kind: BlockKind): Block {
  switch (kind) {
    case BlockKind.TextBlock:
      return {
        Kind: kind,
        Data: <TextBlock>{
          IsHint: false,
          Parts: ""
        }
      };
    case BlockKind.FormulaBlock:
      return {
        Kind: kind,
        Data: <FormulaBlock>{
          Parts: ""
        }
      };
    case BlockKind.FigureBlock:
      return {
        Kind: kind,
        Data: <FigureBlock>{
          ShowGrid: true,
          Bounds: {
            Width: 10,
            Height: 10,
            Origin: { X: 3, Y: 3 }
          },
          Drawings: {
            Lines: [],
            Points: [],
            Segments: []
          }
        }
      };
    default:
      throw "Unexpected Kind";
  }
}

function addBlock(kind: BlockKind) {
  rows.value.push(newBlock(kind));
}

function updateBlock(index: number, data: Block["Data"]) {
  rows.value[index].Data = data;
}

function removeBlock(index: number) {
  rows.value.splice(index, 1);
}

// TODO: fix text erasure when swapping
/** take the block at the index `origin` and insert it right before
the block at index `target` (which is between 0 and nbBlocks) 
 */
function swapBlocks(origin: number, target: number) {
  if (target == origin || target == origin + 1) {
    // nothing to do
    return;
  }

  if (origin < target) {
    const after = rows.value.slice(target);
    const before = rows.value.slice(0, target);
    const originRow = before.splice(origin, 1);
    before.push(...originRow);
    before.push(...after);
    rows.value = before;
  } else {
    const before = rows.value.slice(0, target);
    const originRow = rows.value.splice(origin, 1);
    const after = rows.value.slice(target);
    before.push(...originRow);
    before.push(...after);
    rows.value = before;
  }
}

function addRandomParameter() {
  question.random_parameters?.push({
    variable: "x".codePointAt(0)!,
    expression: ""
  });
}

function updateRandomParameter(index: number, param: randomParameter) {
  question.random_parameters![index] = param;
}

function deleteRandomParameter(index: number) {
  question.random_parameters!.splice(index, 1);
}

let showDropZone = $ref(false);

function onDragStart() {
  showDropZone = true;
}

function onDragEnd(ev: DragEvent) {
  showDropZone = false;
}

async function save() {
  question.enonce = rows.value.map(v => v);
  await controller.EditSaveAndPreview({
    SessionID: props.session_id || "",
    Question: question
  });
}
</script>
