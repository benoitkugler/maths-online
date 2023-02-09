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

  <v-card class="mt-1 pr-1">
    <v-row no-gutters>
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
          <block-bar @add="addBlock"></block-bar>
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
          </v-list>
        </v-menu>
      </v-col>
    </v-row>

    <v-row no-gutters>
      <v-col md="5">
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
          :model-value="question.Enonce || []"
          @update:model-value="onUpdateEnonce"
          @importQuestion="onImportQuestion"
          @add-syntax-hint="addSyntaxHint"
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
import {
  BlockKind,
  type Block,
  type errEnonce,
  type ErrParameters,
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

watch(props, () => {
  question = copy(props.question);
});

onMounted(() => {
  history.addListener();
});

onUnmounted(() => {
  history.clearListener();
});

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

let questionContent = $ref<InstanceType<typeof QuestionContent> | null>(null);
function addBlock(kind: BlockKind) {
  if (questionContent == null) return;
  questionContent.addBlock(kind);
}

function onUpdateEnonce(v: Block[]) {
  question.Enonce = v;
  history.add({ question });
}

let errorEnnonce = $ref<errEnonce | null>(null);

async function save() {
  const res = await controller.EditorSaveQuestionAndPreview({
    Id: question.Id,
    Page: { enonce: question.Enonce, parameters: question.Parameters },
  });
  if (res == undefined) {
    return;
  }

  if (res.IsValid) {
    errorEnnonce = null;
    errorParameters = null;
    // notifie the parent on success
    emit("update", question);
    emit("preview", res.Question);
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
  // only import the data fields
  question.Enonce = imported.Enonce;
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
  errorEnnonce = null;

  errorParameters = out.ErrDefinition.Origin == "" ? null : out.ErrDefinition;
  availableParameters.value = out.Variables || [];
}

async function addSyntaxHint(block: ExpressionFieldBlock) {
  if (questionContent == null) return;

  const res = await controller.EditorGenerateSyntaxHint({
    Block: block,
    SharedParameters: [],
    QuestionParameters: question.Parameters,
  });
  if (res == undefined) return;

  questionContent?.addExistingBlock({ Kind: BlockKind.TextBlock, Data: res });
}
</script>

<style scoped></style>
