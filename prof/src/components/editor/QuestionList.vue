<template>
  <v-dialog
    :model-value="questionToDelete != null"
    @update:model-value="questionToDelete = null"
  >
    <v-card title="Confirmer">
      <v-card-text
        >Etes-vous certain de vouloir supprimer la question
        <i>{{ questionToDelete?.Title }}</i> ? <br />
        Cette opération est irréversible.
      </v-card-text>
      <v-card-actions>
        <v-btn @click="questionToDelete = null">Retour</v-btn>
        <v-spacer></v-spacer>
        <v-btn color="red" @click="deleteQuestion" variant="outlined">
          Supprimer
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>

  <v-dialog
    :model-value="questionToDuplicate != null"
    @update:model-value="questionToDuplicate = null"
  >
    <v-card title="Dupliquer avec difficulté" max-width="600">
      <v-card-text
        >Voulez-vous dupliquer la question
        <i>{{ questionToDuplicate?.Title }}</i> en générant les étiquettes des
        niveaux de difficulté manquants ? <br />
      </v-card-text>
      <v-card-actions>
        <v-btn @click="questionToDuplicate = null">Retour</v-btn>
        <v-spacer></v-spacer>
        <v-btn color="green" @click="duplicateDifficulty" variant="outlined">
          Dupliquer
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>

  <v-card class="pt-2 pb-0">
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

    <v-card-text class="py-1">
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
            clearable
          ></v-text-field>
        </v-col>
        <v-col>
          <v-autocomplete
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
          ></v-autocomplete>
        </v-col>
      </v-row>
      <v-row no-gutters>
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

            <div v-for="questionGroup in questions" :key="questionGroup.Title">
              <question-row
                v-if="questionGroup.Size == 1"
                :question="questionGroup.Questions![0]"
                @clicked="startEdit"
                @delete="(question) => (questionToDelete = question)"
                @duplicate="(question) => (questionToDuplicate = question)"
                @update-public="updatePublic"
              ></question-row>
              <question-group-row
                v-else
                :group="questionGroup"
                :all-tags="props.tags"
                @clicked="startEdit"
                @delete="(question) => (questionToDelete = question)"
                @update-public="updatePublic"
                @update-tags="(tags) => updateGroupTags(questionGroup, tags)"
              ></question-group-row>
            </div>
          </v-list>
          <div class="my-2">
            {{ questions.length }} / {{ serverNbGroups }} questions affichées.
            ({{ displayedNbQuestions }} / {{ serverNbQuestions }}
            variantes)
          </div>
        </v-col>
      </v-row>
    </v-card-text>
  </v-card>
</template>

<script setup lang="ts">
import type {
  Origin,
  Question,
  QuestionGroup,
  QuestionHeader,
} from "@/controller/api_gen";
import { controller, IsDev } from "@/controller/controller";
import { personnalOrigin } from "@/controller/editor";
import { computed, onActivated, onMounted } from "@vue/runtime-core";
import { $ref } from "vue/macros";
import QuestionGroupRow from "./QuestionGroupRow.vue";
import QuestionRow from "./QuestionRow.vue";

interface Props {
  tags: string[];
}
const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "edit", question: Question, tags: string[], origin: Origin): void;
}>();

let questions = $ref<QuestionGroup[]>([]);
let serverNbGroups = $ref(0);
let serverNbQuestions = $ref(0);
const displayedNbQuestions = computed(() => {
  let nb = 0;
  questions.forEach((group) => {
    nb += group.Questions?.length || 0;
  });
  return nb;
});

let querySearch = $ref("");

let queryTags = $ref<string[]>(IsDev ? ["DEV"] : []);

let timerId = 0;

onMounted(fetchQuestions);
onActivated(fetchQuestions);

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
    Tags: queryTags,
  });
  if (result == undefined) {
    return;
  }
  questions = result.Questions || [];
  serverNbGroups = result.NbGroups;
  serverNbQuestions = result.NbQuestions;
}

async function createQuestion() {
  const out = await controller.EditorCreateQuestion();
  if (out == undefined) {
    return;
  }
  emit("edit", out, [], personnalOrigin());
}

let questionToDuplicate: QuestionHeader | null = $ref(null);
async function duplicateDifficulty() {
  await controller.EditorDuplicateQuestionWithDifficulty({
    id: questionToDuplicate!.Id,
  });
  questionToDuplicate = null;
  await fetchQuestions();
}

async function startEdit(question: QuestionHeader) {
  const out = await controller.EditorGetQuestion({ id: question.Id });
  if (out == undefined) {
    return;
  }
  emit("edit", out, question.Tags || [], question.Origin);
}

let questionToDelete: QuestionHeader | null = $ref(null);
async function deleteQuestion() {
  await controller.EditorDeleteQuestion({ id: questionToDelete!.Id });
  await fetchQuestions(); // delete modify the groups
  questionToDelete = null;
}

async function updatePublic(questionID: number, isPublic: boolean) {
  const res = await controller.QuestionUpdateVisiblity({
    QuestionID: questionID,
    Public: isPublic,
  });
  if (res === undefined) {
    return;
  }
  questions.forEach((group) => {
    const index = group.Questions?.findIndex((qu) => qu.Id == questionID);
    if (index !== undefined) {
      group.Questions![index].Origin.IsPublic = isPublic;
    }
  });
}

async function updateGroupTags(group: QuestionGroup, newTags: string[]) {
  const res = await controller.EditorUpdateGroupTags({
    GroupTitle: group.Title,
    CommonTags: newTags,
  });
  if (res == undefined) {
    return;
  }

  group.Questions?.forEach((qu) => (qu.Tags = (res.Tags || {})[qu.Id] || []));
}
</script>
