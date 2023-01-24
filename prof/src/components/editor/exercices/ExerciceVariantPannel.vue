<template>
  <v-card class="px-2">
    <v-row no-gutters class="mb-2">
      <v-col cols="auto" align-self="center">
        <v-pagination
          :length="(exercice.Questions || []).length"
          :total-visible="(exercice.Questions || []).length"
          :model-value="questionIndex + 1"
          @update:model-value="(oneBased) => (questionIndex = oneBased - 1)"
          density="compact"
        ></v-pagination>
      </v-col>
      <v-col cols="auto" align-self="center">
        <v-btn
          icon
          size="x-small"
          variant="tonal"
          color="success"
          title="Ajouter une question"
          @click="createQuestion"
        >
          <v-icon>mdi-plus</v-icon>
        </v-btn>
        <v-btn
          icon
          size="x-small"
          variant="tonal"
          title="Ordre et barème des questions"
          @click="showSkeletonDetails = true"
        >
          <v-icon>mdi-cog</v-icon>
        </v-btn>
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
          <BlockBar @add="addBlock"></BlockBar>
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
            <v-list-item class="ma-0 pa-0">
              <v-btn
                variant="tonal"
                class="ma-2"
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
            <v-list-item class="ma-0 pa-0">
              <v-btn
                variant="tonal"
                class="ma-2"
                @click="download"
                :disabled="!question.Question.Page.enonce?.length"
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

    <v-row no-gutters v-if="question != null">
      <v-col md="5">
        <div style="height: 68vh; overflow-y: auto" class="py-2 px-2">
          <RandomParametersExercice
            :shared-parameters="exercice.Exercice.Parameters.Variables"
            :question-parameters="question.Question.Page.parameters.Variables"
            :is-loading="isCheckingParameters"
            :is-validated="!showErrorParameters"
            @update="updateRandomParameters"
            @done="checkParameters"
          ></RandomParametersExercice>
          <IntrinsicsParametersExercice
            :shared-parameters="exercice.Exercice.Parameters.Intrinsics || []"
            :question-parameters="
              question.Question.Page.parameters.Intrinsics || []
            "
            :is-loading="isCheckingParameters"
            :is-validated="!showErrorParameters"
            @update="updateIntrinsics"
            @done="checkParameters"
          ></IntrinsicsParametersExercice>
        </div>
      </v-col>
      <v-col class="pr-1">
        <QuestionContent
          :model-value="question.Question.Page.enonce || []"
          @update:model-value="updateQuestion"
          @importQuestion="onImportQuestion"
          :available-parameters="[]"
          :errorBlockIndex="errorEnnonce?.Block"
          ref="questionContent"
        >
        </QuestionContent>
      </v-col>
    </v-row>
  </v-card>

  <SnackErrorParameters
    :error="errorParameters"
    @close="errorParameters = null"
  >
  </SnackErrorParameters>

  <SnackErrorEnonce
    :error="errorEnnonce"
    @close="errorEnnonce = null"
  ></SnackErrorEnonce>

  <v-dialog max-width="800" v-model="showSkeletonDetails">
    <SkeletonDetails
      v-if="exercice != null"
      :exercice="exercice"
      :is-readonly="props.isReadonly"
      @update="notifieUpdate"
      @preview="(qu) => emit('preview', qu)"
    ></SkeletonDetails>
  </v-dialog>

  <v-dialog v-model="showEditDescription">
    <DescriptionPannel
      :description="question.Question.Description"
      @save="saveQuestionDescription"
      :readonly="isReadonly"
    >
    </DescriptionPannel>
  </v-dialog>
</template>

<script setup lang="ts">
import {
  BlockKind,
  DifficultyTag,
  type Block,
  type errEnonce,
  type ErrParameters,
  type ExerciceExt,
  type ExerciceHeader,
  type IdExercice,
  type LoopbackShowExercice,
  type Question,
  type RandomParameter,
  type TagsDB,
} from "@/controller/api_gen";
import { controller } from "@/controller/controller";
import { computed, isReadonly, onMounted, onUnmounted, watch } from "vue";
import { $computed, $ref } from "vue/macros";
import SkeletonDetails from "./SkeletonDetails.vue";
import SnackErrorEnonce from "../SnackErrorEnonce.vue";
import SnackErrorParameters from "../parameters/SnackErrorParameters.vue";
import QuestionContent from "../QuestionContent.vue";
import IntrinsicsParametersExercice from "../IntrinsicsParametersExercice.vue";
import RandomParametersExercice from "../RandomParametersExercice.vue";
import DescriptionPannel from "../DescriptionPannel.vue";
import { History } from "@/controller/editor_history";
import { saveData } from "@/controller/editor";
import BlockBar from "../BlockBar.vue";

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

let questionIndex = $ref(0);

let exercice = $ref<ExerciceExt>({
  Exercice: {
    Id: 0,
    IdGroup: 0,
    Subtitle: "",
    Description: "",
    Parameters: { Variables: [], Intrinsics: [] },
    Difficulty: DifficultyTag.DiffEmpty,
  },
  Questions: [],
});

onMounted(() => {
  fetchExercice();
  refreshExercicePreview(props.exerciceHeader.Id);
});

watch(props, async (_) => {
  if (props.exerciceHeader.Id != exercice?.Exercice.Id) {
    await fetchExercice();
    refreshExercicePreview(props.exerciceHeader.Id);

    // reset the question index if needed
    if (questionIndex >= (exercice.Questions?.length || 0)) questionIndex = 0;
  }
});

async function refreshExercicePreview(id: IdExercice) {
  const res = await controller.EditorSaveExerciceAndPreview({
    OnlyPreview: true,
    IdExercice: id,
    Parameters: { Intrinsics: [], Variables: [] }, // ignored
    Questions: [], // ignored
    CurrentQuestion: -1,
  });
  if (res == undefined) return;
  emit("preview", res.Preview);
}

async function fetchExercice() {
  const res = await controller.EditorGetExerciceContent({
    id: props.exerciceHeader.Id,
  });
  if (res == undefined) return;
  exercice = res;
}

function notifieUpdate(ex: ExerciceExt) {
  exercice = ex;
  // reset the question index if needed
  if (questionIndex >= (exercice.Questions?.length || 0)) questionIndex = 0;
  emit("update", {
    Id: ex.Exercice.Id,
    Difficulty: ex.Exercice.Difficulty,
    Subtitle: ex.Exercice.Subtitle,
  });
}

let question = $computed(
  () => (exercice.Questions || [])[questionIndex] || null
);

let history = new History(exercice, controller.showMessage!, restoreHistory);

onMounted(() => {
  history.addListener();
});
onUnmounted(() => {
  history.clearListener();
});

function restoreHistory(snapshot: ExerciceExt) {
  notifieUpdate(snapshot);
}

let showEditDescription = $ref(false);
async function saveQuestionDescription(description: string) {
  showEditDescription = false;
  question.Question.Description = description;
  const res = await controller.EditorSaveQuestionMeta({
    Question: question.Question,
  });
}

function updateQuestion(qu: Block[]) {
  question.Question.Page.enonce = qu;
  history.add(exercice);
}

function updateRandomParameters(
  sharedP: RandomParameter[],
  questionP: RandomParameter[],
  shouldCheck: boolean
) {
  exercice.Exercice.Parameters.Variables = sharedP;
  question.Question.Page.parameters.Variables = questionP;
  if (shouldCheck) {
    checkParameters();
  }
  history.add(exercice);
}

function updateIntrinsics(
  sharedP: string[],
  questionP: string[],
  shouldCheck: boolean
) {
  exercice.Exercice.Parameters.Intrinsics = sharedP;
  question.Question.Page.parameters.Intrinsics = questionP;
  if (shouldCheck) {
    checkParameters();
  }
  history.add(exercice);
}

let questionContent = $ref<InstanceType<typeof QuestionContent> | null>(null);
function addBlock(kind: BlockKind) {
  if (questionContent == null) {
    return;
  }
  questionContent.addBlock(kind);
}

let isCheckingParameters = $ref(false);
let errorParameters = $ref<ErrParameters | null>(null);
const showErrorParameters = computed(() => errorParameters != null);

let errorEnnonce = $ref<errEnonce | null>(null);

async function checkParameters() {
  isCheckingParameters = true;
  const out = await controller.EditorCheckExerciceParameters({
    IdExercice: exercice.Exercice.Id,
    SharedParameters: exercice.Exercice.Parameters,
    QuestionParameters:
      exercice.Questions?.map((q) => q.Question.Page.parameters) || [],
  });
  isCheckingParameters = false;
  if (out === undefined) return;

  // hide previous error
  errorEnnonce = null;

  errorParameters = out.ErrDefinition.Origin == "" ? null : out.ErrDefinition;
  if (errorParameters != null) {
    // go to faulty question
    questionIndex = out.QuestionIndex;
  }
  //   availableParameters.value = out.Variables || [];
}

async function save() {
  const res = await controller.EditorSaveExerciceAndPreview({
    OnlyPreview: false,
    IdExercice: exercice.Exercice.Id,
    Parameters: exercice.Exercice.Parameters,
    Questions: exercice.Questions?.map((qu) => qu.Question) || [],
    CurrentQuestion: questionIndex,
  });
  if (res == undefined) {
    return;
  }

  if (res.IsValid) {
    errorEnnonce = null;
    errorParameters = null;

    notifieUpdate(exercice);
    emit("preview", res.Preview);
  } else {
    if (res.Error.ParametersInvalid) {
      errorEnnonce = null;
      errorParameters = res.Error.ErrParameters;
    } else {
      errorEnnonce = res.Error.ErrEnonce;
      errorParameters = null;
    }
    // go to the faulty question
    questionIndex = res.QuestionIndex;
  }
}

function download() {
  saveData<Question>(
    question.Question,
    `question${questionIndex + 1}.isyro.json`
  );
}

async function onImportQuestion(imported: Question) {
  // only import the data fields
  question.Question.Page = imported.Page;

  history.add(exercice);

  notifieUpdate(exercice);
}

async function createQuestion() {
  const res = await controller.EditorExerciceCreateQuestion({
    IdExercice: exercice.Exercice.Id,
  });
  if (res == undefined) {
    return;
  }
  // go to the new question
  questionIndex = (res.Ex.Questions?.length || 0) - 1;
  notifieUpdate(res.Ex);
  emit("preview", res.Preview);
}

let showSkeletonDetails = $ref(false);
</script>

<style scoped>
.input-small:deep(input) {
  font-size: 14px;
  width: 100%;
}
</style>
