<template>
  <v-card max-height="80vh" class="overflow-y-auto py-2">
    <v-row>
      <v-col>
        <v-card-title>Modifier les réglages Isy'Triv</v-card-title>
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
        <v-col cols="12" md="auto">
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
        <v-spacer></v-spacer>
        <v-col cols="12" md="auto" align-self="center" class="mb-1">
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
                  :key="index"
                  v-for="(tag, index) in tags || []"
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
                :key="index"
                v-for="(tag, index) in hint.Pattern"
                :pointer="false"
              ></TagChip
            ></span>
            sont utilisées.
          </div>
        </v-alert>
        <CategorieRow
          v-for="(categorie, index) in props.edited.Questions.Tags"
          :key="index"
          :index="index"
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
                :last-matiere-level-chapter="lastMatiereLevelChapter"
                @update:model-value="(v) => updateCategorie(index, v)"
              ></tags-selector>
            </v-col>
          </v-row>
        </CategorieRow>
      </v-list>

      <v-row class="mt-2">
        <v-col cols="12" md="6">
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
        <v-col cols="12" md="6">
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
import {
  Section,
  type CheckMissingQuestionsOut,
  type DifficultyTag,
  type QuestionCriterion,
  type Tags,
  type TagsDB,
  type Trivial,
} from "@/controller/api_gen";
import { controller } from "@/controller/controller";
import { ref, computed, onMounted } from "vue";
import TagChip from "../editor/utils/TagChip.vue";
import DifficultyChoices from "./DifficultyChoices.vue";
import TagsSelector from "./TagsSelector.vue";
import CategorieRow from "./CategorieRow.vue";
import type { PrefillTrivialCategorie } from "@/controller/utils";

interface Props {
  edited: Trivial;
  allKnownTags: TagsDB;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "closed"): void;
  (e: "update", v: Trivial): void;
}>();

const lastMatiereLevelChapter = computed<PrefillTrivialCategorie>(() => {
  const all = allTags(props.edited.Questions.Tags);
  if (all.length) {
    const last = all[all.length - 1] || [];
    return {
      matiere: last.find((s) => s.Section == Section.Matiere)?.Tag || "",
      level: last.find((s) => s.Section == Section.Level)?.Tag || "",
      chapter: last.find((s) => s.Section == Section.Chapter)?.Tag || "",
      sublevels: last.filter((s) => s.Section == Section.SubLevel),
    };
  }
  return { matiere: "", level: "", chapter: "", sublevels: [] };
});

const showDifficultyCard = ref(false);

function onEditDifficulties(difficulties: DifficultyTag[]) {
  showDifficultyCard.value = false;
  props.edited.Questions.Difficulties = difficulties;
  fetchHintForMissing();
}

function allTags(tags: QuestionCriterion[]) {
  const all: Tags[] = [];
  tags.forEach((l) => all.push(...(l || []).map((ls) => ls || [])));
  return all;
}

function updateCategorie(index: number, cat: QuestionCriterion) {
  props.edited.Questions.Tags[index] = cat;
  fetchHintForMissing();
}

onMounted(fetchHintForMissing);

const hint = ref<CheckMissingQuestionsOut>({ Pattern: [], Missing: [] });
async function fetchHintForMissing() {
  const criteria = props.edited.Questions;
  // fetch the hint only if the all categories have been filled,
  // to avoid useless queries
  if (!criteria.Tags.every((qu) => qu?.length)) {
    hint.value = { Pattern: [], Missing: [] };
    return;
  }
  const res = await controller.CheckMissingQuestions(criteria);
  if (res == undefined) {
    return;
  }
  hint.value = res;
}
</script>

<style scoped></style>
