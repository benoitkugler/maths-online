<template>
  <v-card max-height="80vh" width="800" class="overflow-y-auto py-2">
    <v-row>
      <v-col>
        <v-card-title>Modifier la session de TrivialPoursuit</v-card-title>
      </v-col>
      <v-col style="text-align: right">
        <v-btn icon flat class="mx-2" @click="emit('close')">
          <v-icon icon="mdi-close"></v-icon>
        </v-btn>
      </v-col>
    </v-row>
    <v-card-text>
      <v-list>
        <v-list-subheader>
          <h3>Choix des questions</h3>
          <small
            >Chaque catégorie est définie par une <i>union</i> d'<i
              >intersections</i
            >
            d'étiquettes.</small
          >
        </v-list-subheader>
        <v-alert
          class="mx-2 mt-1"
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
              <v-list-item v-for="tags in hint.Missing!">
                <TagChip :tag="tag" :key="tag" v-for="tag in tags"></TagChip>
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
              ></TagChip
            ></span>
            sont utilisées.
          </div>
        </v-alert>
        <v-list-item
          v-for="(categorie, index) in props.edited.Questions"
          rounded
          :style="{
            'border-color': colors[index],
            borderWidth: '2px',
            borderStyle: 'solid'
          }"
          class="my-2"
        >
          <v-list-item-subtitle>Catégorie {{ index + 1 }}</v-list-item-subtitle>
          <tags-selector
            :all-tags="allKnownTags"
            :model-value="categorie || []"
            @update:model-value="v => updateCategorie(index, v)"
          ></tags-selector>
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
  CheckMissingQuestionsOut,
  QuestionCriterion,
  TrivialConfig
} from "@/controller/api_gen";
import { controller } from "@/controller/controller";
import { colorsPerCategorie } from "@/controller/trivial";
import { onMounted } from "@vue/runtime-core";
import { $ref } from "vue/macros";
import TagChip from "../editor/utils/TagChip.vue";
import TagsSelector from "./TagsSelector.vue";

interface Props {
  edited: TrivialConfig;
  allKnownTags: string[];
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "close"): void;
  (e: "update", v: TrivialConfig): void;
}>();

const colors = colorsPerCategorie;

function updateCategorie(index: number, cat: QuestionCriterion) {
  props.edited.Questions[index] = cat;
  fetchHint();
}

onMounted(fetchHint);

let hint = $ref<CheckMissingQuestionsOut>({ Pattern: [], Missing: [] });
async function fetchHint() {
  const criteria = props.edited.Questions;
  // fetch the hint if the all categories have been filled,
  // to avoid useless queries
  if (!criteria.every(qu => qu?.length)) {
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
