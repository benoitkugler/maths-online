<template>
  <v-dialog v-model="isEditing">
    <v-card subtitle="Modifier les étiquettes de la question">
      <v-card-text>
        <tag-list-edit
          v-model="tmpList"
          :horizontal="false"
          :all-tags="allTags"
        ></tag-list-edit>
      </v-card-text>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn
          class="my-1"
          color="success"
          @click="endEdit"
          :disabled="!saveEnabled"
          variant="contained"
        >
          Enregistrer
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>

  <v-btn variant="outlined" @click="startEdit" color="secondary">
    <v-chip
      v-for="tag in props.modelValue"
      :key="tag"
      size="small"
      label
      class="ma-1"
      :color="tagColor(tag)"
      style="cursor: pointer"
      >{{ tag }}</v-chip
    >
    <div v-if="props.modelValue.length == 0">Ajouter une étiquette...</div>
  </v-btn>
</template>

<script setup lang="ts">
import { tagColor, tagString } from "@/controller/editor";
import { computed } from "@vue/runtime-core";
import { $ref } from "vue/macros";
import TagListEdit from "./TagListEdit.vue";

interface Props {
  modelValue: string[];
  allTags: string[];
  label?: string;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "update:model-value", v: string[]): void;
}>();

defineExpose({ startEdit });

let isEditing = $ref(false);
let tmpList = $ref<string[]>([]);

function startEdit() {
  isEditing = true;
  tmpList = props.modelValue.map(v => tagString(v));
}

const saveEnabled = computed(() => {
  if (tmpList.length != props.modelValue.length) {
    return true;
  }
  return !tmpList.every((tag, index) => props.modelValue[index] == tag);
});

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
