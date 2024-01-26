<template>
  <v-row style="height: 90dvh">
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
    <v-col>
      <v-card>
        <v-row no-gutters class="pl-2">
          <v-col cols="auto" align-self="center">
            <v-btn
              size="small"
              icon
              title="Retour au tableau"
              @click="emit('back')"
              class="mr-1"
            >
              <v-icon icon="mdi-arrow-left"></v-icon>
            </v-btn>
          </v-col>

          <v-col>
            <StageChip :stage="props.stage"></StageChip>
          </v-col>

          <v-spacer></v-spacer>
          <v-col cols="6" align-self="center">
            <v-tabs
              style="max-width: 90vh"
              density="compact"
              show-arrows
              color="grey"
              v-model="questionIndex"
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
                    <!-- <v-list-item
                      @click="
                        props.readonly ? {} : emit('duplicateVariant', variant)
                      "
                      :link="!props.readonly"
                    >
                      <template v-slot:prepend>
                        <v-icon
                          icon="mdi-content-copy"
                          color="info"
                          size="small"
                        ></v-icon>
                      </template>
                      Dupliquer
                    </v-list-item> -->
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

        <BeltQuestionEditor
          class="mt-2"
          v-if="questions.length"
          :question="question"
          :readonly="props.readonly"
          @save="saveQuestion"
        ></BeltQuestionEditor>
      </v-card>
    </v-col>

    <v-col cols="auto">
      <ClientPreview ref="preview" :hide="false"></ClientPreview>
    </v-col>
  </v-row>
</template>

<script setup lang="ts">
import { Beltquestion, IdBeltquestion, Stage } from "@/controller/api_gen";
import ClientPreview from "../editor/ClientPreview.vue";
import { ref } from "vue";
import StageChip from "./StageChip.vue";
import { controller } from "@/controller/controller";
import { onMounted } from "vue";
import { computed } from "vue";
import BeltQuestionEditor from "./BeltQuestionEditor.vue";

interface Props {
  stage: Stage;
  readonly: boolean;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "back"): void;
}>();

onMounted(fetchQuestions);

const preview = ref<InstanceType<typeof ClientPreview> | null>(null);

const questionIndex = ref(0);
const questions = ref<Beltquestion[]>([]);
async function fetchQuestions() {
  const res = await controller.CeinturesGetQuestions(props.stage);
  if (res === undefined) return;
  questions.value = res || [];
}

async function createQuestion() {
  const res = await controller.CeinturesCreateQuestion(props.stage);
  if (res === undefined) return;

  controller.showMessage("Question créée avec succès.");
  questions.value.push(res);
  questionIndex.value = questions.value.length - 1;
}

async function saveQuestion(qu: Beltquestion) {
  const res = await controller.CeinturesSaveQuestion(qu);
  if (res === undefined) return;

  controller.showMessage("Question modifiée avec succès.");
  questions.value[questionIndex.value] = qu;

  preview.value?.showCeinture(res.Preview);
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
</script>
