<template>
  <v-row>
    <v-col cols="12">
      <v-select
        density="compact"
        variant="outlined"
        :model-value="matiereTag"
        @update:model-value="(t) => updateMatiere(t as string)"
        :color="MatiereColor"
        label="Matière"
        :items="Object.keys(MatiereTagLabels)"
        hide-details
      ></v-select>
    </v-col>
  </v-row>
  <v-row>
    <v-col cols="6">
      <v-combobox
        :model-value="levelTag"
        @update:model-value="(t) => updateLevel(t as string)"
        :items="levelTags"
        label="Niveau (classe)"
        variant="outlined"
        density="compact"
        :color="LevelColor"
        hide-details
      ></v-combobox>
    </v-col>
    <v-col cols="6">
      <v-combobox
        :model-value="subLevelTags"
        @update:model-value="(t) => updateSubLevels(t as string[])"
        label="Filière"
        :items="subLevelItems"
        variant="outlined"
        density="compact"
        :color="SubLevelColor"
        multiple
        hint="Optionnelle."
        persistent-hint
        chips
      ></v-combobox>
    </v-col>
  </v-row>
  <!-- sublevel -->
  <v-row>
    <v-col cols="12">
      <v-combobox
        :model-value="chapterTag"
        @update:model-value="(t) => updateChapter(t as string)"
        :items="chapterTags"
        label="Chapitre"
        variant="outlined"
        density="compact"
        :color="ChapterColor"
        hide-details
      ></v-combobox>
    </v-col>
  </v-row>
  <v-row>
    <v-col cols="12">
      <v-combobox
        :model-value="trivMathTag"
        @update:model-value="(t) => updateIsyTriv(t as string[])"
        label="Isy'Triv"
        :items="trivMathTags"
        variant="outlined"
        density="compact"
        :color="IsyTrivColor"
        multiple
        hint="Etiquettes supplémentaires permettant de définir les catégories de questions dans l'activité Isy'Triv."
        persistent-hint
        chips
        closable-chips
      ></v-combobox>
    </v-col>
  </v-row>

  <v-alert
    v-if="!areTagsDistinct"
    color="warning"
    density="compact"
    class="ml-2"
  >
    Les étiquettes doivent être distinctes.
  </v-alert>
</template>

<script setup lang="ts">
import {
  MatiereTagLabels,
  Section,
  type TagsDB,
  type TagSection,
} from "@/controller/api_gen";
import {
  ChapterColor,
  LevelColor,
  SubLevelColor,
  IsyTrivColor,
  MatiereColor,
  tagString,
} from "@/controller/editor";
import { computed } from "vue";

interface Props {
  modelValue: TagSection[];
  allTags: TagsDB;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "update:model-value", v: TagSection[]): void;
}>();

const matiereTag = computed(
  () => props.modelValue.find((ts) => ts.Section == Section.Matiere)?.Tag || ""
);
function updateMatiere(t: string) {
  const newL = props.modelValue
    .filter((ts) => ts.Section != Section.Matiere)
    .concat([{ Section: Section.Matiere, Tag: t }]);
  emit("update:model-value", newL);
}
const levelTag = computed(
  () => props.modelValue.find((ts) => ts.Section == Section.Level)?.Tag || ""
);
function updateLevel(t: string) {
  const newL = props.modelValue
    .filter((ts) => ts.Section != Section.Level)
    .concat([{ Section: Section.Level, Tag: t }]);
  emit("update:model-value", newL);
}
const chapterTag = computed(
  () => props.modelValue.find((ts) => ts.Section == Section.Chapter)?.Tag || ""
);
function updateChapter(t: string) {
  const newL = props.modelValue
    .filter((ts) => ts.Section != Section.Chapter)
    .concat([{ Section: Section.Chapter, Tag: t }]);
  emit("update:model-value", newL);
}

const trivMathTag = computed(() =>
  props.modelValue
    .filter((ts) => ts.Section == Section.TrivMath)
    .map((ts) => ts.Tag)
);
function updateIsyTriv(t: string[]) {
  const newL = props.modelValue
    .filter((ts) => ts.Section != Section.TrivMath)
    .concat(t.map((s) => ({ Section: Section.TrivMath, Tag: s })));
  emit("update:model-value", newL);
}

const subLevelTags = computed(() =>
  props.modelValue
    .filter((ts) => ts.Section == Section.SubLevel)
    .map((ts) => ts.Tag)
);
function updateSubLevels(t: string[]) {
  const newL = props.modelValue
    .filter((ts) => ts.Section != Section.SubLevel)
    .concat(t.map((s) => ({ Section: Section.SubLevel, Tag: s })));
  emit("update:model-value", newL);
}

const levelTags = computed(() => {
  return props.allTags.Levels || [];
});
const chapterTags = computed(() => {
  return (props.allTags.ChaptersByLevel || {})[levelTag.value] || [];
});
const trivMathTags = computed(() => {
  return (
    ((props.allTags.TrivByChapters || {})[levelTag.value] || {})[
      chapterTag.value
    ] || []
  );
});
const subLevelItems = computed(() => {
  return (props.allTags.SubLevelsByLevel || {})[levelTag.value] || [];
});

const areTagsDistinct = computed(() => {
  const s = new Set<string>(props.modelValue.map((t) => tagString(t.Tag)));
  return s.size == props.modelValue.length;
});
</script>

<style scoped></style>
