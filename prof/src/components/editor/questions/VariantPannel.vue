<template>
  <SnackErrorParameters
    :error="errorParameters"
    @close="errorParameters = null"
  >
  </SnackErrorParameters>

  <SnackErrorEnonce
    :error="errorEnnonce"
    @close="errorEnnonce = null"
  ></SnackErrorEnonce>

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
                    class="my-1"
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
                    class="my-1"
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
                    class="my-1"
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
              :all-tags="props.allTags"
              :model-value="tags"
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
        <div style="height: 68vh; overflow-y: auto" class="py-2 px-2">
          <RandomParametersQuestion
            :parameters="question.page.parameters.Variables"
            :is-loading="isCheckingParameters"
            :is-validated="!showErrorParameters"
            @update="updateRandomParameters"
            @done="checkParameters"
          ></RandomParametersQuestion>
          <IntrinsicsParametersQuestion
            :parameters="question.page.parameters.Intrinsics || []"
            :is-loading="isCheckingParameters"
            :is-validated="!showErrorParameters"
            @update="updateIntrinsics"
            @done="checkParameters"
          ></IntrinsicsParametersQuestion>
        </div>
      </v-col>
      <v-col class="pr-1">
        <QuestionContent
          :model-value="question.page.enonce || []"
          @update:model-value="onUpdateEnonce"
          @importQuestion="onImportQuestion"
          :available-parameters="availableParameters"
          :errorBlockIndex="errorEnnonce?.Block"
          ref="questionContent"
        >
        </QuestionContent>
      </v-col>
    </v-row>
  </v-card>
</template>

<script setup lang="ts">
import type {
  Block,
  BlockKind,
  Question,
  RandomParameter,
  Variable,
} from "@/controller/api_gen";
import {
  Visibility,
  type errEnonce,
  type ErrParameters,
} from "@/controller/api_gen";
import { controller } from "@/controller/controller";
import { saveData } from "@/controller/editor";
import { History } from "@/controller/editor_history";
import { copy } from "@/controller/utils";
import { ref } from "@vue/reactivity";
import { computed, onMounted, onUnmounted } from "@vue/runtime-core";
import { $ref } from "vue/macros";
import BlockBar from "./BlockBar.vue";
import DescriptionPannel from "./DescriptionPannel.vue";
import IntrinsicsParametersQuestion from "./IntrinsicsParametersQuestion.vue";
import SnackErrorParameters from "./parameters/SnackErrorParameters.vue";
import QuestionContent from "./QuestionContent.vue";
import RandomParametersQuestion from "./RandomParametersQuestion.vue";
import SnackErrorEnonce from "./SnackErrorEnonce.vue";
import TagListField from "./TagListField.vue";

interface Props {
  session_id: string;
  question: Question;
  allTags: string[]; // to provide auto completion
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "back"): void;
  (e: "duplicated", question: Question): void;
}>();

let question = $ref(copy(props.question));

onMounted(() => {
  history.addListener();
});

onUnmounted(() => {
  history.clearListener();
});

const isReadonly = computed(
  () => props.origin.Visibility != Visibility.Personnal
);

interface historyEntry {
  question: Question;
  tags: string[];
}

let history = new History(
  { question, tags }, // start with initial question in history
  controller.showMessage!,
  restoreHistory
);

function restoreHistory(snapshot: historyEntry) {
  // do we need to save the tags ?
  const tagsEqual = JSON.stringify(tags) == JSON.stringify(snapshot.tags);
  question = snapshot.question;
  tags = snapshot.tags;
  if (!tagsEqual) {
    controller.EditorUpdateTags({ IdQuestion: question.id, Tags: tags });
  }
}

let questionContent = $ref<InstanceType<typeof QuestionContent> | null>(null);
function addBlock(kind: BlockKind) {
  if (questionContent == null) {
    return;
  }
  questionContent.addBlock(kind);
}

function onUpdateEnonce(v: Block[]) {
  question.page.enonce = v;
  history.add({ question, tags });
}

function updateRandomParameters(l: RandomParameter[], shouldCheck: boolean) {
  question.page.parameters.Variables = l;
  if (shouldCheck) {
    checkParameters();
  }

  history.add({ question, tags });
}

function updateIntrinsics(l: string[], shouldCheck: boolean) {
  question.page.parameters.Intrinsics = l;
  if (shouldCheck) {
    checkParameters();
  }

  history.add({ question, tags });
}

let errorEnnonce = $ref<errEnonce | null>(null);

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

  history.add({ question, tags });
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

async function saveTags(newTags: string[]) {
  await controller.EditorUpdateTags({ IdQuestion: question.id, Tags: newTags });
  tags = newTags;
  history.add({ question, tags });
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
  emit("back", tags);
}

let showEditDescription = $ref(false);
</script>

<style scoped>
.input-small:deep(input) {
  font-size: 14px;
  width: 100%;
}
</style>
