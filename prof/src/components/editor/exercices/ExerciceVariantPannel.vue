<template>
  <v-card class="mt-1 px-2">
    <v-row no-gutters>
      <v-col md="5" align-self="center">
        <v-row no-gutters>
          <v-col cols="auto" align-self="center">
            <v-pagination
              :length="(exercice.Questions || []).length"
              :total-visible="(exercice.Questions || []).length"
              :model-value="questionIndex + 1"
              @update:model-value="(oneBased: number) => (questionIndex = oneBased - 1 as Int)"
              density="compact"
            ></v-pagination>
          </v-col>
          <v-col cols="auto" align-self="center">
            <v-btn
              icon
              size="x-small"
              variant="elevated"
              title="Ajouter une question"
              @click="createQuestion"
              class="mr-1"
              :disabled="props.isReadonly"
            >
              <v-icon color="success">mdi-plus</v-icon>
            </v-btn>
            <v-btn
              icon
              size="x-small"
              variant="elevated"
              title="Ordre et barème des questions"
              @click="showSkeletonDetails = true"
            >
              <v-icon>mdi-cog</v-icon>
            </v-btn>
          </v-col>
        </v-row>
      </v-col>

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
              :disabled="!question"
            >
              <v-icon icon="mdi-plus" color="green"></v-icon>
              Insérer du contenu
            </v-btn>
          </template>
          <BlockBar
            @add="addBlock"
            :simplified="hasEditorSimplified"
            :hide-answer-fields="!modeEnonce"
          ></BlockBar>
        </v-menu>
      </v-col>

      <v-col cols="auto" align-self="center" class="py-1">
        <v-btn
          class="mx-2"
          icon
          @click="save"
          :title="
            props.isReadonly ? 'Visualiser' : 'Enregistrer et prévisualiser'
          "
          size="small"
          :disabled="!question"
        >
          <v-icon
            :icon="props.isReadonly ? 'mdi-eye' : 'mdi-content-save'"
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
                @click="paste"
                title="Coller le bloc"
              >
                <v-icon
                  class="mr-2"
                  icon="mdi-content-paste"
                  size="small"
                ></v-icon>
                Coller
              </v-btn>
            </v-list-item>

            <v-list-item>
              <v-btn
                size="small"
                class="my-1"
                @click="download"
                :disabled="!question?.Question.Enonce?.length"
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
                title="Exporter l'exercice au format LaTeX (.tex)"
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

    <v-row no-gutters v-if="question != null">
      <v-col md="5" v-if="!hasEditorSimplified">
        <ParametersEditor
          :model-value="currentEditedParams"
          @update:model-value="checkParameters"
          :is-loading="isCheckingParameters"
          :is-validated="!showErrorParameters"
          :show-switch="editSharedParams"
          @switch="switchParamsMode"
        ></ParametersEditor>
      </v-col>
      <v-col class="pr-1">
        <QuestionContent
          v-if="modeEnonce"
          :model-value="question.Question.Enonce || []"
          @update:model-value="onUpdateEnonce"
          @importQuestion="onImportQuestion"
          @add-syntax-hint="addSyntaxHint"
          :available-parameters="[]"
          :errorBlockIndex="errorIsCorrection ? undefined : errorContent?.Block"
          ref="questionEnonceNode"
        >
        </QuestionContent>
        <QuestionContent
          v-else
          :model-value="question.Question.Correction || []"
          @update:model-value="onUpdateCorrection"
          @importQuestion="onImportQuestion"
          :available-parameters="[]"
          :errorBlockIndex="errorIsCorrection ? errorContent?.Block : undefined"
          ref="questionCorrectionNode"
        >
        </QuestionContent>
      </v-col>
    </v-row>
    <v-row v-else justify="center" style="height: 70vh">
      <v-col cols="auto" align-self="center">
        <v-btn
          title="Ajouter une question"
          @click="createQuestion"
          :disabled="props.isReadonly"
        >
          <v-icon color="success" class="mr-1">mdi-plus</v-icon>
          Ajouter une question</v-btn
        >
      </v-col>
    </v-row>
  </v-card>

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

  <v-dialog max-width="800" v-model="showSkeletonDetails">
    <SkeletonDetails
      v-if="exercice != null"
      :exercice="exercice"
      :is-readonly="props.isReadonly"
      @update="notifieUpdate"
      @preview="(qu: LoopbackShowExercice) => emit('preview', qu)"
    ></SkeletonDetails>
  </v-dialog>
</template>

<script setup lang="ts">
import {
  BlockKind,
  DifficultyTag,
  ErrorKind,
  type Block,
  type errEnonce,
  type ErrParameters,
  type ErrQuestionInvalid,
  type ExerciceExt,
  type ExerciceHeader,
  type ExerciceQuestionExt,
  type ExpressionFieldBlock,
  type IdExercice,
  type LoopbackShowExercice,
  type Parameters,
  type Question,
  type TagsDB,
  Int,
} from "@/controller/api_gen";
import { controller } from "@/controller/controller";
import { ref, computed, onMounted, onUnmounted, watch } from "vue";
import SkeletonDetails from "./SkeletonDetails.vue";
import SnackErrorEnonce from "../SnackErrorEnonce.vue";
import SnackErrorParameters from "../parameters/SnackErrorParameters.vue";
import QuestionContent from "../QuestionContent.vue";
import { History } from "@/controller/editor_history";
import { saveData, readClipboardForBlock } from "@/controller/editor";
import BlockBar from "../BlockBar.vue";
import ParametersEditor from "../parameters/ParametersEditor.vue";

interface Props {
  exerciceHeader: ExerciceHeader;
  isReadonly: boolean;
  allTags: TagsDB; // to provide auto completion
  showVariantMeta: boolean;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "update", ex: ExerciceHeader): void;
  (e: "preview", ex: LoopbackShowExercice): void;
}>();

const questionIndex = ref(0 as Int);

const exercice = ref<ExerciceExt>({
  Exercice: {
    Id: 0 as Int,
    IdGroup: 0 as Int,
    Subtitle: "",
    Parameters: [],
    Difficulty: DifficultyTag.DiffEmpty,
  },
  Questions: [],
});

const modeEnonce = ref(true); // false for correction

onMounted(async () => {
  await fetchExercice();
  refreshExercicePreview(props.exerciceHeader.Id);
  editSharedParams.value = inferEditSharedParams();
});

watch(props, async (_) => {
  if (props.exerciceHeader.Id != exercice.value?.Exercice.Id) {
    await fetchExercice();
    refreshExercicePreview(props.exerciceHeader.Id);

    // reset the question index if needed
    if (questionIndex.value >= (exercice.value.Questions?.length || 0))
      questionIndex.value = 0 as Int;

    editSharedParams.value = inferEditSharedParams();
  }
});

async function refreshExercicePreview(id: IdExercice) {
  const res = await controller.EditorSaveExerciceAndPreview({
    OnlyPreview: true,
    IdExercice: id,
    Parameters: [], // ignored
    Questions: [], // ignored
    CurrentQuestion: -1 as Int,
    ShowCorrection: false,
  });
  if (res == undefined) return;
  emit("preview", res.Preview);
}

async function fetchExercice() {
  const res = await controller.EditorGetExerciceContent({
    id: props.exerciceHeader.Id,
  });
  if (res == undefined) return;
  exercice.value = res;
}

function notifieUpdate(ex: ExerciceExt) {
  exercice.value = ex;
  // reset the question index if needed
  if (questionIndex.value >= (exercice.value.Questions?.length || 0))
    questionIndex.value = 0 as Int;
  emit("update", {
    Id: ex.Exercice.Id,
    Difficulty: ex.Exercice.Difficulty,
    Subtitle: ex.Exercice.Subtitle,
  });
}

// guess the edit mode, defaulting on shared for empty params
function inferEditSharedParams() {
  if (
    !exercice.value.Exercice.Parameters?.length &&
    !question.value?.Question.Parameters?.length
  ) {
    return true;
  }
  return !!exercice.value.Exercice.Parameters?.length;
}

const editSharedParams = ref(true);

// either the shared params or the params of the current question
const currentEditedParams = computed(() =>
  editSharedParams.value
    ? exercice.value.Exercice.Parameters
    : question.value?.Question.Parameters || []
);

function switchParamsMode(b: boolean) {
  editSharedParams.value = b;
  console.log(currentEditedParams.value);
}

const question = computed<ExerciceQuestionExt | null>(
  () => (exercice.value.Questions || [])[questionIndex.value] || null
);

let history = new History(
  exercice.value,
  controller.showMessage!,
  restoreHistory
);

onMounted(() => {
  history.addListener();
});
onUnmounted(() => {
  history.clearListener();
});

const hasEditorSimplified = computed(
  () => controller.settings.HasEditorSimplified
);

function restoreHistory(snapshot: ExerciceExt) {
  notifieUpdate(snapshot);
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
  if (!question.value) return;
  question.value.Question.Enonce = v;
  history.add(exercice.value);
}
function onUpdateCorrection(v: Block[]) {
  if (!question.value) return;
  question.value.Question.Correction = v;
  history.add(exercice.value);
}

const isCheckingParameters = ref(false);
const errorParameters = ref<ErrParameters | null>(null);
const showErrorParameters = computed(() => errorParameters.value != null);

const errorContent = ref<errEnonce | null>(null);
const errorIsCorrection = ref(false);

async function checkParameters(ps: Parameters) {
  if (editSharedParams.value) {
    exercice.value.Exercice.Parameters = ps;
  } else {
    question.value!.Question.Parameters = ps;
  }
  history.add(exercice.value);

  isCheckingParameters.value = true;
  const out = await controller.EditorCheckExerciceParameters({
    IdExercice: exercice.value.Exercice.Id,
    SharedParameters: exercice.value.Exercice.Parameters,
    QuestionParameters:
      exercice.value.Questions?.map((q) => q.Question.Parameters) || [],
  });
  isCheckingParameters.value = false;
  if (out === undefined) return;

  // hide previous error
  errorContent.value = null;

  errorParameters.value =
    out.ErrDefinition.Origin == "" ? null : out.ErrDefinition;
  if (errorParameters.value != null) {
    // go to faulty question
    questionIndex.value = out.QuestionIndex;
  }
  //   availableParameters.value = out.Variables || [];
}

async function save() {
  const res = await controller.EditorSaveExerciceAndPreview({
    OnlyPreview: false,
    IdExercice: exercice.value.Exercice.Id,
    Parameters: exercice.value.Exercice.Parameters,
    Questions: exercice.value.Questions?.map((qu) => qu.Question) || [],
    CurrentQuestion: questionIndex.value,
    ShowCorrection: !modeEnonce.value,
  });
  if (res == undefined) {
    return;
  }

  if (res.IsValid) {
    errorContent.value = null;
    errorParameters.value = null;

    notifieUpdate(exercice.value);
    emit("preview", res.Preview);
  } else {
    onQuestionError(res.QuestionIndex, res.Error);
  }
}

function onQuestionError(index: Int, err: ErrQuestionInvalid) {
  // reset previous error
  errorParameters.value = null;
  errorContent.value = null;
  switch (err.Kind) {
    case ErrorKind.ErrParameters_:
      errorParameters.value = err.ErrParameters;
      break;
    case ErrorKind.ErrEnonce:
      errorContent.value = err.ErrEnonce;
      errorIsCorrection.value = false;
      modeEnonce.value = true;
      break;
    case ErrorKind.ErrCorrection:
      errorContent.value = err.ErrCorrection;
      errorIsCorrection.value = true;
      modeEnonce.value = false;
      break;
  }

  // go to the faulty question
  questionIndex.value = index;
}

function download() {
  if (!question.value) return;
  saveData<Question>(
    question.value.Question,
    `question${questionIndex.value + 1}.isyro.json`
  );
}

async function onImportQuestion(imported: Question) {
  if (!question.value) return;
  // only import the data fields
  question.value.Question.Enonce = imported.Enonce;
  question.value.Question.Correction = imported.Correction;
  question.value.Question.Parameters = imported.Parameters;

  history.add(exercice.value);

  notifieUpdate(exercice.value);
}

async function createQuestion() {
  const res = await controller.EditorExerciceCreateQuestion({
    IdExercice: exercice.value.Exercice.Id,
  });
  if (res == undefined) {
    return;
  }
  // go to the new question
  questionIndex.value = ((res.Ex.Questions?.length || 0) - 1) as Int;
  notifieUpdate(res.Ex);
  emit("preview", res.Preview);
}

const showSkeletonDetails = ref(false);

async function addSyntaxHint(block: ExpressionFieldBlock) {
  if (questionEnonceNode.value == null || question.value == null) return;

  const res = await controller.EditorGenerateSyntaxHint({
    Block: block,
    SharedParameters: exercice.value.Exercice.Parameters,
    QuestionParameters: question.value.Question.Parameters,
  });
  if (res == undefined) return;

  questionEnonceNode.value?.addExistingBlock({
    Kind: BlockKind.TextBlock,
    Data: res,
  });
}

async function exportLatex() {
  const res = await controller.EditorExerciceExportLateX({
    Questions:
      exercice.value.Questions?.map((qu) => ({
        enonce: qu.Question.Enonce,
        parameters: qu.Question.Parameters,
        correction: qu.Question.Correction,
      })) || [],
    Parameters: exercice.value.Exercice.Parameters,
  });
  if (res == undefined) return;

  if (res.IsValid) {
    try {
      await navigator.clipboard.writeText(res.Latex);
      if (controller.showMessage)
        controller.showMessage("Exercice copié dans le presse-papier");
    } catch (error) {
      if (controller.onError)
        controller.onError(
          "Presse-papier",
          "L'accès au presse-papier a échoué."
        );
    }
  } else {
    onQuestionError(res.QuestionIndex, res.Error);
  }
}

async function paste() {
  const block = await readClipboardForBlock();
  if (block === undefined) return;
  if (modeEnonce.value) {
    questionEnonceNode.value?.addExistingBlock(block);
  } else {
    questionCorrectionNode.value?.addExistingBlock(block);
  }
}
</script>

<style scoped></style>
