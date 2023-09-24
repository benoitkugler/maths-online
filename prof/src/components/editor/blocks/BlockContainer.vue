<template>
  <v-dialog v-model="showDocumentationDefault" max-width="1000px">
    <latex-commands @close="showDocumentationDefault = false"></latex-commands>
  </v-dialog>

  <v-dialog v-model="showDocumentationExpression" max-width="1000px">
    <expression-field-doc
      @close="showDocumentationExpression = false"
    ></expression-field-doc>
  </v-dialog>

  <v-card class="mb-2" elevation="3">
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
          :icon="props.hideContent ? 'mdi-chevron-down' : 'mdi-chevron-up'"
          size="x-small"
          variant="tonal"
          class="mx-2"
          @click="emit('toggleContent')"
        ></v-btn>
        <v-btn
          v-if="props.kind == BlockKind.ExpressionFieldBlock"
          class="mr-2"
          icon
          title="Ajouter un texte avec des conseils de syntaxe."
          size="x-small"
          @click="emit('addSyntaxHint')"
        >
          <v-icon small color="green">mdi-tooltip-plus</v-icon>
        </v-btn>
        <v-btn
          v-if="showLaTeXDoc"
          class="mr-2"
          icon
          title="Documentation de la syntaxe LaTeX"
          size="x-small"
          @click="showDocumentation()"
        >
          <v-icon small color="info">mdi-help</v-icon>
        </v-btn>
        <v-btn icon title="Supprimer" size="x-small" @click="emit('delete')">
          <v-icon small color="red">mdi-close</v-icon>
        </v-btn>
      </v-col>
    </v-row>

    <v-card-text class="pt-1 pb-2" :hidden="props.hideContent">
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
import ExpressionFieldDoc from "./ExpressionFieldDoc.vue";

const emit = defineEmits<{
  (e: "delete"): void;
  (e: "toggleContent"): void;
  (e: "addSyntaxHint"): void;
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
  return isAnswer.value ? "bg-pink-lighten-3" : "bg-purple-lighten-4";
});

const showLaTeXDoc = computed(() => {
  return ![BlockKind.NumberFieldBlock].includes(props.kind);
});

function onDragStart(payload: DragEvent) {
  onDragListItemStart(payload, props.index);
}

let showDocumentationDefault = $ref(false);
let showDocumentationExpression = $ref(false);
function showDocumentation() {
  if (props.kind == BlockKind.ExpressionFieldBlock) {
    showDocumentationExpression = true;
  } else {
    showDocumentationDefault = true;
  }
}
</script>

<style scoped></style>
