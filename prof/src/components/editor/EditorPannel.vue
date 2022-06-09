<template>
  <v-snackbar
    :model-value="showErrorParameters"
    @update:model-value="errorParameters = null"
    color="warning"
  >
    <v-row v-if="errorParameters != null">
      <v-col>
        <v-row no-gutters>
          <v-col>
            Erreur dans la définition <b>{{ errorParameters.Origin }}</b> :
          </v-col>
        </v-row>
        <v-row>
          <v-col>
            <div>{{ errorParameters.Details }}</div>
          </v-col>
        </v-row>
      </v-col>
      <v-col cols="2" align-self="center" style="text-align: right">
        <v-btn icon size="x-small" flat @click="errorParameters = null">
          <v-icon icon="mdi-close" color="warning"></v-icon>
        </v-btn>
      </v-col>
    </v-row>
  </v-snackbar>

  <v-snackbar :model-value="showErrorEnnonce" color="warning">
    <v-row v-if="errorEnnonce != null">
      <v-col>
        <v-row no-gutters>
          <v-col> <b>Erreur dans la contenu de la question</b> </v-col>
        </v-row>
        <v-row>
          <v-col>
            <div>
              <i>{{ errorEnnonce.Error }}</i>
            </div>
          </v-col>
        </v-row>
      </v-col>
      <v-col cols="2" align-self="center" style="text-align: right">
        <v-btn icon size="x-small" flat @click="errorEnnonce = null">
          <v-icon icon="mdi-close" color="warning"></v-icon>
        </v-btn>
      </v-col>
    </v-row>
  </v-snackbar>

  <v-dialog v-model="showEditDescription">
    <description-pannel
      v-model="question.description"
      :readonly="isReadonly"
    ></description-pannel>
  </v-dialog>

  <v-card class="mt-3 px-2">
    <v-row no-gutters class="mb-2">
      <v-col cols="auto" align-self="center" class="pr-2">
        <v-btn
          size="small"
          icon
          title="Retour aux questions"
          @click="backToList"
        >
          <v-icon icon="mdi-arrow-left"></v-icon>
        </v-btn>
      </v-col>

      <v-col>
        <v-row no-gutters>
          <v-col>
            <v-text-field
              class="my-2 input-small"
              variant="outlined"
              density="compact"
              label="Nom de la question"
              v-model="question.page.title"
              :readonly="isReadonly"
              hide-details
            ></v-text-field
          ></v-col>
          <v-col cols="auto" align-self="center">
            <v-btn
              class="mx-1"
              size="x-small"
              icon
              @click="showEditDescription = true"
              title="Editer le commentaire"
            >
              <v-icon icon="mdi-message-reply-text" size="x-small"></v-icon>
            </v-btn>

            <v-btn
              class="mx-2"
              icon
              @click="save"
              :disabled="!session_id"
              :title="
                isReadonly ? 'Visualiser' : 'Enregistrer et prévisualiser'
              "
              size="small"
            >
              <v-icon
                :icon="isReadonly ? 'mdi-eye' : 'mdi-content-save'"
                size="small"
              ></v-icon>
            </v-btn>

            <v-menu offset-y close-on-content-click>
              <template v-slot:activator="{ isActive, props }">
                <v-btn
                  icon
                  title="Plus d'options"
                  v-on="{ isActive }"
                  v-bind="props"
                  size="x-small"
                >
                  <v-icon icon="mdi-dots-vertical"></v-icon>
                </v-btn>
              </template>
              <v-toolbar vertical density="compact" rounded>
                <v-btn
                  size="small"
                  icon
                  class="mr-2"
                  @click="download"
                  :disabled="!session_id"
                  title="Télécharger"
                >
                  <v-icon icon="mdi-download" size="small"></v-icon>
                </v-btn>
                <v-btn
                  class="mx-1"
                  size="small"
                  icon
                  @click="duplicate"
                  title="Dupliquer la question"
                >
                  <v-icon
                    icon="mdi-content-copy"
                    color="info"
                    size="small"
                  ></v-icon>
                </v-btn>
              </v-toolbar>
            </v-menu>
          </v-col>
        </v-row>

        <v-row no-gutters>
          <v-col class="pr-2">
            <tag-list-field
              label="Catégories"
              v-model="tags"
              :all-tags="props.allTags"
              @update:model-value="saveTags"
              :readonly="isReadonly"
            ></tag-list-field
          ></v-col>
          <v-col cols="auto">
            <v-menu offset-y close-on-content-click>
              <template v-slot:activator="{ isActive, props }">
                <v-btn
                  title="Ajouter un contenu"
                  v-on="{ isActive }"
                  v-bind="props"
                  size="small"
                >
                  <v-icon icon="mdi-plus" color="green"></v-icon>
                  Insérer
                </v-btn>
              </template>
              <block-bar @add="addBlock"></block-bar>
            </v-menu>
          </v-col>
        </v-row>
      </v-col>
    </v-row>

    <v-row no-gutters>
      <v-col md="4">
        <div style="height: 70vh; overflow-y: auto" class="py-2 px-2">
          <random-parameters
            :parameters="question.page.parameters.Variables"
            :is-loading="isCheckingParameters"
            :is-validated="!showErrorParameters"
            @add="addRandomParameter"
            @update="updateRandomParameter"
            @delete="deleteRandomParameter"
            @swap="swapRandomParameters"
            @done="checkParameters"
          ></random-parameters>
          <intrinsics
            :parameters="question.page.parameters.Intrinsics || []"
            :is-loading="isCheckingParameters"
            :is-validated="!showErrorParameters"
            @add="addIntrinsic"
            @update="updateIntrinsic"
            @delete="deleteIntrinsic"
            @done="checkParameters"
          ></intrinsics>
        </div>
      </v-col>
      <v-col class="pr-1">
        <div
          @drop="onDropJSON"
          @dragover="onDragoverJSON"
          class="d-flex ma-2"
          style="
            border: 1px solid blue;
            border-radius: 10px;
            height: 96%;
            justify-content: center;
            align-items: center;
          "
          v-if="rows.length == 0"
        >
          Importer une question en faisant glisser un fichier (.isyro.json) ...
        </div>

        <div
          v-else
          style="height: 70vh; overflow-y: auto"
          @dragstart="onDragStart"
          @dragend="onDragEnd"
        >
          <drop-zone
            v-if="showDropZone"
            @drop="(origin) => swapBlocks(origin, 0)"
          ></drop-zone>
          <div
            v-for="(row, index) in rows"
            :key="index"
            :ref="el => (blockWidgets[index] = el as Element)"
          >
            <container
              @delete="removeBlock(index)"
              :index="index"
              :nb-blocks="rows.length"
              :kind="row.Props.Kind"
              :hide-content="showDropZone"
              :has-error="errorEnnonce?.Block == index"
            >
              <component
                :model-value="row.Props.Data"
                @update:model-value="(v: any) => updateBlock(index, v)"
                :is="row.Component"
                :available-parameters="availableParameters"
              ></component>
            </container>
            <drop-zone
              v-if="showDropZone"
              @drop="(origin) => swapBlocks(origin, index + 1)"
            ></drop-zone>
          </div>
        </div>
      </v-col>
    </v-row>
  </v-card>
</template>

<script setup lang="ts">
import {
  Visibility,
  type errEnonce,
  type ErrParameters,
  type Origin,
} from "@/controller/api_gen";
import { controller } from "@/controller/controller";
import { newBlock, saveData, swapItems, xRune } from "@/controller/editor";
import type {
  Block,
  Question,
  RandomParameter,
  Variable,
} from "@/controller/exercice_gen";
import { BlockKind } from "@/controller/exercice_gen";
import { markRaw, ref } from "@vue/reactivity";
import type { Component } from "@vue/runtime-core";
import { computed, nextTick, watch } from "@vue/runtime-core";
import { $ref } from "vue/macros";
import BlockBar from "./BlockBar.vue";
import Container from "./blocks/Container.vue";
import FigureVue from "./blocks/Figure.vue";
import FigureAffineLineFieldVue from "./blocks/FigureAffineLineField.vue";
import FigurePointFieldVue from "./blocks/FigurePointField.vue";
import FigureVectorFieldVue from "./blocks/FigureVectorField.vue";
import FigureVectorPairFieldVue from "./blocks/FigureVectorPairField.vue";
import FormulaVue from "./blocks/Formula.vue";
import FormulaFieldVue from "./blocks/FormulaField.vue";
import FunctionGraphVue from "./blocks/FunctionGraph.vue";
import FunctionPointsFieldVue from "./blocks/FunctionPointsField.vue";
import FunctionVariationGraphVue from "./blocks/FunctionVariationGraph.vue";
import NumberFieldVue from "./blocks/NumberField.vue";
import OrderedListFieldVue from "./blocks/OrderedListField.vue";
import RadioFieldVue from "./blocks/RadioField.vue";
import SignTableVue from "./blocks/SignTable.vue";
import TableVue from "./blocks/Table.vue";
import TableFieldVue from "./blocks/TableField.vue";
import TextVue from "./blocks/Text.vue";
import TreeFieldVue from "./blocks/TreeField.vue";
import VariationTableVue from "./blocks/VariationTable.vue";
import VariationTableFieldVue from "./blocks/VariationTableField.vue";
import DescriptionPannel from "./DescriptionPannel.vue";
import DropZone from "./DropZone.vue";
import Intrinsics from "./Intrinsics.vue";
import RandomParameters from "./RandomParameters.vue";
import TagListField from "./TagListField.vue";

interface Props {
  session_id: string;
  question: Question;
  origin: Origin;
  tags: string[];
  allTags: string[]; // to provide auto completion
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "back"): void;
  (e: "duplicated", question: Question): void;
}>();

let question = $ref(props.question);
let tags = $ref(props.tags);

watch(props, () => {
  question = props.question;
  tags = props.tags;
});

const isReadonly = computed(
  () => props.origin.Visibility != Visibility.Personnal
);

const rows = computed(() => props.question.page.enonce?.map(dataToBlock) || []);

interface block {
  Props: Block;
  Component: Component;
}

function dataToBlock(data: Block): block {
  switch (data.Kind) {
    case BlockKind.TextBlock:
      return { Props: data, Component: markRaw(TextVue) };
    case BlockKind.FormulaBlock:
      return { Props: data, Component: markRaw(FormulaVue) };
    case BlockKind.FigureBlock:
      return { Props: data, Component: markRaw(FigureVue) };
    case BlockKind.FunctionGraphBlock:
      return { Props: data, Component: markRaw(FunctionGraphVue) };
    case BlockKind.FunctionVariationGraphBlock:
      return { Props: data, Component: markRaw(FunctionVariationGraphVue) };
    case BlockKind.VariationTableBlock:
      return { Props: data, Component: markRaw(VariationTableVue) };
    case BlockKind.SignTableBlock:
      return { Props: data, Component: markRaw(SignTableVue) };
    case BlockKind.TableBlock:
      return { Props: data, Component: markRaw(TableVue) };
    case BlockKind.NumberFieldBlock:
      return { Props: data, Component: markRaw(NumberFieldVue) };
    case BlockKind.ExpressionFieldBlock:
      return { Props: data, Component: markRaw(FormulaFieldVue) };
    case BlockKind.RadioFieldBlock:
      return { Props: data, Component: markRaw(RadioFieldVue) };
    case BlockKind.OrderedListFieldBlock:
      return { Props: data, Component: markRaw(OrderedListFieldVue) };
    case BlockKind.FigurePointFieldBlock:
      return { Props: data, Component: markRaw(FigurePointFieldVue) };
    case BlockKind.FigureVectorFieldBlock:
      return { Props: data, Component: markRaw(FigureVectorFieldVue) };
    case BlockKind.VariationTableFieldBlock:
      return { Props: data, Component: markRaw(VariationTableFieldVue) };
    case BlockKind.FunctionPointsFieldBlock:
      return { Props: data, Component: markRaw(FunctionPointsFieldVue) };
    case BlockKind.FigureVectorPairFieldBlock:
      return { Props: data, Component: markRaw(FigureVectorPairFieldVue) };
    case BlockKind.FigureAffineLineFieldBlock:
      return { Props: data, Component: markRaw(FigureAffineLineFieldVue) };
    case BlockKind.TreeFieldBlock:
      return { Props: data, Component: markRaw(TreeFieldVue) };
    case BlockKind.TableFieldBlock:
      return { Props: data, Component: markRaw(TableFieldVue) };
    default:
      throw "Unexpected Kind";
  }
}

const blockWidgets = ref<(Element | null)[]>([]);

function addBlock(kind: BlockKind) {
  question.page.enonce!.push(newBlock(kind));
  nextTick(() => {
    console.log(blockWidgets.value);

    const L = blockWidgets.value?.length;
    if (L) {
      blockWidgets.value[L - 1]?.scrollIntoView();
    }
  });
}

function updateBlock(index: number, data: Block["Data"]) {
  question.page.enonce![index].Data = data;
}

function removeBlock(index: number) {
  question.page.enonce!.splice(index, 1);
}

/** take the block at the index `origin` and insert it right before
the block at index `target` (which is between 0 and nbBlocks)
 */
function swapBlocks(origin: number, target: number) {
  question.page.enonce = swapItems(origin, target, question.page.enonce!);
}

function addRandomParameter() {
  const l = question.page.parameters.Variables || [];
  l.push({
    variable: { Name: xRune, Indice: "" },
    expression: "randint(1;10)",
  });
  question.page.parameters.Variables = l;
}

function updateRandomParameter(index: number, param: RandomParameter) {
  question.page.parameters.Variables![index] = param;
}

function deleteRandomParameter(index: number) {
  question.page.parameters.Variables!.splice(index, 1);

  checkParameters();
}

function swapRandomParameters(origin: number, target: number) {
  question.page.parameters.Variables = swapItems(
    origin,
    target,
    question.page.parameters.Variables!
  );
}

function addIntrinsic() {
  const l = question.page.parameters.Intrinsics || [];
  l.push("a,b,c = pythagorians()");
  question.page.parameters.Intrinsics = l;
}

function updateIntrinsic(index: number, param: string) {
  question.page.parameters.Intrinsics![index] = param;
}

function deleteIntrinsic(index: number) {
  question.page.parameters.Intrinsics!.splice(index, 1);

  checkParameters();
}

let showDropZone = $ref(false);

function onDragStart() {
  setTimeout(() => (showDropZone = true), 100); // workaround bug
}

function onDragEnd(ev: DragEvent) {
  showDropZone = false;
}

const errorEnnonce = ref<errEnonce | null>(null);
const showErrorEnnonce = computed(() => errorEnnonce.value != null);

async function save() {
  question.page.enonce = rows.value.map((v) => v.Props);
  const res = await controller.EditorSaveAndPreview({
    SessionID: props.session_id || "",
    Question: question,
  });
  if (res == undefined) {
    return;
  }

  if (res.IsValid) {
    errorEnnonce.value = null;
    errorParameters.value = null;
  } else {
    if (res.Error.ParametersInvalid) {
      errorEnnonce.value = null;
      errorParameters.value = res.Error.ErrParameters;
    } else {
      errorEnnonce.value = res.Error.ErrEnonce;
      errorParameters.value = null;

      blockWidgets.value[res.Error.ErrEnonce.Block]?.scrollIntoView();
    }
  }
}

function download() {
  saveData(question, "question.isyro.json");
}

async function onDropJSON(ev: DragEvent) {
  if (ev.dataTransfer?.files.length) {
    ev.preventDefault();
    const content = await ev.dataTransfer?.files[0].text();
    // keep the current ID
    const trueID = question.id;
    question = JSON.parse(content!);
    question.id = trueID;
  }
}

function onDragoverJSON(ev: DragEvent) {
  if (ev.dataTransfer?.files.length || ev.dataTransfer?.items.length) {
    ev.preventDefault();
  }
}

const errorParameters = ref<ErrParameters | null>(null);
const showErrorParameters = computed(() => errorParameters.value != null);
const availableParameters = ref<Variable[]>([]);
let isCheckingParameters = $ref(false);

async function checkParameters() {
  isCheckingParameters = true;
  const out = await controller.EditorCheckParameters({
    SessionID: props.session_id || "",
    Parameters: question.page.parameters,
  });
  isCheckingParameters = false;
  if (out === undefined) return;

  errorParameters.value =
    out.ErrDefinition.Origin == "" ? null : out.ErrDefinition;
  availableParameters.value = out.Variables || [];
}

async function saveTags() {
  await controller.EditorUpdateTags({ IdQuestion: question.id, Tags: tags });
}

async function duplicate() {
  const newQuestion = await controller.EditorDuplicateQuestion({
    id: question.id,
  });
  if (newQuestion == undefined) {
    return;
  }
  emit("duplicated", newQuestion);
}

function backToList() {
  emit("back");
}

let showEditDescription = $ref(false);
</script>

<style scoped>
.input-small:deep(input) {
  font-size: 14px;
  width: 100%;
}
</style>
