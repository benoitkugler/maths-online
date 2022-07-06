<template>
  <SnackErrorParameters
    :error="errorParameters"
    @close="errorParameters = null"
  >
  </SnackErrorParameters>

  <!-- <v-snackbar :model-value="showErrorEnnonce" color="warning">
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
  </v-snackbar> -->

  <!-- <v-dialog v-model="showErrVarsDetails">
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
  </v-dialog> -->

  <!-- <v-dialog v-model="showEditDescription">
    <description-pannel
      v-model="question.description"
      :readonly="isReadonly"
    ></description-pannel>
  </v-dialog> -->

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

      <v-col>
        <!-- <v-row no-gutters>
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
              </v-list>
            </v-menu>
          </v-col>
        </v-row> -->

        <!-- <v-row no-gutters>
          <v-col class="pr-2" align-self="center"> ></v-col>
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
        </v-row> -->
      </v-col>

      <v-col cols="auto">
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
        <div style="height: 70vh; overflow-y: auto" class="py-2 px-2">
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
        <!-- <QuestionContent
          :model-value="question.page.enonce || []"
          @update:model-value="(v) => (question.page.enonce = v)"
          @importQuestion="onImportQuestion"
          :available-parameters="availableParameters"
          :errorBlockIndex="errorEnnonce?.Block"
          ref="questionContent"
        >
        </QuestionContent> -->
      </v-col>
    </v-row>
  </v-card>
</template>

<script setup lang="ts">
import type {
  errEnonce,
  ErrParameters,
  ExerciceExt,
} from "@/controller/api_gen";
import { controller } from "@/controller/controller";
import type { RandomParameter } from "@/controller/exercice_gen";
import { computed } from "vue";
import { $computed, $ref } from "vue/macros";
import IntrinsicsParametersExercice from "../editor/IntrinsicsParametersExercice.vue";
import SnackErrorParameters from "../editor/parameters/SnackErrorParameters.vue";
import RandomParametersExercice from "../editor/RandomParametersExercice.vue";

interface Props {
  session_id: string;
  exercice: ExerciceExt;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "back"): void;
}>();

let questionIndex = $ref(0);

let question = $computed(() => (props.exercice.Questions || [])[questionIndex]);

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

let isCheckingParameters = $ref(false);
let errorParameters = $ref<ErrParameters | null>(null);
const showErrorParameters = computed(() => errorParameters != null);

let errorEnnonce = $ref<errEnonce | null>(null);
const showErrorEnnonce = computed(() => errorEnnonce != null);

async function checkParameters() {
  isCheckingParameters = true;
  const out = await controller.EditorCheckExerciceParameters({
    IdExercice: props.exercice.Exercice.Id,
    SharedParameters: props.exercice.Exercice.Parameters,
    QuestionParameters:
      props.exercice.Questions?.map((q) => q.Question.page.parameters) || [],
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
</script>

<style scoped>
.input-small:deep(input) {
  font-size: 14px;
  width: 100%;
}
</style>
