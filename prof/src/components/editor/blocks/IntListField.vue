<template>
  <v-row>
    <v-col cols="8">
      <v-autocomplete
        variant="outlined"
        density="compact"
        multiple
        :label="label"
        chips
        :model-value="
          props.modelValue.map((v, i) => ({ value: i, title: String(v) }))
        "
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
        class="fix-input-width"
        variant="outlined"
        density="compact"
        label="Ajouter"
        color="green"
        v-model="entry"
        @click:append-inner="add"
        :hint="entryHint"
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
  disallowIntervals?: boolean;
  disallowRepeat?: boolean;
  sorted: boolean;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "update:model-value", v: number[]): void;
}>();

function onDelete(s: number[]) {
  emit(
    "update:model-value",
    s.map(index => props.modelValue[index])
  );
}

let entry = $ref("");

const isEntryValid = computed(() => parseEntry(entry) !== undefined);

const entryHint = computed(() => {
  let out = "";
  if (!props.disallowIntervals) {
    out += "Ajouter un intervalle avec a:b. ";
  }
  if (!props.disallowRepeat) {
    out += "Ajouter a fois le nombre b avec a*b. ";
  }
  return out;
});

function parseEntry(s: string): number[] {
  s = s.trim();
  if (s == "") {
    return [];
  }
  let chunks = s.split(",");
  if (chunks.length == 1) {
    chunks = s.split(";");
  }
  if (chunks.length > 1) {
    return chunks.map(v => Number(v.trim())).filter(n => !isNaN(n));
  }

  chunks = s.split(":");
  if (chunks.length == 2) {
    const start = Number(chunks[0].trim());
    const end = Number(chunks[1].trim());
    if (isNaN(start) || isNaN(end)) {
      return [];
    }
    const out = [];
    for (let index = start; index <= end; index++) {
      out.push(index);
    }
    return out;
  }

  chunks = s.split("*");
  if (chunks.length == 2) {
    const a = Number(chunks[0].trim());
    const b = Number(chunks[1].trim());
    if (isNaN(a) || isNaN(b)) {
      return [];
    }
    return Array.from({ length: a }, (v, i) => b);
  }

  const v = Number(s);
  if (isNaN(v)) {
    return [];
  }
  return [v];
}

function add() {
  const values = parseEntry(entry)!;
  props.modelValue.push(...values);
  if (props.sorted) {
    props.modelValue.sort((a, b) => a - b);
  }
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

.fix-input-width:deep(input) {
  width: 100%;
}
</style>
