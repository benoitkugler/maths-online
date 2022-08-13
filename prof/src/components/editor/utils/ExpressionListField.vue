<template>
  <v-row no-gutters>
    <v-col cols="8" align-self="center">
      <v-autocomplete
        variant="outlined"
        density="compact"
        multiple
        :label="label"
        chips
        :model-value="props.modelValue.map((v, i) => ({ value: i, title: v }))"
        @update:model-value="onDelete"
        hide-no-data
        closable-chips
        readonly
        clearable
        append-inner-icon=""
        :hint="props.hint"
        :persistent-hint="!!props.hint"
      >
      </v-autocomplete>
    </v-col>
    <v-col cols="4">
      <v-text-field
        class="fix-input-width"
        variant="outlined"
        density="compact"
        label="Ajouter"
        :color="color"
        v-model="entry"
        @click:append-inner="add"
        hint="Ajouter une expression"
        @keyup="onEnter"
      >
        <template v-slot:append-inner>
          <v-btn icon size="x-small" :disabled="!isEntryValid" @click="add">
            <v-icon icon="mdi-plus" color="green"></v-icon>
          </v-btn>
        </template>
      </v-text-field>
    </v-col>
  </v-row>
</template>

<script setup lang="ts">
import { ExpressionColor } from "@/controller/editor";
import { computed } from "@vue/runtime-core";
import { $ref } from "vue/macros";

interface Props {
  modelValue: string[];
  label?: string;
  hint?: string;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "update:model-value", v: string[]): void;
}>();

const color = ExpressionColor;

function onDelete(indices: number[]) {
  emit(
    "update:model-value",
    indices.map((index) => props.modelValue[index])
  );
}

let entry = $ref("");
const isEntryValid = computed(() => !!entry.length);

function add() {
  props.modelValue.push(entry);
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
