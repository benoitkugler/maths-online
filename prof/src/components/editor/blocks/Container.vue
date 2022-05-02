<template>
  <v-card class="my-2" elevation="3">
    <v-row
      no-gutters
      :class="'px-2 rounded ' + colorClass"
      @dragstart="onDragStart"
      draggable="true"
    >
      <v-col align-self="center" cols="8">
        <v-card-subtitle>
          <b>{{ kindLabels[props.kind].label }}</b>
          <span v-if="isAnswer" class="ml-1">(Champ de réponse)</span>
        </v-card-subtitle>
      </v-col>
      <v-col cols="4" style="text-align: right" class="my-2">
        <v-btn icon flat title="Supprimer" size="x-small">
          <v-icon small color="red" @click="emit('delete')">mdi-close</v-icon>
        </v-btn>
      </v-col>
    </v-row>

    <v-card-text class="pt-1 pb-2" :hidden="hideContent">
      <slot></slot>
    </v-card-text>
  </v-card>
</template>

<script setup lang="ts">
import { BlockKindLabels, onDragListItemStart } from "@/controller/editor";
import type { BlockKind } from "@/controller/exercice_gen";
import { computed } from "@vue/runtime-core";

const emit = defineEmits<{
  (e: "delete"): void;
}>();

interface ContainerProps {
  index: number;
  kind: BlockKind;
  hideContent: boolean;
  hasError: boolean;
}

const props = defineProps<ContainerProps>();

const kindLabels = BlockKindLabels;

const isAnswer = computed(() => kindLabels[props.kind].isAnswerField);

const colorClass = computed(() => {
  if (props.hasError) {
    return "bg-red";
  }
  return isAnswer.value ? "bg-pink-lighten-3" : "bg-purple-lighten-3";
});

function onDragStart(payload: DragEvent) {
  onDragListItemStart(payload, props.index);
}
</script>

<style scoped>
.small-slider:deep(.v-input__control) {
  min-height: 200px !important;
}
</style>
