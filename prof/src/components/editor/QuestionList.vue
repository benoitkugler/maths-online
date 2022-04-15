<template>
  <v-dialog
    :model-value="questionToDelete != null"
    @update:model-value="questionToDelete = null"
  >
    <v-card title="Confirmer">
      <v-card-text
        >Etes-vous certain de vouloir supprimer la question
        {{ questionToDelete?.Title }} ? <br />
        Cette opération est irréversible.
      </v-card-text>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn color="red" @click="deleteQuestion"> Supprimer </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
  <v-card class="pt-2">
    <v-row>
      <v-col> <v-card-title>Liste des questions</v-card-title> </v-col>

      <v-col align-self="center" style="text-align: right" cols="4">
        <v-btn
          class="mx-2"
          @click="createQuestion"
          title="Créer une nouvelle question"
        >
          <v-icon icon="mdi-plus" color="success"></v-icon>
          Créer
        </v-btn>
      </v-col>
    </v-row>

    <v-card-text>
      <v-row>
        <v-col>
          <v-text-field
            label="Rechercher"
            hint="Rechercher une question par nom."
            variant="outlined"
            density="comfortable"
            v-model="querySearch"
            @update:model-value="updateQuerySearch"
            persistent-hint
          ></v-text-field>
        </v-col>
        <v-col>
          <v-select
            variant="outlined"
            density="comfortable"
            multiple
            chips
            closable-chips
            :items="props.tags"
            color="primary"
            label="Catégories"
            no-data-text="Aucune catégorie n'est encore utilisée."
            v-model="queryTags"
            @update:model-value="updateQueryTags"
            @blur="updateQueryTags"
            hint="Restreint la recherche à l'intersection des catégories sélectionnées."
            persistent-hint
          ></v-select>
        </v-col>
      </v-row>
      <v-row>
        <v-col>
          <v-list style="height: 52vh" class="overflow-y-auto">
            <div
              v-if="questions.length == 0"
              style="width: 100%; text-align: center"
            >
              <i>
                {{
                  querySearch == "" && queryTags.length == 0
                    ? "Entrer une recherche..."
                    : "Aucun résultat"
                }}
              </i>
            </div>
            <v-list-item
              dense
              class="py-0"
              v-for="question in questions"
              @click="startEdit(question)"
            >
              <v-list-item-media class="mr-3">
                <v-btn
                  size="x-small"
                  icon
                  @click.stop="questionToDelete = question"
                >
                  <v-icon icon="mdi-delete" color="red"></v-icon>
                </v-btn>
              </v-list-item-media>
              <v-list-item-title>
                {{ question.Title ? question.Title : "..." }}
              </v-list-item-title>
              <v-spacer></v-spacer>
              <v-list-item-media>
                <v-chip
                  v-for="tag in question.Tags"
                  :key="tag"
                  size="small"
                  label
                  class="ma-1"
                  color="primary"
                  >{{ tag.toUpperCase() }}</v-chip
                >
              </v-list-item-media>
            </v-list-item>
          </v-list>
        </v-col>
      </v-row>
    </v-card-text>
  </v-card>
</template>

<script setup lang="ts">
import type { QuestionHeader } from "@/controller/api_gen";
import { controller } from "@/controller/controller";
import type { Question } from "@/controller/exercice_gen";
import { onMounted } from "@vue/runtime-core";
import { $ref } from "vue/macros";

interface Props {
  tags: string[];
}
const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "edit", question: Question, tags: string[]): void;
}>();

let questions = $ref<QuestionHeader[]>([]);

let querySearch = $ref("");
let queryTags = $ref<string[]>([]);

let timerId = 0;

onMounted(fetchQuestions);

function updateQuerySearch() {
  const debounceDelay = 200;
  // cancel pending call
  clearTimeout(timerId);

  // delay new call 500ms
  timerId = setTimeout(() => {
    fetchQuestions();
  }, debounceDelay);
}

async function updateQueryTags() {
  await fetchQuestions();
}

async function fetchQuestions() {
  const result = await controller.EditorSearchQuestions({
    TitleQuery: querySearch,
    Tags: queryTags
  });
  questions = result || [];
}

async function createQuestion() {
  const out = await controller.EditorCreateQuestion({});
  if (out == undefined) {
    return;
  }
  emit("edit", out, []);
}

async function startEdit(question: QuestionHeader) {
  const out = await controller.EditorGetQuestion({ id: String(question.Id) });
  if (out == undefined) {
    return;
  }
  emit("edit", out, question.Tags || []);
}

let questionToDelete: QuestionHeader | null = $ref(null);
async function deleteQuestion() {
  await controller.EditorDeleteQuestion({ id: String(questionToDelete?.Id) });
  questions = questions.filter(qu => qu.Id != questionToDelete?.Id);
  questionToDelete = null;
}
</script>
