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
import { filterTags, tagColor, tagString } from "@/controller/editor";
import { computed } from "@vue/runtime-core";
import { $ref } from "vue/macros";

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
  setTimeout(() => {
    onSearch("");
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

let tagItems = $ref(filterTags(props.allTags, "", props.modelValue));

function onSearch(s: string) {
  search = s;
  tagItems = filterTags(props.allTags, s, props.modelValue);
}
</script>

<style scoped></style>
