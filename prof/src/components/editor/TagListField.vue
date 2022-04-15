<template>
  <v-dialog v-model="isEditing">
    <v-card subtitle="Modifier les étiquettes de la question">
      <v-card-text>
        <v-row>
          <v-col>
            <v-chip
              v-for="(tag, index) in tmpList"
              :key="tag"
              label
              closable
              class="ma-1"
              color="primary"
              @click:close="e => onDelete(e, index)"
              >{{ tag }}</v-chip
            >
          </v-col>
        </v-row>
        <v-row>
          <v-col>
            <v-autocomplete
              :items="allTags"
              variant="outlined"
              density="comfortable"
              hide-details
              label="Ajouter..."
              :search="entry"
              @update:search="s => (entry = s)"
              @keyup="onEnterKey"
              append-icon=""
              hide-no-data
              auto-select-first
              hide-selected
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
                  <v-icon icon="mdi-plus"></v-icon>
                </v-btn>
              </template>
            </v-autocomplete>
          </v-col>
        </v-row>
        <!-- <v-row>
          <v-col>
            <v-text-field
              hide-details
              label="Ajouter..."
              v-model="entry"
              @keyup="onEnter"
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
                  <v-icon icon="mdi-plus"></v-icon>
                </v-btn>
              </template>
            </v-text-field>
          </v-col>
        </v-row> -->
      </v-card-text>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn class="my-1" color="success" @click="endEdit">
          Enregistrer
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>

  <v-btn variant="outlined" class="mt-2" @click="onEdit" color="secondary">
    <v-chip
      v-for="tag in props.modelValue"
      :key="tag"
      size="small"
      label
      class="ma-1"
      color="primary"
      style="cursor: pointer"
      >{{ tag }}</v-chip
    >
    <div v-if="props.modelValue.length == 0">Ajouter une étiquette...</div>
  </v-btn>
</template>

<script setup lang="ts">
import { computed } from "@vue/runtime-core";
import { $ref } from "vue/macros";

interface Props {
  modelValue: string[];
  allTags: string[];
  label?: string;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "update:model-value", v: string[]): void;
}>();

let isEditing = $ref(false);
let tmpList = $ref<string[]>([]);

let entry = $ref("");

const isEntryValid = computed(() => {
  const e = entry.trim();
  if (e.length == 0) {
    return false;
  }
  return !tmpList.includes(e);
});

function onEdit() {
  isEditing = true;
  tmpList = props.modelValue.map(v => v.toUpperCase());
}

function add() {
  tmpList.push(entry.toUpperCase());
  entry = "";
}

function onEnterKey(key: KeyboardEvent) {
  if (key.key == "Enter" && isEntryValid.value) {
    add();
  }
}

function onDelete(e: Event, index: number) {
  tmpList = tmpList.filter((_, i) => i != index);
}

function endEdit() {
  isEditing = false;
  emit("update:model-value", tmpList);
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
