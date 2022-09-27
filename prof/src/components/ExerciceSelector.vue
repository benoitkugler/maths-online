<template>
  <v-card class="mt-3 px-2" title="Choisir un exercice">
    <v-card-text>
      <v-row>
        <v-col>
          <v-text-field
            label="Rechercher"
            hint="Rechercher un exercice par nom."
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

      <div style="height: 47vh; width: 800px" class="overflow-y-auto">
        <v-expansion-panels class="pa-2">
          <v-expansion-panel v-for="(group, index) in exercices" :key="index">
            <v-expansion-panel-title>
              {{ group.Group.Title }}
            </v-expansion-panel-title>
            <v-expansion-panel-text>
              <v-list>
                <v-list-item
                  v-for="(exercice, index) in group.Variants"
                  :key="index"
                  @click="emit('selected', exercice)"
                >
                  <v-row>
                    <v-col>
                      <small>({{ exercice.Id }})</small> {{ exercice.Subtitle }}
                    </v-col>
                    <v-spacer></v-spacer>
                    <v-col cols="auto">
                      <tag-chip
                        :tag="exercice.Difficulty || 'Aucune difficulté'"
                      ></tag-chip>
                    </v-col>
                  </v-row>
                </v-list-item>
              </v-list>
            </v-expansion-panel-text>
          </v-expansion-panel>
        </v-expansion-panels>
      </div>
      <div class="my-2">
        {{ exercices.length }} / {{ serverNbExercices }} variantes d'exercices
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
import type { ExercicegroupExt, ExerciceHeader } from "../controller/api_gen";
import { controller } from "../controller/controller";
import TagChip from "./editor/utils/TagChip.vue";

export interface ExerciceQuery {
  search: string;
  tags: string[];
}

interface Props {
  tags: string[];
  query: ExerciceQuery;
}

const props = defineProps<Props>();
const emit = defineEmits<{
  (e: "closed"): void;
  (e: "selected", question: ExerciceHeader): void;
  (e: "update:query", query: ExerciceQuery): void;
}>();

let exercices = $ref<ExercicegroupExt[]>([]);
let serverNbExercices = $ref(0);

let timerId = 0;

onMounted(() => {
  if (props.query.search || props.query.tags.length) {
    fetchExercices();
  }
});

function updateQuerySearch() {
  const debounceDelay = 200;
  // cancel pending call
  clearTimeout(timerId);

  // delay new call 500ms
  timerId = setTimeout(() => {
    fetchExercices();
  }, debounceDelay);
}

async function updateQueryTags() {
  await fetchExercices();
}

async function fetchExercices() {
  const result = await controller.EditorSearchExercices({
    TitleQuery: props.query.search,
    Tags: props.query.tags,
  });
  if (result == undefined) return;
  exercices = result.Groups || [];
  serverNbExercices = result.NbExercices;
}
</script>

<style></style>
