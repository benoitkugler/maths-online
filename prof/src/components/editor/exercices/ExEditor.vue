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
    <DescriptionPannel
      :description="question.Question.Description"
      @save="saveQuestionDescription"
      :readonly="isReadonly"
    >
    </DescriptionPannel>
  </v-dialog>

  <v-card class="px-2">
    <v-row no-gutters class="mb-2">
      <v-col cols="auto" align-self="center">
        <v-btn
          size="small"
          icon
          :title="
            questionIndex == 0 ? 'Retour' : 'Aller à la question précédente'
          "
          @click="goToPrevious"
        >
          <v-icon icon="mdi-arrow-left" size="small"></v-icon>
        </v-btn>
        Question {{ questionIndex + 1 }} /
        {{ (props.exercice.Questions || []).length }}
        <v-btn
          size="small"
          icon
          title="Aller à la question suivante"
          @click="questionIndex = questionIndex + 1"
          :disabled="
            questionIndex >= (props.exercice.Questions || []).length - 1
          "
        >
          <v-icon icon="mdi-arrow-right" size="small"></v-icon>
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
          :title="isReadonly ? 'Visualiser' : 'Enregistrer et prévisualiser'"
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

    <v-row no-gutters>
      <v-col md="5">
        <div style="height: 68vh; overflow-y: auto" class="py-2 px-2">
          <RandomParametersExercice
            :shared-parameters="props.exercice.Exercice.Parameters.Variables"
            :question-parameters="question.Question.Page.parameters.Variables"
            :is-loading="isCheckingParameters"
            :is-validated="!showErrorParameters"
            @update="updateRandomParameters"
            @done="checkParameters"
          ></RandomParametersExercice>
          <IntrinsicsParametersExercice
            :shared-parameters="
              props.exercice.Exercice.Parameters.Intrinsics || []
            "
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
</template>

<script setup lang="ts">
import type {
  Block,
  BlockKind,
  errEnonce,
  ErrParameters,
  ExerciceExt,
  LoopbackShowExercice,
  Question,
  RandomParameter,
} from "@/controller/api_gen";
import { controller } from "@/controller/controller";
import { saveData } from "@/controller/editor";
import { History } from "@/controller/editor_history";
import { computed, onMounted, onUnmounted } from "vue";
import { $computed, $ref } from "vue/macros";
import BlockBar from "../BlockBar.vue";
import DescriptionPannel from "../DescriptionPannel.vue";
import IntrinsicsParametersExercice from "../IntrinsicsParametersExercice.vue";
import SnackErrorParameters from "../parameters/SnackErrorParameters.vue";
import QuestionContent from "../QuestionContent.vue";
import RandomParametersExercice from "../RandomParametersExercice.vue";
import SnackErrorEnonce from "../SnackErrorEnonce.vue";

interface Props {
  exercice: ExerciceExt;
  questionIndex: number;
  isReadonly: boolean;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "update", ex: ExerciceExt): void;
  (e: "preview", qu: LoopbackShowExercice): void;
  (e: "back"): void;
}>();

let questionIndex = $ref(props.questionIndex);

let question = $computed(() => (props.exercice.Questions || [])[questionIndex]);

let history = new History(
  props.exercice,
  controller.showMessage!,
  restoreHistory
);

onMounted(() => {
  history.addListener();
});
onUnmounted(() => {
  history.clearListener();
});

function restoreHistory(snapshot: ExerciceExt) {
  emit("update", snapshot);
}

function goToPrevious() {
  if (questionIndex == 0) {
    emit("back");
  } else {
    questionIndex = questionIndex - 1;
  }
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
  history.add(props.exercice);
}

function updateRandomParameters(
  sharedP: RandomParameter[],
  questionP: RandomParameter[],
  shouldCheck: boolean
) {
  props.exercice.Exercice.Parameters.Variables = sharedP;
  question.Question.Page.parameters.Variables = questionP;
  if (shouldCheck) {
    checkParameters();
  }
  history.add(props.exercice);
}

function updateIntrinsics(
  sharedP: string[],
  questionP: string[],
  shouldCheck: boolean
) {
  props.exercice.Exercice.Parameters.Intrinsics = sharedP;
  question.Question.Page.parameters.Intrinsics = questionP;
  if (shouldCheck) {
    checkParameters();
  }
  history.add(props.exercice);
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
    IdExercice: props.exercice.Exercice.Id,
    SharedParameters: props.exercice.Exercice.Parameters,
    QuestionParameters:
      props.exercice.Questions?.map((q) => q.Question.Page.parameters) || [],
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
    IdExercice: props.exercice.Exercice.Id,
    Parameters: props.exercice.Exercice.Parameters,
    Questions: props.exercice.Questions?.map((qu) => qu.Question) || [],
    CurrentQuestion: questionIndex,
  });
  if (res == undefined) {
    return;
  }

  if (res.IsValid) {
    errorEnnonce = null;
    errorParameters = null;

    emit("update", props.exercice);
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

  history.add(props.exercice);

  emit("update", props.exercice);
}
</script>

<style scoped>
.input-small:deep(input) {
  font-size: 14px;
  width: 100%;
}
</style>
