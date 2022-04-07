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
      <v-col md="4" class="mr-2">
        <random-parameters
          :parameters="question.random_parameters"
          @add="addRandomParameter"
          @update="updateRandomParameter"
          @delete="deleteRandomParameter"
        ></random-parameters>
      </v-col>
      <v-col>
        <v-card style="height: 70vh; overflow-y: auto">
          <component
            v-for="(row, index) in rows"
            :Data="row.props.Data"
            :Pos="blockPosition(index)"
            :is="row.component"
            @delete="removeBlock(index)"
            @swap="swapBlocks"
          ></component>
        </v-card>
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
import { markRaw, reactive } from "@vue/reactivity";
import type { Component } from "@vue/runtime-core";
import BlockBar from "./BlockBar.vue";
import type { ContainerProps } from "./blocks/Container.vue";
import FigureVue from "./blocks/Figure.vue";
import FormulaVue from "./blocks/Formula.vue";
import TextVue from "./blocks/Text.vue";
import RandomParameters from "./RandomParameters.vue";

const props = defineProps({
  session_id: { type: String, required: true }
});

const question: Question = reactive({
  title: "Nouvelle question",
  enonce: [],
  random_parameters: []
});

type block = {
  component: Component;
  props: Block;
};

let rows = reactive(<block[]>[]); // TODO

function blockPosition(index: number): ContainerProps {
  return { index: index, nbBlocks: rows.length };
}

function newBlock(kind: BlockKind): block {
  switch (kind) {
    case BlockKind.TextBlock:
      return {
        component: markRaw(TextVue),
        props: {
          Kind: kind,
          Data: <TextBlock>{
            IsHint: false,
            Parts: ""
          }
        }
      };
    case BlockKind.FormulaBlock:
      return {
        component: markRaw(FormulaVue),
        props: {
          Kind: kind,
          Data: <FormulaBlock>{
            Parts: ""
          }
        }
      };
    case BlockKind.FigureBlock:
      return {
        component: markRaw(FigureVue),
        props: {
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
        }
      };
    default:
      throw "Unexpected Kind";
  }
}

function addBlock(kind: BlockKind) {
  rows.push(newBlock(kind));
}

function removeBlock(index: number) {
  rows.splice(index, 1);
}

// TODO: fix text erasure when swapping
// TODO: insert and shift instead of swapping
function swapBlocks(origin: number, target: number) {
  if (origin == target) {
    return;
  }
  var b = rows[origin];
  rows[origin] = rows[target];
  rows[target] = b;
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

async function save() {
  question.enonce = rows.map(v => v.props);
  await controller.EditSaveAndPreview({
    SessionID: props.session_id || "",
    Question: question
  });
}
</script>
