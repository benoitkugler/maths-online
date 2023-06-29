<template>
  <v-row>
    <v-col align-self="center">
      <v-text-field
        label="Rechercher une ressource par son titre"
        variant="outlined"
        density="compact"
        v-model="props.modelValue.TitleQuery"
        @update:model-value="updateQuerySearch"
        hide-details
        clearable
      ></v-text-field>
    </v-col>
    <v-col cols="3">
      <v-select
        multiple
        :items="levelTags"
        :model-value="props.modelValue.LevelTags || []"
        @update:model-value="
          (v) => {
            props.modelValue.LevelTags = v;
            emit('update:model-value', props.modelValue);
          }
        "
        label="Niveau (classe)"
        variant="outlined"
        density="compact"
        :color="LevelColor"
        chips
        closable-chips
        hide-details
      ></v-select>
    </v-col>
    <v-col cols="4">
      <v-select
        hide-details
        multiple
        :items="chapterTags"
        :model-value="props.modelValue.ChapterTags || []"
        @update:model-value="
          (v) => {
            props.modelValue.ChapterTags = v;
            emit('update:model-value', props.modelValue);
          }
        "
        label="Chapitre"
        variant="outlined"
        density="compact"
        :color="ChapterColor"
        chips
        closable-chips
        hide-no-data
      ></v-select>
    </v-col>
    <v-col cols="auto" align-self="center">
      <OriginSelect
        :origin="props.modelValue.Origin"
        @update:origin="
          (o) => {
            props.modelValue.Origin = o;
            emit('update:model-value', props.modelValue);
          }
        "
      ></OriginSelect>
    </v-col>
  </v-row>
</template>

<script setup lang="ts">
import type { Query, TagsDB } from "@/controller/api_gen";
import { ChapterColor, LevelColor } from "@/controller/editor";
import { computed } from "vue";
import OriginSelect from "../OriginSelect.vue";

interface Props {
  modelValue: Query;
  allTags: TagsDB;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "update:model-value", v: Query): void;
}>();

const levelTags = computed(() => {
  return props.allTags.Levels || [];
});
const chapterTags = computed(() => {
  const all = new Set<string>();
  props.modelValue.LevelTags?.forEach((levelTag) => {
    const l = (props.allTags.ChaptersByLevel || {})[levelTag];
    l?.forEach((v) => all.add(v));
  });
  const out = Array.from(all.values());
  out.sort();
  return out;
});

// debounce feature for text field
let timerId = 0;
function updateQuerySearch() {
  const debounceDelay = 300;
  // cancel pending call
  clearTimeout(timerId);

  // delay new call 500ms
  timerId = setTimeout(() => {
    emit("update:model-value", props.modelValue);
  }, debounceDelay);
}
</script>

<style scoped></style>
