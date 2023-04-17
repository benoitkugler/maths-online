<template>
  <SnackErrorParameters
    :error="errorParameters"
    @close="errorParameters = null"
  >
  </SnackErrorParameters>

  <SnackErrorEnonce
    :error="errorContent"
    :is-correction="errorIsCorrection"
    @close="errorContent = null"
  ></SnackErrorEnonce>

  <v-card class="mt-1 px-2 pr-1">
    <v-row no-gutters>
      <v-col md="5" v-if="!hasEditorSimplified"> </v-col>

      <v-col cols="auto" align-self="center">
        <v-btn-toggle
          variant="tonal"
          class="py-1"
          density="compact"
          :model-value="modeEnonce ? 0 : 1"
          @update:model-value="(i:number) => (modeEnonce = i == 0)"
        >
          <v-btn size="small">énoncé</v-btn>
          <v-btn size="small">Correction</v-btn>
        </v-btn-toggle>
      </v-col>

      <v-spacer></v-spacer>

      <v-col cols="auto" align-self="center">
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
          <block-bar
            @add="addBlock"
            :simplified="hasEditorSimplified"
            :hide-answer-fields="!modeEnonce"
          ></block-bar>
        </v-menu>
      </v-col>

      <v-col cols="auto" align-self="center">
        <v-btn
          class="my-1 mx-2"
          icon
          @click="save"
          :title="
            props.readonly ? 'Visualiser' : 'Enregistrer et prévisualiser'
          "
          size="small"
        >
          <v-icon
            :icon="props.readonly ? 'mdi-eye' : 'mdi-content-save'"
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
                @click="download"
                title="Télécharger la question au format .json"
              >
                <v-icon class="mr-2" icon="mdi-download" size="small"></v-icon>
                Télécharger
              </v-btn>
            </v-list-item>
            <v-list-item>
              <v-btn
                class="my-1"
                size="small"
                @click="exportLatex"
                title="Exporter la question au format LaTeX (.tex)"
              >
                <v-icon
                  class="mr-2"
                  icon="mdi-file-export"
                  size="small"
                ></v-icon>
                Exporter en LaTeX
              </v-btn>
            </v-list-item>
          </v-list>
        </v-menu>
      </v-col>
    </v-row>

    <v-row no-gutters>
      <v-col md="5" v-if="!hasEditorSimplified">
        <ParametersEditor
          :model-value="question.Parameters"
          @update:model-value="checkParameters"
          :is-loading="isCheckingParameters"
          :is-validated="!showErrorParameters"
          :show-switch="null"
        ></ParametersEditor>
      </v-col>
      <v-col class="pr-1">
        <QuestionContent
          v-if="modeEnonce"
          :model-value="question.Enonce || []"
          @update:model-value="onUpdateEnonce"
          @importQuestion="onImportQuestion"
          @add-syntax-hint="addSyntaxHint"
          :available-parameters="availableParameters"
          :errorBlockIndex="errorIsCorrection ? undefined : errorContent?.Block"
          ref="questionEnonceNode"
        >
        </QuestionContent>
        <QuestionContent
          v-else
          :model-value="question.Correction || []"
          @update:model-value="onUpdateCorrection"
          @importQuestion="onImportQuestion"
          :available-parameters="availableParameters"
          :errorBlockIndex="errorIsCorrection ? errorContent?.Block : undefined"
          ref="questionCorrectionNode"
        >
        </QuestionContent>
      </v-col>
    </v-row>
  </v-card>
</template>

<script setup lang="ts">
import {
  BlockKind,
  ErrorKind,
  type Block,
  type errEnonce,
  type ErrParameters,
  type ErrQuestionInvalid,
  type ExpressionFieldBlock,
  type LoopbackShowQuestion,
  type Parameters,
  type Question,
  type Variable,
} from "@/controller/api_gen";
import { controller } from "@/controller/controller";
import { saveData } from "@/controller/editor";
import { History } from "@/controller/editor_history";
import { copy } from "@/controller/utils";
import { ref } from "@vue/reactivity";
import { computed, onMounted, onUnmounted, watch } from "@vue/runtime-core";
import { $ref } from "vue/macros";
import BlockBar from "../BlockBar.vue";
import SnackErrorParameters from "../parameters/SnackErrorParameters.vue";
import QuestionContent from "../QuestionContent.vue";
import SnackErrorEnonce from "../SnackErrorEnonce.vue";
import ParametersEditor from "../parameters/ParametersEditor.vue";

interface Props {
  question: Question;
  readonly: boolean;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "update", question: Question): void;
  (e: "preview", preview: LoopbackShowQuestion): void;
}>();

let question = $ref(copy(props.question));

let modeEnonce = $ref(true); // false for correction

watch(props, () => {
  question = copy(props.question);
});

onMounted(() => {
  history.addListener();
});

onUnmounted(() => {
  history.clearListener();
});

const hasEditorSimplified = computed(
  () => controller.settings.HasEditorSimplified
);

interface historyEntry {
  question: Question;
}

let history = new History(
  { question }, // start with initial question in history
  controller.showMessage!,
  restoreHistory
);

function restoreHistory(snapshot: historyEntry) {
  question = snapshot.question;
}

let questionEnonceNode = $ref<InstanceType<typeof QuestionContent> | null>(
  null
);
let questionCorrectionNode = $ref<InstanceType<typeof QuestionContent> | null>(
  null
);
function addBlock(kind: BlockKind) {
  if (modeEnonce) {
    questionEnonceNode?.addBlock(kind);
  } else {
    questionCorrectionNode?.addBlock(kind);
  }
}

function onUpdateEnonce(v: Block[]) {
  question.Enonce = v;
  history.add({ question });
}
function onUpdateCorrection(v: Block[]) {
  question.Correction = v;
  history.add({ question });
}

let errorContent = $ref<errEnonce | null>(null);
let errorIsCorrection = $ref(false);

async function save() {
  const res = await controller.EditorSaveQuestionAndPreview({
    Id: question.Id,
    Page: {
      enonce: question.Enonce,
      parameters: question.Parameters,
      correction: question.Correction,
    },
    ShowCorrection: !modeEnonce,
  });
  if (res == undefined) {
    return;
  }

  if (res.IsValid) {
    errorParameters = null;
    errorContent = null;
    // notifie the parent on success
    emit("update", question);
    emit("preview", res.Question);
  } else {
    onQuestionError(res.Error);
  }
}

function onQuestionError(err: ErrQuestionInvalid) {
  // reset previous error
  errorParameters = null;
  errorContent = null;
  switch (err.Kind) {
    case ErrorKind.ErrParameters_:
      errorParameters = err.ErrParameters;
      return;
    case ErrorKind.ErrEnonce:
      errorContent = err.ErrEnonce;
      errorIsCorrection = false;
      modeEnonce = true;
      return;
    case ErrorKind.ErrCorrection:
      errorContent = err.ErrCorrection;
      errorIsCorrection = true;
      modeEnonce = false;
      return;
  }
}

function download() {
  saveData(question, "question.isyro.json");
}

async function onImportQuestion(imported: Question) {
  // only import the data fields
  question.Enonce = imported.Enonce;
  question.Correction = imported.Correction;
  question.Parameters = imported.Parameters;

  history.add({ question });
}

let errorParameters = $ref<ErrParameters | null>(null);
const showErrorParameters = computed(() => errorParameters != null);
const availableParameters = ref<Variable[]>([]);
let isCheckingParameters = $ref(false);

async function checkParameters(params: Parameters) {
  question.Parameters = params;
  history.add({ question });

  isCheckingParameters = true;
  const out = await controller.EditorCheckQuestionParameters({
    Parameters: params,
  });
  isCheckingParameters = false;
  if (out === undefined) return;

  // hide previous error
  errorContent = null;
  errorParameters = out.ErrDefinition.Origin == "" ? null : out.ErrDefinition;

  availableParameters.value = out.Variables || [];
}

async function addSyntaxHint(block: ExpressionFieldBlock) {
  if (questionEnonceNode == null) return;

  const res = await controller.EditorGenerateSyntaxHint({
    Block: block,
    SharedParameters: [],
    QuestionParameters: question.Parameters,
  });
  if (res == undefined) return;

  questionEnonceNode?.addExistingBlock({
    Kind: BlockKind.TextBlock,
    Data: res,
  });
}

async function exportLatex() {
  const res = await controller.EditorQuestionExportLateX({
    parameters: question.Parameters,
    enonce: question.Enonce,
    correction: question.Correction,
  });
  if (res == undefined) return;

  if (res.IsValid) {
    try {
      await navigator.clipboard.writeText(res.Latex);
      if (controller.showMessage)
        controller.showMessage("Question copiée dans le presse-papier");
    } catch (error) {
      if (controller.onError)
        controller.onError(
          "Presse-papier",
          "L'accès au presse-papier a échoué."
        );
    }
  } else {
    onQuestionError(res.Error);
  }
}
</script>

<style scoped></style>
