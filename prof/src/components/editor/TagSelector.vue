<template>
  <v-autocomplete
    variant="outlined"
    density="compact"
    multiple
    chips
    closable-chips
    :items="tagItems"
    :search="search"
    @update:search="onSearch"
    color="primary"
    label="Catégories"
    no-data-text="Aucune étiquette ne correspond à votre recherche."
    :model-value="props.modelValue"
    @update:model-value="onUpdate"
    hint="Restreint la recherche à l'intersection des catégories sélectionnées."
    persistent-hint
  ></v-autocomplete>
</template>

<script setup lang="ts">
import type { TagsDB } from "@/controller/api_gen";
import { filterTags } from "@/controller/editor";
import { watch } from "vue";
import { $ref } from "vue/macros";

interface Props {
  modelValue: string[];
  allTags: TagsDB;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "update:model-value", v: string[]): void;
}>();

// let search = $ref("");

// function onUpdate(v: any) {
//   emit("update:model-value", v);
// }

// let tagItems = $ref(filterTags(props.allTags, "", []));

// function onSearch(s: string) {
//   search = s;
//   tagItems = filterTags(props.allTags, s, []);
// }

// watch(props, () => {
//   onSearch("");
// });
</script>

<style scoped></style>
