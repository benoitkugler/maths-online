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
              <i v-html="errorEnnonce.Error"></i>
            </div>
          </v-col>
        </v-row>
      </v-col>
      <v-col
        v-if="errVars.length > 0"
        cols="3"
        align-self="center"
        class="px-1"
      >
        <v-btn variant="outlined" @click="showErrVarsDetails = true">
          Détails
        </v-btn>
      </v-col>
      <v-col
        cols="2"
        align-self="center"
        style="text-align: right"
        class="px-1"
      >
        <v-btn icon size="x-small" @click="errorEnnonce = null">
          <v-icon icon="mdi-close" color="warning"></v-icon>
        </v-btn>
      </v-col>
    </v-row>
  </v-snackbar>

  <v-dialog v-model="showErrVarsDetails">
    <v-card subtitle="Valeurs des paramètres aléatoires">
      <v-card-text>
        L'erreur est rencontrée pour les valeurs suivantes :
        <v-list>
          <v-list-item v-for="(entry, index) in errVars" :key="index">
            <v-row no-gutters>
              <v-col>
                {{ entry[0] }}
              </v-col>
              <v-col class="text-grey">
                {{ entry[1] }}
              </v-col>
            </v-row>
          </v-list-item>
        </v-list>
      </v-card-text>
    </v-card>
  </v-dialog>

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
              <v-list>
                <v-list-item>
                  <v-btn
                    size="small"
                    @click="showEditDescription = true"
                    title="Editer le commentaire"
                  >
                    <v-icon
                      class="mr-2"
                      icon="mdi-message-reply-text"
                      size="small"
                    ></v-icon>
                    Commentaire
                  </v-btn>
                </v-list-item>
                <v-list-item>
                  <v-btn
                    size="small"
                    @click="download"
                    :disabled="!session_id"
                    title="Télécharger la question au format .json"
                  >
                    <v-icon
                      class="mr-2"
                      icon="mdi-download"
                      size="small"
                    ></v-icon>
                    Télécharger
                  </v-btn>
                </v-list-item>
                <v-list-item>
                  <v-btn
                    size="small"
                    @click="duplicate"
                    title="Dupliquer la question"
                  >
                    <v-icon
                      class="mr-2"
                      icon="mdi-content-copy"
                      color="info"
                      size="small"
                    ></v-icon>
                    Dupliquer
                  </v-btn>
                </v-list-item>
              </v-list>
            </v-menu>
          </v-col>
        </v-row>

        <v-row no-gutters>
          <v-col class="pr-2" align-self="center">
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
                  title="Ajouter un bloc de contenu (énoncé ou champ de réponse)"
                  v-on="{ isActive }"
                  v-bind="props"
                  size="small"
                >
                  <v-icon icon="mdi-plus" color="green"></v-icon>
                  Insérer du contenu
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
        <QuestionContent
          :model-value="question.page.enonce || []"
          @update:model-value="(v) => (question.page.enonce = v)"
          @importQuestion="onImportQuestion"
          :available-parameters="availableParameters"
          :errorBlockIndex="errorEnnonce?.Block"
          ref="questionContent"
        >
        </QuestionContent>
        <!-- <div
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
            <BlockContainer
              @delete="removeBlock(index)"
              :index="index"
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
            </BlockContainer>
            <drop-zone
              v-if="showDropZone"
              @drop="(origin) => swapBlocks(origin, index + 1)"
            ></drop-zone>
          </div>
        </div> -->
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
import { saveData, swapItems, xRune } from "@/controller/editor";
import type {
  BlockKind,
  Question,
  RandomParameter,
  Variable,
} from "@/controller/exercice_gen";
import { ref } from "@vue/reactivity";
import { computed, watch } from "@vue/runtime-core";
import { $ref } from "vue/macros";
import BlockBar from "./BlockBar.vue";
import DescriptionPannel from "./DescriptionPannel.vue";
import Intrinsics from "./Intrinsics.vue";
import QuestionContent from "./QuestionContent.vue";
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

let questionContent = $ref<InstanceType<typeof QuestionContent> | null>(null);

function addBlock(kind: BlockKind) {
  if (questionContent == null) {
    return;
  }
  questionContent.addBlock(kind);
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

let errorEnnonce = $ref<errEnonce | null>(null);
const showErrorEnnonce = computed(() => errorEnnonce != null);
const errVars = computed(() => {
  const out = Object.entries(errorEnnonce?.Vars || {});
  out.sort((a, b) => a[0].localeCompare(b[0]));
  return out;
});
let showErrVarsDetails = $ref(false);

async function save() {
  const res = await controller.EditorSaveQuestionAndPreview({
    SessionID: props.session_id || "",
    Question: question,
  });
  if (res == undefined) {
    return;
  }

  if (res.IsValid) {
    errorEnnonce = null;
    errorParameters = null;
  } else {
    if (res.Error.ParametersInvalid) {
      errorEnnonce = null;
      errorParameters = res.Error.ErrParameters;
    } else {
      errorEnnonce = res.Error.ErrEnonce;
      errorParameters = null;
    }
  }
}

function download() {
  saveData(question, "question.isyro.json");
}

async function onImportQuestion(imported: Question) {
  // keep the current ID
  imported.id = question.id;
  question = imported;
}

let errorParameters = $ref<ErrParameters | null>(null);
const showErrorParameters = computed(() => errorParameters != null);
const availableParameters = ref<Variable[]>([]);
let isCheckingParameters = $ref(false);

async function checkParameters() {
  isCheckingParameters = true;
  const out = await controller.EditorCheckQuestionParameters({
    SessionID: props.session_id || "",
    Parameters: question.page.parameters,
  });
  isCheckingParameters = false;
  if (out === undefined) return;

  // hide previous error
  errorEnnonce = null;

  errorParameters = out.ErrDefinition.Origin == "" ? null : out.ErrDefinition;
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
