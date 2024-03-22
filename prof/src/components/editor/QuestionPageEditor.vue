<template>
  <v-row no-gutters>
    <v-col>
      <v-card class="mt-1 px-2">
        <v-row no-gutters>
          <v-col md="5" align-self="center">
            <slot name="top-left"></slot>
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
                props.readonly ? 'Visualiser' : 'Enregistrer et prévisualiser'
              "
              size="small"
              :disabled="!question"
            >
              <v-badge color="pink" dot v-model="isDirty">
                <v-icon
                  :icon="props.readonly ? 'mdi-eye' : 'mdi-content-save'"
                ></v-icon>
              </v-badge>
            </v-btn>

            <v-menu offset-y close-on-content-click>
              <template v-slot:activator="{ isActive, props }">
                <v-btn
                  flat
                  icon
                  title="Plus d'options"
                  v-on="{ isActive }"
                  v-bind="props"
                  size="small"
                >
                  <v-icon icon="mdi-dots-vertical"></v-icon>
                </v-btn>
              </template>
              <v-list>
                <v-list-item
                  @click="paste"
                  title="Coller"
                  subtitle="un bloc"
                  prepend-icon="mdi-content-paste"
                >
                </v-list-item>

                <v-divider></v-divider>
                <v-list-item
                  @click="exportJSON"
                  title="Exporter"
                  subtitle="au format isyro.json"
                  prepend-icon="mdi-download"
                >
                </v-list-item>

                <v-list-item
                  @click="showImportJSON = true"
                  title="Importer"
                  prepend-icon="mdi-upload"
                  subtitle="depuis un fichier isyro.json"
                >
                </v-list-item>

                <v-divider></v-divider>

                <v-list-item
                  @click="exportLatex"
                  prepend-icon="mdi-file-export"
                  title="Exporter en LaTeX"
                  subtitle="(Expérimental)"
                >
                </v-list-item>
              </v-list>
            </v-menu>
          </v-col>
        </v-row>

        <v-row no-gutters>
          <v-col md="5" v-if="!hasEditorSimplified">
            <ParametersEditor
              :dual="props.showDualParameters"
              :parameters="inner.parameters"
              :shared-parameters="inner.sharedParameters"
              @update="checkParameters"
              :is-loading="isCheckingParameters"
              :is-validated="!showErrorParameters"
            ></ParametersEditor>
          </v-col>
          <v-col class="pr-1">
            <QuestionContent
              v-if="modeEnonce"
              :model-value="inner.enonce || []"
              @update:model-value="onUpdateEnonce"
              @import-question="doImportJSON"
              @add-syntax-hint="addSyntaxHint"
              :available-parameters="[]"
              :errorBlockIndex="
                errorIsCorrection ? undefined : errorContent?.Block
              "
              ref="questionEnonceNode"
            >
            </QuestionContent>
            <QuestionContent
              v-else
              :model-value="inner.correction || []"
              @update:model-value="onUpdateCorrection"
              @import-question="doImportJSON"
              :available-parameters="[]"
              :errorBlockIndex="
                errorIsCorrection ? errorContent?.Block : undefined
              "
              ref="questionCorrectionNode"
            >
            </QuestionContent>
          </v-col>
        </v-row>
      </v-card>
    </v-col>
    <v-col cols="auto">
      <ClientPreview ref="preview"></ClientPreview>
    </v-col>
  </v-row>

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

  <v-dialog v-model="showImportJSON" max-width="800">
    <v-card
      title="Importer une question"
      subtitle="Importer un fichier généré depuis l'éditeur, au format .isyro.json"
    >
      <v-card-text>
        <v-file-input v-model="importedFiles" label="Fichier" accept=".json">
        </v-file-input>
      </v-card-text>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn :disabled="!importedFiles.length" @click="onImportJSON"
          >Importer</v-btn
        >
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script setup lang="ts">
import {
  BlockKind,
  Enonce,
  ErrParameters,
  ErrQuestionInvalid,
  ErrorKind,
  ExportQuestionLatexOut,
  ExpressionFieldBlock,
  Parameters,
  Variable,
  errEnonce,
} from "@/controller/api_gen";
import BlockBar from "./BlockBar.vue";
import ClientPreview from "./ClientPreview.vue";
import QuestionContent from "./QuestionContent.vue";
import SnackErrorEnonce from "./SnackErrorEnonce.vue";
import ParametersEditor from "./parameters/ParametersEditor.vue";
import SnackErrorParameters from "./parameters/SnackErrorParameters.vue";
import {
  QuestionPage,
  SaveQuestionOut,
  readClipboardForBlock,
  saveData,
} from "@/controller/editor";
import { computed, onMounted, onUnmounted, ref } from "vue";
import { History } from "@/controller/editor_history";
import { controller } from "@/controller/controller";
import { copy } from "@/controller/utils";
import { watch } from "vue";

const props = defineProps<{
  question: QuestionPage;
  onSave: (showCorrection: boolean) => Promise<SaveQuestionOut | undefined>;
  onExportLatex: () => Promise<ExportQuestionLatexOut | undefined>;
  readonly: boolean;
  showDualParameters: boolean;
}>();

const emit = defineEmits<{
  (e: "update", page: QuestionPage): void;
}>();

watch(
  () => props.question,
  (newV, oldV) => {
    inner.value = copy(props.question);
    if (newV.id != oldV.id) {
      isDirty.value = false; // reset since it is unknown
    }
  }
);

const inner = ref(copy(props.question));
// updated on saved, used to implement dirty feature
const isDirty = ref(false);

const modeEnonce = ref(true); // false for correction

let history = new History(inner.value, controller.showMessage, restoreHistory);

onMounted(() => {
  history.addListener();
});
onUnmounted(() => {
  history.clearListener();
});

const hasEditorSimplified = computed(
  () => controller.settings.HasEditorSimplified
);

function update() {
  emit("update", inner.value);
  isDirty.value = true;
}

function restoreHistory(snapshot: QuestionPage) {
  inner.value = snapshot;
  update();
}

const preview = ref<InstanceType<typeof ClientPreview> | null>(null);

const questionEnonceNode = ref<InstanceType<typeof QuestionContent> | null>(
  null
);
const questionCorrectionNode = ref<InstanceType<typeof QuestionContent> | null>(
  null
);
function addBlock(kind: BlockKind) {
  // this triggers an update event
  if (modeEnonce.value) {
    questionEnonceNode.value?.addBlock(kind);
  } else {
    questionCorrectionNode.value?.addBlock(kind);
  }
}

function onUpdateEnonce(v: Enonce) {
  inner.value.enonce = v;
  history.add(copy(inner.value));
  update();
}
function onUpdateCorrection(v: Enonce) {
  inner.value.correction = v;
  history.add(copy(inner.value));
  update();
}

const availableParameters = ref<Variable[]>([]);
const isCheckingParameters = ref(false);
const errorParameters = ref<ErrParameters | null>(null);
const showErrorParameters = computed(() => errorParameters.value != null);

const errorContent = ref<errEnonce | null>(null);
const errorIsCorrection = ref(false);

async function checkParameters(ps: Parameters, shared: Parameters) {
  inner.value.parameters = ps;
  inner.value.sharedParameters = shared;
  history.add(copy(inner.value));
  update();

  isCheckingParameters.value = true;
  const out = await controller.EditorCheckQuestionParameters({
    Parameters: (ps || []).concat(shared || []),
  });
  isCheckingParameters.value = false;
  if (out === undefined) return;

  // hide previous error
  errorContent.value = null;

  errorParameters.value =
    out.ErrDefinition.Origin == "" ? null : out.ErrDefinition;
  availableParameters.value = out.Variables || [];
}

async function save() {
  const res = await props.onSave(!modeEnonce.value);
  if (res == undefined) return;

  if (res.IsValid) {
    isDirty.value = false;

    controller.showMessage(
      props.readonly
        ? "Question générée avec succès."
        : "Question enregistrée avec succès."
    );
    errorContent.value = null;
    errorParameters.value = null;

    preview.value?.preview(res.Preview); // TODO:
  } else {
    onQuestionError(res.Error);
  }
}

function onQuestionError(err: ErrQuestionInvalid) {
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
}

function exportJSON() {
  saveData(inner.value, `question.isyro.json`);
}

function doImportJSON(json: string) {
  const imported: QuestionPage = JSON.parse(json);
  // do not erase the id
  imported.id = props.question.id;
  inner.value = imported;
  history.add(copy(inner.value));
  update();
}

const showImportJSON = ref(false);
const importedFiles = ref<File[]>([]);
async function onImportJSON() {
  showImportJSON.value = false;
  if (!importedFiles.value.length) return;

  const file = importedFiles.value[0];
  const content = await file.text();
  doImportJSON(content);
}

async function addSyntaxHint(block: ExpressionFieldBlock) {
  if (questionEnonceNode.value == null) return;

  const res = await controller.EditorGenerateSyntaxHint({
    Block: block,
    QuestionParameters: inner.value.parameters,
    SharedParameters: inner.value.sharedParameters,
  });
  if (res == undefined) return;

  questionEnonceNode.value?.addExistingBlock({
    Kind: BlockKind.TextBlock,
    Data: res,
  });
}

async function exportLatex() {
  const res = await props.onExportLatex();
  if (res == undefined) return;

  if (res.IsValid) {
    try {
      await navigator.clipboard.writeText(res.Latex);
      controller.showMessage("Exercice copié dans le presse-papier");
    } catch (error) {
      controller.onError("Presse-papier", "L'accès au presse-papier a échoué.");
    }
  } else {
    onQuestionError(res.Error);
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
