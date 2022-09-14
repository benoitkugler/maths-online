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
      :description="question.Description"
      :readonly="props.readonly"
      @save="saveDescription"
    ></description-pannel>
  </v-dialog>

  <v-card class="mt-3 px-2">
    <v-row no-gutters class="mb-2">
      <v-col>
        <v-text-field
          class="my-2 input-small"
          variant="outlined"
          density="compact"
          label="Sous-titre de la variante (optionnel)"
          v-model="question.Subtitle"
          :readonly="props.readonly"
          hide-details
          @blur="saveMeta"
        ></v-text-field
      ></v-col>

      <v-col align-self="center">
        <DifficultyField
          class="px-1"
          v-model="question.Difficulty"
          @update:model-value="saveMeta"
          :readonly="props.readonly"
        ></DifficultyField>
      </v-col>

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
          class="mx-2"
          icon
          @click="save"
          :disabled="!session_id"
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
                <v-icon class="mr-2" icon="mdi-download" size="small"></v-icon>
                Télécharger
              </v-btn>
            </v-list-item>
          </v-list>
        </v-menu>
      </v-col>
    </v-row>

    <v-row no-gutters>
      <v-col md="4">
        <div style="height: 66vh; overflow-y: auto" class="py-2 px-2">
          <RandomParametersQuestion
            :parameters="question.Page.parameters.Variables"
            :is-loading="isCheckingParameters"
            :is-validated="!showErrorParameters"
            @update="updateRandomParameters"
            @done="checkParameters"
          ></RandomParametersQuestion>
          <IntrinsicsParametersQuestion
            :parameters="question.Page.parameters.Intrinsics || []"
            :is-loading="isCheckingParameters"
            :is-validated="!showErrorParameters"
            @update="updateIntrinsics"
            @done="checkParameters"
          ></IntrinsicsParametersQuestion>
        </div>
      </v-col>
      <v-col class="pr-1">
        <QuestionContent
          :model-value="question.Page.enonce || []"
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
  errEnonce,
  ErrParameters,
  Question,
  RandomParameter,
  Variable,
} from "@/controller/api_gen";
import { controller } from "@/controller/controller";
import { saveData } from "@/controller/editor";
import { History } from "@/controller/editor_history";
import { copy } from "@/controller/utils";
import { ref } from "@vue/reactivity";
import { computed, onMounted, onUnmounted, watch } from "@vue/runtime-core";
import { $ref } from "vue/macros";
import BlockBar from "../BlockBar.vue";
import DescriptionPannel from "../DescriptionPannel.vue";
import IntrinsicsParametersQuestion from "../IntrinsicsParametersQuestion.vue";
import SnackErrorParameters from "../parameters/SnackErrorParameters.vue";
import QuestionContent from "../QuestionContent.vue";
import RandomParametersQuestion from "../RandomParametersQuestion.vue";
import SnackErrorEnonce from "../SnackErrorEnonce.vue";
import DifficultyField from "../utils/DifficultyField.vue";

interface Props {
  session_id: string;
  question: Question;
  readonly: boolean;
  allTags: string[]; // to provide auto completion
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "update", question: Question): void;
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
  if (questionContent == null) {
    return;
  }
  questionContent.addBlock(kind);
}

function onUpdateEnonce(v: Block[]) {
  question.Page.enonce = v;
  history.add({ question });
}

function updateRandomParameters(l: RandomParameter[], shouldCheck: boolean) {
  question.Page.parameters.Variables = l;
  if (shouldCheck) {
    checkParameters();
  }

  history.add({ question });
}

function updateIntrinsics(l: string[], shouldCheck: boolean) {
  question.Page.parameters.Intrinsics = l;
  if (shouldCheck) {
    checkParameters();
  }

  history.add({ question });
}

let errorEnnonce = $ref<errEnonce | null>(null);

async function saveMeta() {
  if (props.readonly) {
    return;
  }
  await controller.EditorSaveQuestionMeta({
    Question: question,
  });
  emit("update", question);
}

async function saveDescription(desc: string) {
  showEditDescription = false;
  question.Description = desc;
  saveMeta();
}

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
  imported.Id = question.Id;
  question = imported;

  history.add({ question });
}

let errorParameters = $ref<ErrParameters | null>(null);
const showErrorParameters = computed(() => errorParameters != null);
const availableParameters = ref<Variable[]>([]);
let isCheckingParameters = $ref(false);

async function checkParameters() {
  isCheckingParameters = true;
  const out = await controller.EditorCheckQuestionParameters({
    SessionID: props.session_id || "",
    Parameters: question.Page.parameters,
  });
  isCheckingParameters = false;
  if (out === undefined) return;

  // hide previous error
  errorEnnonce = null;

  errorParameters = out.ErrDefinition.Origin == "" ? null : out.ErrDefinition;
  availableParameters.value = out.Variables || [];
}

let showEditDescription = $ref(false);
</script>

<style scoped>
.input-small:deep(input) {
  font-size: 14px;
  width: 100%;
}
</style>
