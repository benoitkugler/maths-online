<template>
  <QuestionPageEditor
    v-if="page != null"
    :question="page"
    :readonly="props.readonly"
    :show-dual-parameters="true"
    :on-save="saveQuestion"
    :on-export-latex="exportLatex"
    @update="writeChanges"
    ref="editor"
  >
    <!-- navigation between questions -->
    <template v-slot:top-left>
      <v-row no-gutters class="pr-1">
        <v-col cols="auto" align-self="center">
          <small class="text-grey"> Questions : </small>
        </v-col>
        <v-col cols="8" align-self="center">
          <v-pagination
            :length="(exercice.Questions || []).length"
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
            :disabled="props.readonly"
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
    </template>
  </QuestionPageEditor>

  <v-dialog max-width="800" v-model="showSkeletonDetails">
    <SkeletonDetails
      v-if="exercice != null"
      :exercice="exercice"
      :is-readonly="props.readonly"
      @update="notifieUpdate"
      @preview="updatePreview"
    ></SkeletonDetails>
  </v-dialog>
</template>

<script setup lang="ts">
import {
  DifficultyTag,
  type ExerciceExt,
  type ExerciceHeader,
  type ExerciceQuestionExt,
  type LoopbackShowExercice,
  type TagsDB,
  Int,
} from "@/controller/api_gen";
import { controller } from "@/controller/controller";
import { ref, computed, onMounted, watch } from "vue";
import SkeletonDetails from "./SkeletonDetails.vue";
import { QuestionPage } from "@/controller/editor";
import QuestionPageEditor from "../QuestionPageEditor.vue";
import { LoopbackServerEventKind } from "@/controller/loopback_gen";

interface Props {
  exerciceHeader: ExerciceHeader;
  readonly: boolean;
  allTags: TagsDB; // to provide auto completion
  showVariantMeta: boolean;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "update", ex: ExerciceHeader): void;
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

defineExpose({ refreshExercicePreview });

onMounted(async () => {
  await fetchExercice();
  refreshExercicePreview();
});

watch(
  () => props.exerciceHeader,
  async (newV, oldV) => {
    if (newV.Id != oldV.Id) {
      await fetchExercice();
      refreshExercicePreview();

      // reset the question index if needed
      if (questionIndex.value >= (exercice.value.Questions?.length || 0))
        questionIndex.value = 0 as Int;
    }
  }
);

const showSkeletonDetails = ref(false);

async function refreshExercicePreview() {
  const res = await controller.EditorSaveExerciceAndPreview({
    OnlyPreview: true,
    IdExercice: props.exerciceHeader.Id,
    Parameters: [], // ignored
    Questions: [], // ignored
    CurrentQuestion: -1 as Int,
    ShowCorrection: false,
  });
  if (res == undefined) return;
  updatePreview(res.Preview);
}

const editor = ref<InstanceType<typeof QuestionPageEditor> | null>(null);
function updatePreview(content: LoopbackShowExercice) {
  editor.value?.updatePreview({
    Kind: LoopbackServerEventKind.LoopbackShowExercice,
    Data: content,
  });
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
}

const question = computed<ExerciceQuestionExt | null>(
  () => (exercice.value.Questions || [])[questionIndex.value] || null
);

const page = computed<QuestionPage | null>(() =>
  question.value == null
    ? null
    : {
        id: question.value.Question.Id,
        parameters: question.value.Question.Parameters,
        sharedParameters: exercice.value.Exercice.Parameters,
        enonce: question.value.Question.Enonce,
        correction: question.value.Question.Correction,
      }
);

function writeChanges(page: QuestionPage) {
  const qu = (exercice.value.Questions || [])[questionIndex.value].Question;
  qu.Parameters = page.parameters;
  qu.Enonce = page.enonce;
  qu.Correction = page.correction;
  exercice.value.Exercice.Parameters = page.sharedParameters;
}

async function saveQuestion(showCorrection: boolean) {
  const res = await controller.EditorSaveExerciceAndPreview({
    OnlyPreview: false,
    IdExercice: exercice.value.Exercice.Id,
    Parameters: exercice.value.Exercice.Parameters,
    Questions: exercice.value.Questions?.map((qu) => qu.Question) || [],
    CurrentQuestion: questionIndex.value,
    ShowCorrection: showCorrection,
  });
  if (res === undefined) return;

  if (!res.IsValid) {
    // jump to the invalid question
    // go to the faulty question
    questionIndex.value = res.QuestionIndex;
  }

  return {
    IsValid: res.IsValid,
    Error: res.Error,
    Preview: {
      Kind: LoopbackServerEventKind.LoopbackShowExercice,
      Data: res.Preview,
    },
  };
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

  if (!res.IsValid) {
    // jump to the invalid question
    // go to the faulty question
    questionIndex.value = res.QuestionIndex;
  }

  return { IsValid: res.IsValid, Error: res.Error, Latex: res.Latex };
}
</script>

<style scoped></style>
