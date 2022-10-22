<template>
  <v-card max-height="80vh" class="overflow-y-auto py-2">
    <v-row>
      <v-col>
        <v-card-title>Modifier la session de Triv'Maths</v-card-title>
      </v-col>
      <v-col style="text-align: right">
        <v-btn icon flat class="mx-2" @click="emit('closed')">
          <v-icon icon="mdi-close"></v-icon>
        </v-btn>
      </v-col>
    </v-row>
    <v-card-text>
      <v-row>
        <v-col>
          <v-text-field
            label="Nom"
            density="compact"
            variant="outlined"
            v-model="props.edited.Name"
            hint="Usage interne, non visible par les élèves."
          >
          </v-text-field>
        </v-col>
      </v-row>
      <v-row>
        <v-col>
          <v-list-subheader>
            <h3>Choix des questions</h3>
            <small
              >Chaque catégorie est définie par une <i>union</i> d'<i
                >intersections</i
              >
              d'étiquettes.</small
            >
          </v-list-subheader>
        </v-col>
        <v-col cols="auto" align-self="center">
          <v-menu
            offset-y
            :close-on-content-click="false"
            v-model="showDifficultyCard"
          >
            <template v-slot:activator="{ isActive, props }">
              <v-btn
                class="mr-1"
                title="Sélectionner la difficulté"
                v-on="{ isActive }"
                v-bind="props"
                size="small"
              >
                Difficulté
              </v-btn>
            </template>
            <DifficultyChoices
              :difficulties="props.edited.Questions.Difficulties || []"
              @update:difficulties="onEditDifficulties"
            ></DifficultyChoices>
          </v-menu>

          <v-menu offset-y close-on-content-click>
            <template v-slot:activator="{ isActive, props }">
              <v-btn
                title="Utiliser une configuration de référence"
                v-on="{ isActive }"
                v-bind="props"
                size="small"
              >
                <v-icon icon="mdi-plus" color="green"></v-icon>
                Importer un modèle
              </v-btn>
            </template>
            <v-list>
              <v-list-item
                v-for="prop in propositions"
                :key="prop.name"
                rounded
                class="py-0"
                @click="importQuestions(prop.Questions)"
              >
                {{ prop.name }}
              </v-list-item>
            </v-list>
          </v-menu>
        </v-col>
      </v-row>
      <v-list>
        <v-alert
          class="mt-1 py-2 px-3"
          variant="outlined"
          v-if="hint.Pattern?.length"
          :color="hint.Missing?.length ? 'info' : 'success'"
          closable
        >
          <div
            v-if="hint.Missing?.length"
            style="max-height: 50vh"
            class="overflow-y-auto mr-2"
          >
            Les catégories suivantes ne sont pas utilisées :
            <v-list>
              <v-list-item v-for="(tags, index) in hint.Missing!" :key="index">
                <TagChip
                  :tag="tag"
                  :key="tag"
                  v-for="tag in tags"
                  :pointer="false"
                ></TagChip>
              </v-list-item>
            </v-list>
          </div>
          <div v-else>
            Toutes les questions inclusent dans
            <span>
              <TagChip
                :tag="tag"
                :key="tag"
                v-for="tag in hint.Pattern"
                :pointer="false"
              ></TagChip
            ></span>
            sont utilisées.
          </div>
        </v-alert>
        <v-list-item
          v-for="(categorie, index) in props.edited.Questions.Tags"
          :key="index"
          rounded
          :style="{
            'border-color': colors[index],
            borderWidth: '2px',
            borderStyle: 'solid',
          }"
          class="my-2"
        >
          <v-row>
            <v-col align-self="center" cols="2">
              <v-list-item-subtitle>
                Catégorie {{ index + 1 }}
              </v-list-item-subtitle>
            </v-col>
            <v-col class="my-1">
              <tags-selector
                :all-tags="allKnownTags"
                :model-value="categorie || []"
                @update:model-value="(v) => updateCategorie(index, v)"
              ></tags-selector>
            </v-col>
          </v-row>
        </v-list-item>
      </v-list>

      <v-row class="mt-2">
        <v-col cols="6">
          <v-text-field
            density="compact"
            variant="outlined"
            label="Durée limite pour une question"
            type="number"
            min="1"
            suffix="sec"
            v-model.number="props.edited.QuestionTimeout"
          ></v-text-field>
        </v-col>
        <v-col cols="6">
          <v-checkbox
            density="compact"
            label="Afficher le décrassage en fin de partie"
            v-model.number="props.edited.ShowDecrassage"
          ></v-checkbox>
        </v-col>
      </v-row>
    </v-card-text>

    <v-card-actions>
      <v-spacer></v-spacer>
      <v-btn color="success" @click="emit('update', props.edited)">
        Enregistrer les modifications
      </v-btn>
    </v-card-actions>
  </v-card>
</template>

<script setup lang="ts">
import type {
  CategoriesQuestions,
  CheckMissingQuestionsOut,
  DifficultyTag,
  QuestionCriterion,
  Trivial,
} from "@/controller/api_gen";
import { controller } from "@/controller/controller";
import { colorsPerCategorie, questionPropositions } from "@/controller/trivial";
import { onMounted } from "@vue/runtime-core";
import { $ref } from "vue/macros";
import TagChip from "../editor/utils/TagChip.vue";
import DifficultyChoices from "./DifficultyChoices.vue";
import TagsSelector from "./TagsSelector.vue";

interface Props {
  edited: Trivial;
  allKnownTags: string[];
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "closed"): void;
  (e: "update", v: Trivial): void;
}>();

const colors = colorsPerCategorie;

const propositions = questionPropositions;

let showDifficultyCard = $ref(false);

function onEditDifficulties(difficulties: DifficultyTag[]) {
  showDifficultyCard = false;
  props.edited.Questions.Difficulties = difficulties;
  fetchHintForMissing();
}

function updateCategorie(index: number, cat: QuestionCriterion) {
  props.edited.Questions.Tags[index] = cat;
  fetchHintForMissing();
}

function importQuestions(criteria: CategoriesQuestions) {
  props.edited.Questions = criteria;
  fetchHintForMissing();
}

onMounted(fetchHintForMissing);

let hint = $ref<CheckMissingQuestionsOut>({ Pattern: [], Missing: [] });
async function fetchHintForMissing() {
  const criteria = props.edited.Questions;
  // fetch the hint only if the all categories have been filled,
  // to avoid useless queries
  if (!criteria.Tags.every((qu) => qu?.length)) {
    hint = { Pattern: [], Missing: [] };
    return;
  }
  const res = await controller.CheckMissingQuestions(criteria);
  if (res == undefined) {
    return;
  }
  hint = res;
}
</script>

<style scoped></style>
