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
      v-model="question.Question.description"
      :readonly="isReadonly"
    >
    </DescriptionPannel>
  </v-dialog>

  <v-card class="mt-3 px-2">
    <v-row no-gutters class="mb-2">
      <v-col cols="auto" align-self="center" class="pr-2">
        <v-btn
          size="small"
          icon
          title="Retour à la liste des questions"
          @click="emit('back')"
        >
          <v-icon icon="mdi-arrow-left"></v-icon>
        </v-btn>
      </v-col>

      <v-col align-self="center">
        <v-text-field
          class="my-2 input-small"
          variant="outlined"
          density="compact"
          label="Nom de la question"
          v-model="question.Question.page.title"
          :readonly="isReadonly"
          hide-details
        ></v-text-field
      ></v-col>

      <v-col cols="3" align-self="center" class="px-1">
        <v-row no-gutters justify="center">
          <v-col cols="auto" align-self="center" class="py-1">
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
                    :disabled="!question.Question.page.enonce?.length"
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
              </v-list>
            </v-menu>
          </v-col>
        </v-row>

        <v-row no-gutters>
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
              <BlockBar @add="addBlock"></BlockBar>
            </v-menu>
          </v-col>
        </v-row>
      </v-col>

      <v-col cols="auto" align-self="center">
        <v-btn
          size="small"
          icon
          title="Aller à la question précédente"
          @click="questionIndex = questionIndex - 1"
          :disabled="questionIndex <= 0"
        >
          <v-icon icon="mdi-arrow-left" size="small"></v-icon>
        </v-btn>
        {{ questionIndex + 1 }} / {{ (props.exercice.Questions || []).length }}
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
    </v-row>

    <v-row no-gutters>
      <v-col md="5">
        <div style="height: 68vh; overflow-y: auto" class="py-2 px-2">
          <RandomParametersExercice
            :shared-parameters="props.exercice.Exercice.Parameters.Variables"
            :question-parameters="question.Question.page.parameters.Variables"
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
              question.Question.page.parameters.Intrinsics || []
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
          :model-value="question.Question.page.enonce || []"
          @update:model-value="(v) => (question.Question.page.enonce = v)"
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
import {
  BlockKind,
  Visibility,
  type errEnonce,
  type ErrParameters,
  type ExerciceExt,
  type Question,
} from "@/controller/api_gen";
import { controller } from "@/controller/controller";
import { saveData } from "@/controller/editor";
import type { RandomParameter } from "@/controller/exercice_gen";
import { computed } from "vue";
import { $computed, $ref } from "vue/macros";
import BlockBar from "../editor/BlockBar.vue";
import DescriptionPannel from "../editor/DescriptionPannel.vue";
import IntrinsicsParametersExercice from "../editor/IntrinsicsParametersExercice.vue";
import SnackErrorParameters from "../editor/parameters/SnackErrorParameters.vue";
import QuestionContent from "../editor/QuestionContent.vue";
import RandomParametersExercice from "../editor/RandomParametersExercice.vue";
import SnackErrorEnonce from "../editor/SnackErrorEnonce.vue";

interface Props {
  session_id: string;
  exercice: ExerciceExt;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "back"): void;
  (e: "update", ex: ExerciceExt): void;
}>();

let questionIndex = $ref(0);

let question = $computed(() => {
  const questionID = (props.exercice.Questions || [])[questionIndex]
    .id_question;
  return getQuestion(questionID);
});

function getQuestion(questionID: number) {
  return props.exercice.QuestionsSource![questionID];
}

const isReadonly = computed(
  () => question.Origin.Visibility != Visibility.Personnal
);

function updateRandomParameters(
  sharedP: RandomParameter[],
  questionP: RandomParameter[],
  shouldCheck: boolean
) {
  props.exercice.Exercice.Parameters.Variables = sharedP;
  question.Question.page.parameters.Variables = questionP;
  if (shouldCheck) {
    checkParameters();
  }
}

function updateIntrinsics(
  sharedP: string[],
  questionP: string[],
  shouldCheck: boolean
) {
  props.exercice.Exercice.Parameters.Intrinsics = sharedP;
  question.Question.page.parameters.Intrinsics = questionP;
  if (shouldCheck) {
    checkParameters();
  }
}

let showEditDescription = $ref(false);

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
      props.exercice.Questions?.map(
        (q) => getQuestion(q.id_question).Question.page.parameters
      ) || [],
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
    SessionID: props.session_id || "",
    IdExercice: props.exercice.Exercice.Id,
    Parameters: props.exercice.Exercice.Parameters,
    Questions: Object.fromEntries(
      Object.entries(props.exercice.QuestionsSource || {}).map((k) => [
        k[0],
        k[1].Question,
      ])
    ),
  });
  if (res == undefined) {
    return;
  }

  if (res.IsValid) {
    errorEnnonce = null;
    errorParameters = null;

    emit("update", props.exercice);
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
  // keep the current ID
  imported.id = question.Question.id;
  question.Question = imported;
  emit("update", props.exercice);
}
</script>

<style scoped>
.input-small:deep(input) {
  font-size: 14px;
  width: 100%;
}
</style>
