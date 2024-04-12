<template>
  <v-text-field
    :hide-details="!props.hint"
    :hint="props.hint"
    :label="props.label"
    readonly
    variant="outlined"
    density="compact"
    model-value=" "
  >
    <template v-slot:append>
      <v-text-field
        class="fix-input-width"
        variant="outlined"
        density="comfortable"
        label="Ajouter"
        :color="color"
        v-model="entry"
        @click:append-inner="add"
        hint="Ajouter une expression"
        persistent-hint
        @keyup="onEnter"
      >
        <template v-slot:append-inner>
          <v-btn icon size="x-small" :disabled="!isEntryValid" @click="add">
            <v-icon icon="mdi-plus" color="green"></v-icon>
          </v-btn>
        </template>
      </v-text-field>
    </template>

    <v-chip
      v-for="(v, i) in props.modelValue"
      :key="v"
      closable
      size="small"
      @click:close="onDelete(i)"
    >
      {{ v }}
    </v-chip>
  </v-text-field>
</template>

<script setup lang="ts">
import { ExpressionColor } from "@/controller/editor";
import { computed } from "vue";
import { ref } from "vue";

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

function onDelete(indice: number) {
  emit("update:model-value", props.modelValue.toSpliced(indice, 1));
}

const entry = ref("");
const isEntryValid = computed(() => !!entry.value.length);

function add() {
  emit("update:model-value", props.modelValue.concat(entry.value));
  entry.value = "";
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
