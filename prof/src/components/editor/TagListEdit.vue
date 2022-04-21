<template>
  <v-row>
    <v-col :cols="props.horizontal ? 7 : 12">
      <v-chip
        v-for="(tag, index) in props.modelValue"
        :key="tag"
        label
        closable
        class="ma-1"
        color="primary"
        @click:close="e => onDelete(e, index)"
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
        :search="entry"
        @update:search="s => (entry = s)"
        @keyup="onEnterKey"
        @update:model-value="onSelectItem"
        append-icon=""
        hide-no-data
        auto-select-first
      >
        <template v-slot:appendInner>
          <v-btn
            icon
            size="x-small"
            class="my-1"
            color="success"
            :disabled="!isEntryValid"
            @click="add"
          >
            <v-icon icon="mdi-plus" size="x-small"></v-icon>
          </v-btn>
        </template>
      </v-autocomplete>
    </v-col>
  </v-row>
</template>

<script setup lang="ts">
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

let entry = $ref("");

const tagItems = computed(() => {
  return props.allTags.filter(t => !props.modelValue.includes(t));
});

const isEntryValid = computed(() => {
  const e = entry.trim();
  if (e.length == 0) {
    return false;
  }
  return !props.modelValue.includes(e);
});

function add() {
  props.modelValue.push(entry.toUpperCase());
  entry = "";
  emit("update:model-value", props.modelValue);
}

function onEnterKey(key: KeyboardEvent) {
  if (key.key == "Enter" && isEntryValid.value) {
    add();
  }
}

function onSelectItem(s: string) {
  s = s.toUpperCase();
  if (props.modelValue.includes(s)) {
    entry = "";
    return;
  }
  entry = s;
  add();
}

function onDelete(e: Event, index: number) {
  props.modelValue.splice(index, 1);
  emit("update:model-value", props.modelValue);
}
</script>

<style scoped></style>
