<template>
  <v-row>
    <v-col :cols="props.horizontal ? 7 : 12">
      <v-chip
        v-for="(tag, index) in props.modelValue"
        :key="tag"
        label
        closable
        class="ma-1"
        :color="tagColor(tag)"
        @click:close="(e) => onDelete(e, index)"
        >{{ tag }}</v-chip
      >
    </v-col>
    <v-col :cols="props.horizontal ? 5 : 12">
      <v-autocomplete
        :items="tagItems"
        variant="outlined"
        :density="props.horizontal ? 'compact' : 'comfortable'"
        hide-details
        label="Ajouter..."
        placeholder="Rechercher une étiquette..."
        :search="search"
        @update:search="onSearch"
        @keyup="onEnterKey"
        @update:model-value="onSelectItem"
        append-icon=""
        no-data-text="Aucune étiquette ne correspond à votre recherche."
        hide-selected
        auto-select-first
        :custom-filter="filter"
      >
        <template v-slot:append-inner>
          <v-btn
            icon
            size="x-small"
            class="my-1"
            color="success"
            :disabled="!isEntryValid"
            @click="add(search)"
          >
            <v-icon icon="mdi-plus" size="x-small"></v-icon>
          </v-btn>
        </template>
      </v-autocomplete>
    </v-col>
  </v-row>
</template>

<script setup lang="ts">
import { tagColor, tagString } from "@/controller/editor";
import { computed } from "@vue/runtime-core";
import { $computed, $ref } from "vue/macros";

interface Props {
  modelValue: string[];
  allTags: string[];
  horizontal: boolean;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "update:model-value", v: string[]): void;
}>();

let search = $ref("");

const isEntryValid = computed(() => {
  const e = search.trim();
  if (e.length == 0) {
    return false;
  }
  return !props.modelValue.includes(e);
});

function add(s: string) {
  props.modelValue.push(tagString(s));
  emit("update:model-value", props.modelValue);

  // reset the search field
  search = "";
  setTimeout(() => {
    tagItems = candidatesTags.slice(0, maxRes);
    search = "";
  }, 100);
}

function onEnterKey(key: KeyboardEvent) {
  if (key.key == "Enter" && isEntryValid.value) {
    add(search);
  }
}

function onSelectItem(s: string) {
  s = tagString(s);
  if (props.modelValue.includes(s)) {
    search = "";
    return;
  }
  add(s);
}

function onDelete(e: Event, index: number) {
  props.modelValue.splice(index, 1);
  emit("update:model-value", props.modelValue);
  onSearch("");
}

// this is used to limit the number of items returned by the search
let nbResults = 0;
const maxRes = 10;

let candidatesTags = $computed(() => {
  return props.allTags.filter((t) => !props.modelValue.includes(t));
});

let tagItems = $ref(candidatesTags.slice(0, maxRes));

function onSearch(s: string) {
  search = s;
  if (s.length) {
    tagItems = candidatesTags.map((v) => v);
  } else {
    tagItems = candidatesTags.slice(0, maxRes);
  }
  nbResults = 0;
}

function filter(value: string, query: string, item: string) {
  if (nbResults >= maxRes) {
    return false;
  }
  query = query.toUpperCase();
  const start = value.indexOf(query);
  if (start != -1) {
    nbResults += 1;
    return start;
  }
  return false;
}
</script>

<style scoped></style>
