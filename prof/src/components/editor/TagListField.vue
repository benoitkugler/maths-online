<template>
  <v-dialog v-model="isEditing" :retain-focus="false" max-width="800">
    <v-card title="Modifier les étiquettes de la question">
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
          variant="outlined"
        >
          Enregistrer
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>

  <v-sheet
    variant="outlined"
    rounded
    border="secondary"
    :style="{
      'border-width': '2px',
      cursor: props.readonly ? '' : 'pointer',
      'text-align': 'center',
    }"
    :class="props.yPadding ? 'py-1' : ''"
  >
    <v-btn
      v-if="props.modelValue.length == 0"
      @click.stop="startEdit"
      flat
      size="x-small"
      class="py-0 px-2"
      :disabled="props.readonly"
      block
    >
      {{ props.readonly ? "Aucune étiquette" : "Ajouter une étiquette..." }}
    </v-btn>

    <v-row no-gutters v-else justify="center" @click.stop="startEdit">
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
  yPadding?: boolean;
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
