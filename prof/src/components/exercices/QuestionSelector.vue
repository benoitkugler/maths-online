<template>
  <v-card class="mt-3 px-2" title="Choisir une question">
    <v-card-text>
      <v-row>
        <v-col>
          <v-text-field
            label="Rechercher"
            hint="Rechercher une question par nom."
            variant="outlined"
            density="comfortable"
            v-model="props.query.search"
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
            v-model="props.query.tags"
            @update:model-value="updateQueryTags"
            @blur="updateQueryTags"
            hint="Restreint la recherche à l'intersection des catégories sélectionnées."
            persistent-hint
          ></v-autocomplete>
        </v-col>
      </v-row>

      <v-list style="height: 47vh" class="overflow-y-auto">
        <v-list-item
          v-for="(question, index) in questions"
          :key="index"
          @click="emit('selected', question)"
        >
          <v-row>
            <v-col>
              <small>({{ question.Id }})</small> {{ question.Title }}
            </v-col>
            <v-col cols="5">
              <tag-list-field
                readonly
                :model-value="question.Tags || []"
                :all-tags="[]"
              ></tag-list-field>
            </v-col>
          </v-row>
        </v-list-item>
      </v-list>
      <div class="my-2">
        {{ questions.length }} / {{ serverNbQuestions }} variantes de questions
        affichées
      </div>
    </v-card-text>
    <v-card-actions>
      <v-btn @click="emit('closed')">Retour</v-btn>
      <v-spacer></v-spacer>
    </v-card-actions>
  </v-card>
</template>

<script setup lang="ts">
import { onMounted } from "vue";
import { $ref } from "vue/macros";
import type { QuestionHeader } from "../../controller/api_gen";
import { controller } from "../../controller/controller";
import TagListField from "../editor/TagListField.vue";

export interface Query {
  search: string;
  tags: string[];
}

interface Props {
  tags: string[];
  query: Query;
}

const props = defineProps<Props>();
const emit = defineEmits<{
  (e: "closed"): void;
  (e: "selected", question: QuestionHeader): void;
  (e: "update:query", query: Query): void;
}>();

let questions = $ref<QuestionHeader[]>([]);
let serverNbQuestions = $ref(0);

let timerId = 0;

onMounted(() => {
  if (props.query.search || props.query.tags.length) {
    fetchQuestions();
  }
});

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
    TitleQuery: props.query.search,
    Tags: props.query.tags,
  });
  if (result == undefined) {
    return;
  }
  const qus: QuestionHeader[] = [];
  (result.Questions || []).forEach((gr) => qus.push(...(gr.Questions || [])));
  questions = qus;
  serverNbQuestions = result.NbQuestions;
}
</script>

<style></style>
