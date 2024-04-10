<template>
  <v-card width="100%">
    <!-- Suppression confirmation -->
    <v-dialog
      :model-value="toDelete != null"
      @update:model-value="toDelete = null"
      width="600px"
    >
      <v-card title="Confirmer la suppression">
        <v-card-text>
          Confirmez-vous la suppression de la question ? <br /><br />

          Cette opération est irréversible.
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn color="red" @click="deleteQuestion">Supprimer</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- Question settings -->
    <v-dialog
      :model-value="settingsToEdit != null"
      @update:model-value="settingsToEdit = null"
      width="600px"
    >
      <v-card title="Réglages de la question" v-if="settingsToEdit != null">
        <v-card-text>
          <v-row
            ><v-col>
              <v-select
                label="Nombre de répétitions"
                variant="outlined"
                density="compact"
                v-model="settingsToEdit.Repeat"
                :items="[
                  { title: '1', value: 1 },
                  { title: '2', value: 2 },
                  { title: '3', value: 3 },
                  { title: '4', value: 4 },
                  { title: '5', value: 5 },
                ]"
              >
              </v-select> </v-col
          ></v-row>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn color="success" @click="updateQuestion">Enregistrer</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <v-card-text class="py-1 px-2">
      <v-row no-gutters>
        <v-col cols="auto" align-self="center">
          <v-btn
            size="small"
            icon
            title="Retour à la vue d'ensemble"
            @click="emit('back')"
            class="mr-1"
          >
            <v-icon icon="mdi-arrow-left"></v-icon>
          </v-btn>
        </v-col>

        <v-spacer></v-spacer>

        <v-col cols="auto">
          <v-btn
            icon
            flat
            size="small"
            :disabled="props.stage.Rank == Rank.Blanche"
            @click="emit('goTo', (props.stage.Rank - 1) as Rank)"
            ><v-icon>mdi-chevron-left</v-icon></v-btn
          >
          <v-menu>
            <template v-slot:activator="{ isActive, props: innerProps }">
              <v-chip
                v-on="{ isActive }"
                v-bind="innerProps"
                variant="elevated"
                :color="rankColors[props.stage.Rank]"
                class="ma-2"
                label
              >
                {{ DomainLabels[props.stage.Domain] }}
              </v-chip>
            </template>
            <v-list>
              <v-list-item
                v-for="r in rankItems"
                :key="r"
                :title="RankLabels[r]"
                @click="emit('goTo', r)"
              >
                <template v-slot:prepend>
                  <RankIcon :rank="r" small></RankIcon>
                </template>
              </v-list-item>
            </v-list>
          </v-menu>
          <v-btn
            icon
            flat
            size="small"
            :disabled="props.stage.Rank == Rank.Noire"
            @click="emit('goTo', (props.stage.Rank + 1) as Rank)"
            ><v-icon>mdi-chevron-right</v-icon></v-btn
          >
        </v-col>

        <v-spacer></v-spacer>
        <v-col cols="6" align-self="center">
          <v-tabs
            density="compact"
            show-arrows
            color="grey"
            v-model="questionIndex"
            @update:model-value="refreshPreview"
            align-tabs="end"
          >
            <v-tab
              v-for="(question, index) in questions"
              :key="index"
              class="text-subtitle-2 font-weight-light"
            >
              Question {{ index + 1 }}

              <v-menu
                offset-y
                close-on-content-click
                v-if="index == questionIndex"
              >
                <template v-slot:activator="{ isActive, props: innerProps2 }">
                  <v-btn
                    v-on="{ isActive }"
                    v-bind="innerProps2"
                    icon
                    size="x-small"
                    flat
                    class="pr-0 mr-0"
                  >
                    <v-icon>mdi-dots-vertical</v-icon>
                  </v-btn>
                </template>
                <v-list>
                  <v-list-item
                    @click="
                      props.readonly ? {} : (settingsToEdit = copy(question))
                    "
                    :link="!props.readonly"
                  >
                    <template v-slot:prepend>
                      <v-icon icon="mdi-cog" color="info" size="small"></v-icon>
                    </template>
                    Réglages
                  </v-list-item>
                  <v-divider></v-divider>
                  <v-list-item
                    @click="props.readonly ? {} : (toDelete = question.Id)"
                    :link="!props.readonly"
                  >
                    <template v-slot:prepend>
                      <v-icon
                        icon="mdi-delete"
                        color="red"
                        size="small"
                      ></v-icon>
                    </template>

                    Supprimer
                  </v-list-item>
                </v-list>
              </v-menu>
            </v-tab>
          </v-tabs>
        </v-col>
        <v-col cols="auto" align-self="center">
          <v-btn icon size="small" class="mx-2" @click="createQuestion">
            <v-icon color="green">mdi-plus</v-icon>
          </v-btn>
        </v-col>
      </v-row>

      <div v-if="!questions.length" class="text-center my-4">
        <v-btn @click="createQuestion">
          <template v-slot:prepend>
            <v-icon color="green">mdi-plus</v-icon>
          </template>
          Créer une question</v-btn
        >
      </div>

      <QuestionPageEditor
        v-if="questions.length"
        :question="page"
        :readonly="props.readonly"
        :show-dual-parameters="false"
        @update="writeChanges"
        @save="saveQuestion"
        @export-latex="exportLatex"
        ref="editor"
      ></QuestionPageEditor>
    </v-card-text>
  </v-card>
</template>

<script setup lang="ts">
import {
  Beltquestion,
  IdBeltquestion,
  Rank,
  RankLabels,
  DomainLabels,
  Stage,
  LoopbackShowCeinture,
  Int,
} from "@/controller/api_gen";
import { ref, watch } from "vue";
import { controller } from "@/controller/controller";
import { onMounted } from "vue";
import { computed } from "vue";
import { copy, rankColors } from "@/controller/utils";
import RankIcon from "./RankIcon.vue";
import QuestionPageEditor from "../editor/QuestionPageEditor.vue";
import { QuestionPage, SaveQuestionOut } from "@/controller/editor";
import { LoopbackServerEventKind } from "@/controller/loopback_gen";

interface Props {
  stage: Stage;
  readonly: boolean;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "back"): void;
  (e: "goTo", rank: Rank): void;
}>();

onMounted(fetchQuestions);

watch(() => props.stage, fetchQuestions);

const rankItems = Object.keys(RankLabels)
  .map((r) => Number(r) as Rank)
  .filter((r) => r != Rank.StartRank);

const questionIndex = ref(0);
const questions = ref<Beltquestion[]>([]);
async function fetchQuestions() {
  const res = await controller.CeinturesGetQuestions(props.stage);
  if (res === undefined) return;
  // reset questionIndex..
  questionIndex.value = 0;
  // .. and preview
  editor.value?.updatePreview({
    Kind: LoopbackServerEventKind.LoopbackPaused,
    Data: {},
  });
  questions.value = res || [];
}

async function createQuestion() {
  const res = await controller.CeinturesCreateQuestion(props.stage);
  if (res === undefined) return;

  controller.showMessage("Question créée avec succès.");
  questions.value.push(res);
  questionIndex.value = questions.value.length - 1;
}

const toDelete = ref<IdBeltquestion | null>(null);
async function deleteQuestion() {
  const id = toDelete.value;
  if (id == null) return;
  const res = await controller.CeinturesDeleteQuestion({ id: id });
  if (res === undefined) return;

  controller.showMessage("Question supprimée avec succès.");
  questions.value = questions.value.filter((qu) => qu.Id != id);
  questionIndex.value = 0;
  toDelete.value = null;
}

const question = computed(() => questions.value[questionIndex.value]);
const page = computed<QuestionPage>(() => ({
  id: question.value.Id,
  parameters: question.value.Parameters,
  sharedParameters: [],
  enonce: question.value.Enonce,
  correction: question.value.Correction,
}));

async function writeChanges(qu: QuestionPage) {
  const question = questions.value[questionIndex.value];
  question.Parameters = qu.parameters;
  question.Enonce = qu.enonce;
  question.Correction = qu.correction;
}

async function saveQuestion(
  isCorrection: boolean
): Promise<SaveQuestionOut | undefined> {
  const res = await controller.CeinturesSaveQuestion({
    Question: question.value,
    ShowCorrection: isCorrection,
  });
  if (res === undefined) return;
  return {
    IsValid: res.IsValid,
    Error: res.Error,
    Preview: {
      Kind: LoopbackServerEventKind.LoopbackShowCeinture,
      Data: res.Preview,
    },
  };
}

async function exportLatex() {
  const res = await controller.EditorQuestionExportLateX(page.value);
  return res;
}

const settingsToEdit = ref<Beltquestion | null>(null);
async function updateQuestion() {
  const qu = settingsToEdit.value;
  if (qu == null) return;
  settingsToEdit.value = null;
  const res = await controller.CeinturesUpdateQuestion({
    Id: qu.Id,
    Repeat: qu.Repeat,
  });
  if (res === undefined) return;
  questions.value.find((q) => q.Id == qu.Id)!.Repeat = qu.Repeat;
  controller.showMessage("Réglages modifiés avec succès.");
}

const editor = ref<InstanceType<typeof QuestionPageEditor> | null>(null);

function refreshPreview() {
  const data = editor.value?.previewData();

  if (
    data === undefined ||
    data.Kind != LoopbackServerEventKind.LoopbackShowCeinture
  ) {
    return;
  }
  const d = data.Data as LoopbackShowCeinture;
  d.QuestionIndex = questionIndex.value as Int;
  data.Data = d;
  editor.value?.updatePreview(data);
}
</script>
