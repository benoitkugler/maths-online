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

  <v-sheet
    variant="outlined"
    @click.stop="startEdit"
    rounded
    border="secondary"
    :style="{ 'border-width': '2px', cursor: props.readonly ? '' : 'pointer' }"
    :class="props.yPadding ? 'py-1' : ''"
  >
    <div
      v-if="props.modelValue.length == 0"
      style="text-align: center; font-style: italic"
      class="pa-2"
    >
      Ajouter une étiquette...
    </div>
    <v-row no-gutters v-else justify="center">
      <v-col v-for="tag in props.modelValue" :key="tag" cols="auto">
        <tag-chip :tag="tag" :pointer="!props.readonly"> </tag-chip>
      </v-col>
    </v-row>
  </v-sheet>
</template>

<script setup lang="ts">
import { tagString } from "@/controller/editor";
import { computed } from "@vue/runtime-core";
import { $ref } from "vue/macros";
import TagListEdit from "./TagListEdit.vue";
import TagChip from "./utils/TagChip.vue";

interface Props {
  modelValue: string[];
  allTags: string[];
  readonly: boolean;
  label?: string;
  yPadding: boolean;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: "update:model-value", v: string[]): void;
}>();

defineExpose({ startEdit });

let isEditing = $ref(false);
let tmpList = $ref<string[]>([]);

function startEdit() {
  if (props.readonly) {
    return;
  }
  isEditing = true;
  tmpList = props.modelValue.map((v) => tagString(v));
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
