<template>
  <v-dialog v-model="showDocumentation" max-width="1000px">
    <latex-commands @close="showDocumentation = false"></latex-commands>
  </v-dialog>

  <v-card class="my-2" elevation="3">
    <v-row
      no-gutters
      :class="'rounded ' + colorClass"
      @dragstart="onDragStart"
      draggable="true"
      style="cursor: grab"
    >
      <v-col cols="auto" align-self="center">
        <v-icon size="large" icon="mdi-drag-vertical"></v-icon>
      </v-col>
      <v-col align-self="center" cols="7">
        <v-card-subtitle>
          <b>{{ kindLabels[props.kind].label }}</b>
          <span v-if="isAnswer" class="ml-1">(Champ de réponse)</span>
        </v-card-subtitle>
      </v-col>
      <v-spacer></v-spacer>
      <v-col cols="auto" style="text-align: right" class="my-1 mr-2">
        <v-btn
          v-if="showLaTeXDoc"
          class="mr-2"
          icon
          title="Documentation de la syntaxe LaTeX"
          size="x-small"
        >
          <v-icon small color="info" @click="showDocumentation = true"
            >mdi-help</v-icon
          >
        </v-btn>
        <v-btn icon title="Supprimer" size="x-small">
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
import { BlockKind } from "@/controller/api_gen";
import { BlockKindLabels } from "@/controller/editor";
import { onDragListItemStart } from "@/controller/utils";
import { computed } from "@vue/runtime-core";
import { $ref } from "vue/macros";
import LatexCommands from "./LatexCommands.vue";

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

const showLaTeXDoc = computed(() => {
  return ![BlockKind.ExpressionFieldBlock, BlockKind.NumberFieldBlock].includes(
    props.kind
  );
});

function onDragStart(payload: DragEvent) {
  onDragListItemStart(payload, props.index);
}

let showDocumentation = $ref(false);
</script>

<style scoped></style>
