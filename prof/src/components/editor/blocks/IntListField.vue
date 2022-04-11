<template>
  <v-row>
    <v-col>
      <v-autocomplete
        variant="outlined"
        density="compact"
        multiple
        :label="label"
        chips
        :model-value="props.modelValue.map(v => String(v))"
        @update:model-value="onDelete"
        hide-no-data
        closable-chips
        readonly
        clearable
        append-inner-icon=""
      >
      </v-autocomplete>
    </v-col>
    <v-col cols="4">
      <v-text-field
        variant="outlined"
        density="compact"
        label="Ajouter"
        color="green"
        v-model="entry"
        @click:append-inner="add"
        hint="Ajouter un intervalle avec a:b"
        @keyup="onEnter"
      >
        <template v-slot:appendInner>
          <v-btn icon size="x-small" :disabled="!isEntryValid" @click="add">
            <v-icon icon="mdi-plus" color="green"></v-icon>
          </v-btn>
        </template>
      </v-text-field>
    </v-col>
  </v-row>
</template>

<script setup lang="ts">
import { computed } from "@vue/runtime-core";
import { $ref } from "vue/macros";

interface Props {
  modelValue: number[];
  label?: string;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "update:model-value", v: number[]): void;
}>();

function onDelete(s: string[]) {
  emit(
    "update:model-value",
    s.map(s => Number(s))
  );
}

let entry = $ref("");

const isEntryValid = computed(() => parseEntry(entry) !== undefined);

function parseEntry(s: string) {
  s = s.trim();
  if (s == "") {
    return;
  }
  let chunks = s.split(",");
  if (chunks.length != 2) {
    chunks = s.split(";");
  }
  if (chunks.length != 2) {
    chunks = s.split(":");
  }
  if (chunks.length == 2) {
    const start = Number(chunks[0].trim());
    const end = Number(chunks[1].trim());
    if (isNaN(start) || isNaN(end)) {
      return;
    }
    return { start, end };
  }

  const v = Number(s);
  if (isNaN(v)) {
    return;
  }
  return { start: v, end: v };
}

function add() {
  const range = parseEntry(entry)!;
  for (let index = range.start; index <= range.end; index++) {
    props.modelValue.push(index);
  }
  props.modelValue.sort((a, b) => a - b);
  entry = "";
  emit("update:model-value", props.modelValue);
}

function onEnter(key: KeyboardEvent) {
  if (key.key == "Enter" && isEntryValid.value) {
    add();
  }
}
</script>

<style scoped>
.centered-input:deep(input) {
  text-align: center;
}

:deep(.v-field__append-inner) {
  padding-top: 4px;
}
</style>
