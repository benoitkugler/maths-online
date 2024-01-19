<template>
  <div>
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

    <v-row no-gutters>
      <v-col md="5"> </v-col>
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
            :simplified="false"
            :hide-answer-fields="!modeEnonce"
          ></block-bar>
        </v-menu>
      </v-col>

      <v-col cols="auto" align-self="center">
        <v-btn
          class="my-1 mx-2"
          icon
          @click="emit('save', question)"
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
  </div>
</template>

<script setup lang="ts">
import {
  Beltquestion,
  Block,
  BlockKind,
  ErrParameters,
  ExpressionFieldBlock,
  Parameters,
  Question,
  Variable,
  errEnonce,
} from "@/controller/api_gen";
import { ref } from "vue";
import { controller } from "@/controller/controller";
import { onMounted } from "vue";
import BlockBar from "../editor/BlockBar.vue";
import ParametersEditor from "../editor/parameters/ParametersEditor.vue";
import QuestionContent from "../editor/QuestionContent.vue";
import { computed } from "vue";
import { History } from "@/controller/editor_history";
import { copy } from "@/controller/utils";
import SnackErrorParameters from "../editor/parameters/SnackErrorParameters.vue";
import SnackErrorEnonce from "../editor/SnackErrorEnonce.vue";
import { onUnmounted } from "vue";
import { watch } from "vue";

interface Props {
  question: Beltquestion;
  readonly: boolean;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "save", qu: Beltquestion): void;
}>();

const question = ref(copy(props.question));
watch(props, () => (question.value = copy(props.question)));

interface historyEntry {
  question: Beltquestion;
}

let history = new History<historyEntry>(
  { question: copy(question.value) }, // start with initial question in history
  controller.showMessage,
  restoreHistory
);

function restoreHistory(snapshot: historyEntry) {
  question.value = snapshot.question;
}

onMounted(() => {
  history.addListener();
});

onUnmounted(() => {
  history.clearListener();
});

const modeEnonce = ref(true); // false for correction

const errorContent = ref<errEnonce | null>(null);
const errorIsCorrection = ref(false);
const errorParameters = ref<ErrParameters | null>(null);
const showErrorParameters = computed(() => errorParameters.value != null);
const availableParameters = ref<Variable[]>([]);
const isCheckingParameters = ref(false);
async function checkParameters(params: Parameters) {
  question.value.Parameters = params;
  history.add({ question: copy(question.value) });

  isCheckingParameters.value = true;
  const out = await controller.EditorCheckQuestionParameters({
    Parameters: params,
  });
  isCheckingParameters.value = false;
  if (out === undefined) return;

  // hide previous error
  errorContent.value = null;
  errorParameters.value =
    out.ErrDefinition.Origin == "" ? null : out.ErrDefinition;

  availableParameters.value = out.Variables || [];
}

const questionEnonceNode = ref<InstanceType<typeof QuestionContent> | null>(
  null
);
const questionCorrectionNode = ref<InstanceType<typeof QuestionContent> | null>(
  null
);
function addBlock(kind: BlockKind) {
  if (modeEnonce.value) {
    questionEnonceNode.value?.addBlock(kind);
  } else {
    questionCorrectionNode.value?.addBlock(kind);
  }
}
function onUpdateEnonce(v: Block[]) {
  question.value.Enonce = v;
  history.add({ question: copy(question.value) });
}
function onUpdateCorrection(v: Block[]) {
  question.value.Correction = v;
  history.add({ question: copy(question.value) });
}
async function onImportQuestion(imported: Question) {
  // only import the data fields
  question.value.Enonce = imported.Enonce;
  question.value.Correction = imported.Correction;
  question.value.Parameters = imported.Parameters;

  history.add({ question: copy(question.value) });
}
async function addSyntaxHint(block: ExpressionFieldBlock) {
  if (questionEnonceNode.value == null) return;

  const res = await controller.EditorGenerateSyntaxHint({
    Block: block,
    SharedParameters: [],
    QuestionParameters: question.value.Parameters,
  });
  if (res == undefined) return;

  questionEnonceNode.value?.addExistingBlock({
    Kind: BlockKind.TextBlock,
    Data: res,
  });
}
</script>
