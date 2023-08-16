<template>
  <v-card title="Modifier les étiquettes">
    <v-card-text class="my-1">
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
      <!-- sublevel -->
      <v-row>
        <v-col cols="12">
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
      <v-row>
        <v-col cols="12">
          <v-combobox
            :model-value="trivMathTag"
            @update:model-value="(t) => updateIsyTriv(t as string[])"
            label="Isy'Triv"
            :items="trivMathTags"
            variant="outlined"
            density="compact"
            :color="TrivMathColor"
            multiple
            hint="Etiquettes supplémentaires permettant de définir les catégories de questions dans l'activité Isy'Triv."
            persistent-hint
            chips
          ></v-combobox>
        </v-col>
      </v-row>
    </v-card-text>
    <v-card-actions>
      <v-spacer></v-spacer>
      <v-btn
        class="my-1"
        color="success"
        @click="emit('save')"
        :disabled="!saveEnabled"
        variant="outlined"
      >
        Enregistrer
      </v-btn>
    </v-card-actions>
  </v-card>
</template>

<script setup lang="ts">
import { Section, type TagsDB, type TagSection } from "@/controller/api_gen";
import {
  ChapterColor,
  LevelColor,
  SubLevelColor,
  TrivMathColor
} from "@/controller/editor";
import { computed } from "@vue/runtime-core";
import { $computed } from "vue/macros";

interface Props {
  modelValue: TagSection[];
  allTags: TagsDB;
  saveEnabled: boolean;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "update:model-value", v: TagSection[]): void;
  (e: "save"): void;
}>();

let levelTag = $computed(
  () => props.modelValue.find(ts => ts.Section == Section.Level)?.Tag || ""
);
function updateLevel(t: string) {
  const newL = props.modelValue
    .filter(ts => ts.Section != Section.Level)
    .concat([{ Section: Section.Level, Tag: t }]);
  emit("update:model-value", newL);
}
let chapterTag = $computed(
  () => props.modelValue.find(ts => ts.Section == Section.Chapter)?.Tag || ""
);
function updateChapter(t: string) {
  const newL = props.modelValue
    .filter(ts => ts.Section != Section.Chapter)
    .concat([{ Section: Section.Chapter, Tag: t }]);
  emit("update:model-value", newL);
}

let trivMathTag = $computed(() =>
  props.modelValue
    .filter(ts => ts.Section == Section.TrivMath)
    .map(ts => ts.Tag)
);
function updateIsyTriv(t: string[]) {
  const newL = props.modelValue
    .filter(ts => ts.Section != Section.TrivMath)
    .concat(t.map(s => ({ Section: Section.TrivMath, Tag: s })));
  emit("update:model-value", newL);
}

const subLevelTags = $computed(() =>
  props.modelValue
    .filter(ts => ts.Section == Section.SubLevel)
    .map(ts => ts.Tag)
);
function updateSubLevels(t: string[]) {
  const newL = props.modelValue
    .filter(ts => ts.Section != Section.SubLevel)
    .concat(t.map(s => ({ Section: Section.SubLevel, Tag: s })));
  emit("update:model-value", newL);
}

const levelTags = computed(() => {
  return props.allTags.Levels || [];
});
const chapterTags = computed(() => {
  return (props.allTags.ChaptersByLevel || {})[levelTag] || [];
});
const trivMathTags = computed(() => {
  return (
    ((props.allTags.TrivByChapters || {})[levelTag] || {})[chapterTag] || []
  );
});
const subLevelItems = computed(() => {
  return (props.allTags.SubLevelsByLevel || {})[levelTag] || [];
});
</script>

<style scoped></style>
