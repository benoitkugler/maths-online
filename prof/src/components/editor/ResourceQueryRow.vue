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
    <v-col cols="2">
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
    <v-col cols="2">
      <v-select
        label="Filière"
        variant="outlined"
        density="compact"
        multiple
        :items="subLevelTags"
        :model-value="props.modelValue.SubLevelTags || []"
        @update:model-value="
          (v) => {
            props.modelValue.SubLevelTags = v;
            emit('update:model-value', props.modelValue);
          }
        "
        :color="SubLevelColor"
        chips
        closable-chips
        hide-details
        no-data-text="Aucune filière pour le niveau choisi."
      ></v-select>
    </v-col>
    <v-col cols="3">
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
        :no-data-text="noChaptersText"
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
      <MatiereSelect
        :matiere="props.modelValue.Matiere"
        @update:matiere="
          (o) => {
            props.modelValue.Matiere = o;
            emit('update:model-value', props.modelValue);
          }
        "
      ></MatiereSelect>
    </v-col>
  </v-row>
</template>

<script setup lang="ts">
import type { LevelTag, Query, TagsDB } from "@/controller/api_gen";
import { ChapterColor, LevelColor, SubLevelColor } from "@/controller/editor";
import { computed } from "vue";
import OriginSelect from "../OriginSelect.vue";
import MatiereSelect from "../MatiereSelect.vue";

interface Props {
  modelValue: Query;
  allTags: TagsDB;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "update:model-value", v: Query): void;
}>();

const levelTags = computed(() => {
  return (props.allTags.Levels || [])
    .map((tag) => ({
      title: tag as LevelTag | "Non classé",
      value: tag as string,
    }))
    .concat({ title: "Non classé", value: "" });
});

const chapterTags = computed(() => {
  const all = new Set<string>();
  props.modelValue.LevelTags?.forEach((levelTag) => {
    const l = (props.allTags.ChaptersByLevel || {})[levelTag as LevelTag];
    l?.forEach((v) => all.add(v));
  });
  const out = Array.from(all.values());
  out.sort();
  return out
    .map((tag) => ({
      title: tag,
      value: tag,
    }))
    .concat({ title: "Non classé", value: "" });
});

const subLevelTags = computed(() => {
  const all = new Set<string>();
  props.modelValue.LevelTags?.forEach((levelTag) => {
    const l = (props.allTags.SubLevelsByLevel || {})[levelTag as LevelTag];
    l?.forEach((v) => all.add(v));
  });
  const out = Array.from(all.values());
  out.sort();
  return out.map((tag) => ({
    title: tag,
    value: tag,
  }));
  // .concat({ title: "Non classé", value: "" });
});

const noChaptersText = computed(() => {
  return props.modelValue.LevelTags?.length
    ? "Aucun chapitre n'est disponible"
    : "Selectionnez d'abord un niveau.";
});

// debounce feature for text field
let timerId: ReturnType<typeof setTimeout>;
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
