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
              :kind="row.Props.Kind"
            >
              <component
                :model-value="row.Props.Data"
                @update:model-value="(v: any) => updateBlock(index, v)"
                :is="row.Component"
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
import type { TypedBlock } from "@/controller/editor";
import type { Block, Question } from "@/controller/exercice_gen";
import { BlockKind, SignSymbol, TextKind } from "@/controller/exercice_gen";
import { markRaw, reactive, ref } from "@vue/reactivity";
import type { Component } from "@vue/runtime-core";
import { $ref } from "vue/macros";
import BlockBar from "./BlockBar.vue";
import Container from "./blocks/Container.vue";
import FigureVue from "./blocks/Figure.vue";
import FormulaVue from "./blocks/Formula.vue";
import FunctionGraphVue from "./blocks/FunctionGraph.vue";
import FunctionVariationGraphVue from "./blocks/FunctionVariationGraph.vue";
import SignTableVue from "./blocks/SignTable.vue";
import TableVue from "./blocks/Table.vue";
import TextVue from "./blocks/Text.vue";
import VariationTableVue from "./blocks/VariationTable.vue";
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

const rows = ref(<block[]>[]); // TODO

interface block {
  Props: Block;
  Component: Component;
}

function newBlock(kind: BlockKind): block {
  switch (kind) {
    case BlockKind.TextBlock: {
      const out: TypedBlock<typeof kind> = {
        Kind: kind,
        Data: {
          IsHint: false,
          Parts: ""
        }
      };
      return { Props: out, Component: markRaw(TextVue) };
    }
    case BlockKind.FormulaBlock: {
      const out: TypedBlock<typeof kind> = {
        Kind: kind,
        Data: {
          Parts: ""
        }
      };
      return { Props: out, Component: markRaw(FormulaVue) };
    }
    case BlockKind.FigureBlock: {
      const out: TypedBlock<typeof kind> = {
        Kind: kind,
        Data: {
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
      return { Props: out, Component: markRaw(FigureVue) };
    }
    case BlockKind.FunctionGraphBlock: {
      const out: TypedBlock<typeof kind> = {
        Kind: kind,
        Data: {
          Function: "",
          Label: "f",
          Variable: "x".codePointAt(0)!,
          Range: [-5, 5]
        }
      };
      return { Props: out, Component: markRaw(FunctionGraphVue) };
    }
    case BlockKind.FunctionVariationGraphBlock: {
      const out: TypedBlock<typeof kind> = {
        Kind: kind,
        Data: {
          Xs: ["-5", "0", "5"],
          Fxs: ["-3", "2", "-1"]
        }
      };
      return { Props: out, Component: markRaw(FunctionVariationGraphVue) };
    }
    case BlockKind.VariationTableBlock: {
      const out: TypedBlock<typeof kind> = {
        Kind: kind,
        Data: {
          Xs: ["-5", "0", "5"],
          Fxs: ["-3", "2", "-1"]
        }
      };
      return { Props: out, Component: markRaw(VariationTableVue) };
    }
    case BlockKind.SignTableBlock: {
      const out: TypedBlock<typeof kind> = {
        Kind: kind,
        Data: {
          FxSymbols: [
            SignSymbol.Nothing,
            SignSymbol.Zero,
            SignSymbol.ForbiddenValue,
            SignSymbol.Nothing
          ],
          Xs: ["\\infty", "3", "5", "+\\infty"],
          Signs: [true, false, true]
        }
      };
      return { Props: out, Component: markRaw(SignTableVue) };
    }
    case BlockKind.TableBlock: {
      const out: TypedBlock<typeof kind> = {
        Kind: kind,
        Data: {
          VerticalHeaders: [
            { Kind: TextKind.Text, Content: "Ligne 1" },
            { Kind: TextKind.Text, Content: "Ligne 2" }
          ],
          HorizontalHeaders: [
            { Kind: TextKind.Text, Content: "Colonne 1" },
            { Kind: TextKind.Text, Content: "Colonne 2" }
          ],
          Values: [
            [
              { Kind: TextKind.Text, Content: "Case" },
              { Kind: TextKind.StaticMath, Content: "\\frac{1}{2}" }
            ],
            [
              { Kind: TextKind.Expression, Content: "2x + 3" },
              { Kind: TextKind.StaticMath, Content: "18" }
            ]
          ]
        }
      };
      return { Props: out, Component: markRaw(TableVue) };
    }
    default:
      throw "Unexpected Kind";
  }
}

function addBlock(kind: BlockKind) {
  rows.value.push(newBlock(kind));
}

function updateBlock(index: number, data: Block["Data"]) {
  rows.value[index].Props.Data = data;
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
    expression: "randint(1;10)"
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
  question.enonce = rows.value.map(v => v.Props);
  await controller.EditSaveAndPreview({
    SessionID: props.session_id || "",
    Question: question
  });
}
</script>
